package parser

import (
	"io"

	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

// ParseDocument parses the content of the reader identitied by the filename
func ParseDocument(r io.Reader, config configuration.Configuration, options ...Option) (types.Document, error) {
	fragments, err := ParseDocumentFragments(r, config, options...)
	if err != nil {
		return types.Document{}, err
	}

	elements, err := assemble(fragments)

	// draftDoc, err := ApplySubstitutions(rawDoc, config)
	// if err != nil {
	// 	return types.Document{}, err
	// }
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debug("draft doc:")
	// 	spew.Fdump(log.StandardLogger().Out, draftDoc)
	// }

	// now, merge list items into proper lists
	blocks, err := rearrangeListItems(elements, false)
	if err != nil {
		return types.Document{}, err
	}
	// filter out blocks not needed in the final doc
	blocks = filter(blocks, allMatchers...)

	blocks, footnotes := processFootnotes(blocks)
	// now, rearrange elements in a hierarchical manner
	doc, err := rearrangeSections(blocks)
	if err != nil {
		return types.Document{}, err
	}
	// also, set the footnotes
	doc.Footnotes = footnotes
	// insert the preamble at the right location
	doc = includePreamble(doc)
	// doc.Attributes = doc.Attributes.SetAll(draftDoc.Attributes)
	// also insert the table of contents
	doc = includeTableOfContentsPlaceHolder(doc)
	// finally
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("final document:")
		spew.Fdump(log.StandardLogger().Out, doc)
	}
	return doc, nil
}

// ContextKey a non-built-in type for keys in the context
type ContextKey string

// LevelOffset the key for the level offset of the file to include
const LevelOffset ContextKey = "leveloffset"

// // ParseRawDocument parses a document's content and applies the preprocessing directives (file inclusions)
// func ParseRawDocument(r io.Reader, config configuration.Configuration, options ...Option) (types.RawDocument, error) {
// 	// first, let's find all file inclusions and replace with the actual content to include
// 	source, err := ParseRawSource(r, config, options...)
// 	if err != nil {
// 		return types.RawDocument{}, err
// 	}
// 	if log.IsLevelEnabled(log.DebugLevel) {
// 		log.Debug("source to parse:")
// 		fmt.Fprintf(log.StandardLogger().Out, "'%s'\n", source)
// 	}
// 	// then let's parse the "source" to detect raw blocks
// 	options = append(options, Entrypoint("RawDocument"), GlobalStore(usermacrosKey, config.Macros))
// 	if result, err := Parse(config.Filename, source, options...); err != nil {
// 		return types.RawDocument{}, err
// 	} else if doc, ok := result.(types.RawDocument); ok {
// 		return doc, nil
// 	} else {
// 		return types.RawDocument{}, fmt.Errorf("unexpected type of content: '%T'", result)
// 	}
// }

func assemble(fragments []interface{}) ([]interface{}, error) {
	return nil, nil
}
