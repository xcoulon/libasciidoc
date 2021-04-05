package parser

import (
	"fmt"
	"io"

	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	log "github.com/sirupsen/logrus"
)

// ParseDocumentFragments parses a document's content and applies the preprocessing directives (file inclusions)
// Returns the document fragments (to be assembled) or an error
func ParseDocumentFragments(r io.Reader, config configuration.Configuration, options ...Option) (types.DocumentFragments, error) {
	// first phase: parse lines with substitutions defined by the block attribute previously detected
	return parseDocumentFragments(newParserContext(config), r, options...)
}

func parseDocumentFragments(ctx *parserContext, r io.Reader, options ...Option) (types.DocumentFragments, error) {
	content, err := ParseReader(ctx.config.Filename, r, append(options, Entrypoint("DocumentFragments"), GlobalStore(parseContextKey, ctx))...)
	if err != nil {
		log.Errorf("failed to parse raw document: %s", err)
		return nil, err
	}
	fragments, ok := content.(types.DocumentFragments)
	if !ok {
		return nil, fmt.Errorf("unexpected type of raw lines: '%T'", content)
	}
	// look-up file inclusions and replace with the file content
	result := make([]interface{}, 0, len(fragments)) // assume the initial capacity as if there is no file inclusion to process
	for _, fragment := range fragments {
		switch fragment := fragment.(type) {
		case types.FileInclusion:
			fileContent, err := parseFileToInclude(ctx.clone(), fragment, options...) // clone the context so that further level offset in this "child" doc are not applied to the rest of this "parent" doc
			if err != nil {
				return nil, err
			}
			result = append(result, fileContent...)
		case types.AttributeDeclaration:
			// immediatly process the attribute substitutions in the value (if there is any)
			value := substituteAttributes(fragment.Value, ctx.attributes)
			switch value := value.(type) {
			case types.StringElement:
				ctx.attributes[fragment.Name] = value.Content
			case string:
				ctx.attributes[fragment.Name] = value
			default:
				return nil, fmt.Errorf("unexpected type of value after substituing attributes: '%T'", value)
			}
			// add the attribute to the current content
			result = append(result, fragment)

		case types.Section:
			// apply the level offsets on the sections of the doc to include (by default, there is no offset)
			oldLevel := fragment.Level
			for _, offset := range ctx.levelOffsets {
				oldLevel := offset.apply(&fragment)
				// replace the absolute level offset with a relative one,
				// taking into account the offset that was applied on this section,
				// so that following sections in the document have the same relative offset
				// applied, but keep the same hierarchy
				if offset.absolute {
					// replace current offset and all previous ones with a single relative offset
					ctx.levelOffsets = []levelOffset{relativeOffset(fragment.Level - oldLevel)}
				}
			}
			if log.IsLevelEnabled(log.DebugLevel) {
				log.Debugf("applied level offset on section: level %d -> %d", oldLevel, fragment.Level)
				// log.Debug("level offsets:")
				// spew.Fdump(log.StandardLogger().Out, ctx.levelOffsets)
			}
			result = append(result, fragment)
		case types.BlockDelimiter:
			ctx.onBlockDelimiter(fragment)
			result = append(result, fragment)
		case types.InlineElements:
			fragment, err := onInlineElements(ctx, fragment)
			if err != nil {
				return nil, err
			}
			result = append(result, fragment)
		default:
			result = append(result, fragment)
		}
	}
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debug("parsed file to include")
	// 	spew.Fdump(log.StandardLogger().Out, result)
	// }
	return result, nil
}
