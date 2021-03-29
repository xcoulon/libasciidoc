package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// ParseDocumentFragments parses a document's content and applies the preprocessing directives (file inclusions)
// Returns the document fragments (to be assembled) or an error
func ParseDocumentFragments(r io.Reader, config configuration.Configuration, options ...Option) (types.DocumentFragments, error) {
	return parseDocumentFragments(&fileinclusionContext{
		config:       config,
		attributes:   types.Attributes{},
		levelOffsets: []levelOffset{},
		blockLevels:  newStack(),
	}, r, options...)
}

type fileinclusionContext struct {
	config       configuration.Configuration
	attributes   types.Attributes
	levelOffsets []levelOffset
	blockLevels  *stack
}

func (c *fileinclusionContext) clone() *fileinclusionContext {
	return &fileinclusionContext{
		config:       c.config,
		attributes:   c.attributes, // TODO: should we clone this too? ie, can an attribute declared in a child doc be used in rest of the parent doc?
		levelOffsets: append([]levelOffset{}, c.levelOffsets...),
		blockLevels:  c.blockLevels,
	}
}

func (ctx *fileinclusionContext) isSectionRuleEnabled() bool {
	return ctx.blockLevels.empty()
}

func (ctx *fileinclusionContext) onBlockDelimiter(d types.BlockDelimiter) {
	currentLevel := ctx.blockLevels.get()
	if currentLevel == d.Kind {
		ctx.blockLevels.pop() // discard current level, assuming we've just parsed the ending delimiter of the current block
		return
	}
	ctx.blockLevels.push(d.Kind) // push current delimiter kind, assuming we've just parsed the starting delimiter of a new block
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("current substitution context is now: %v", ctx.blockLevels.get())
	}
}

// levelOffset a func that applies a given offset to the sections of a child document to include in a parent doc (the caller)
type levelOffset struct {
	absolute bool
	value    int
	apply    func(*types.Section) int // returns the former level of the section
}

func relativeOffset(offset int) levelOffset {
	return levelOffset{
		absolute: false,
		value:    offset,
		apply: func(s *types.Section) int {
			old := s.Level
			s.Level += offset
			return old
		},
	}
}

func absoluteOffset(offset int) levelOffset {
	return levelOffset{
		absolute: true,
		value:    offset,
		apply: func(s *types.Section) int {
			old := s.Level
			s.Level = offset
			return old
		},
	}
}

func parseDocumentFragments(ctx *fileinclusionContext, r io.Reader, options ...Option) (types.DocumentFragments, error) {
	content, err := ParseReader(ctx.config.Filename, r, append(options, Entrypoint("DocumentFragments"), GlobalStore(substitutionContextKey, ctx))...)
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
			// clone the context so that further level offset in this "child" doc are not applied to the rest of this "parent" doc
			ctx := ctx.clone()
			fileContent, err := parseFileToInclude(ctx, fragment, options...)
			if err != nil {
				return nil, err
			}
			result = append(result, fileContent...)
		case types.AttributeDeclaration:
			// add the attribute to the current content
			ctx.attributes[fragment.Name] = fragment.Value
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

func parseFileToInclude(ctx *fileinclusionContext, incl types.FileInclusion, options ...Option) ([]interface{}, error) {
	incl, err := applySubstitutionsOnFileInclusionPath(substitutionContextDeprecated{
		attributes: types.AttributesWithOverrides{
			Content: ctx.attributes,
		},
		config: ctx.config,
	}, incl)
	if err != nil {
		return nil, err
	}
	path := incl.Location.Stringify()
	currentDir := filepath.Dir(ctx.config.Filename)
	f, absPath, done, err := open(filepath.Join(currentDir, path))
	defer done()
	if err != nil {
		return nil, fmt.Errorf("Unresolved directive in %s - %s", ctx.config.Filename, incl.RawText)
	}
	content := bytes.NewBuffer(nil)
	scanner := bufio.NewScanner(bufio.NewReader(f))
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("parsing file to %s", incl.RawText)
	}
	if lr, ok := lineRanges(incl, ctx.config); ok {
		if err := readWithinLines(scanner, content, lr); err != nil {
			return nil, fmt.Errorf("Unresolved directive in %s - %s", ctx.config.Filename, incl.RawText)
		}
	} else if tr, ok := tagRanges(incl, ctx.config); ok {
		if err := readWithinTags(path, scanner, content, tr); err != nil {
			return nil, err // keep the underlying error here
		}
	} else {
		if err := readAll(scanner, content); err != nil {
			return nil, fmt.Errorf("Unresolved directive in %s - %s", ctx.config.Filename, incl.RawText)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Unresolved directive in %s - %s", ctx.config.Filename, incl.RawText)
	}
	// if the file to include is not an Asciidoc document, just return the content as "raw lines"
	if !IsAsciidoc(absPath) {
		return []interface{}{
			types.RawContent(content.Bytes()),
		}, nil
	}
	// parse the content, and returns the corresponding elements
	lvl, found, err := incl.Attributes.GetAsString(types.AttrLevelOffset)
	if err != nil {
		return nil, err
	}
	if found {
		offset, err := strconv.Atoi(lvl)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read file to include")
		}
		if strings.HasPrefix(lvl, "+") || strings.HasPrefix(lvl, "-") {
			ctx.levelOffsets = append(ctx.levelOffsets, relativeOffset(offset))
		} else {
			ctx.levelOffsets = []levelOffset{absoluteOffset(offset)}

		}
	}
	actualFilename := ctx.config.Filename
	defer func() {
		// restore actual filename after exiting
		ctx.config.Filename = actualFilename
	}()
	ctx.config.Filename = absPath
	// now, let's parse this content and process nested file inclusions (recursively)
	return parseDocumentFragments(ctx, content, options...)
}

// applies the attribute substitutions on the location path of the file to include
func applySubstitutionsOnFileInclusionPath(ctx substitutionContextDeprecated, f types.FileInclusion) (types.FileInclusion, error) {
	elements := [][]interface{}{f.Location.Path} // wrap to
	elements, err := substituteAttributes(ctx, elements)
	if err != nil {
		return types.FileInclusion{}, err
	}
	f.Location.Path = elements[0]
	return f, nil
}

// lineRanges parses the `lines` attribute if it exists in the given FileInclusion, and returns
// a corresponding `LineRanges` (or `false` if parsing failed to invalid input)
func lineRanges(incl types.FileInclusion, config configuration.Configuration) (types.LineRanges, bool) {
	lineRanges, exists, err := incl.Attributes.GetAsString(types.AttrLineRanges)
	if err != nil {
		log.Errorf("Unresolved directive in %s - %s", config.Filename, incl.RawText)
		return types.LineRanges{}, false
	}
	if exists {
		lr, err := Parse("", []byte(lineRanges), Entrypoint("LineRanges"))
		if err != nil {
			log.Errorf("Unresolved directive in %s - %s", config.Filename, incl.RawText)
			return types.LineRanges{}, false
		}
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("line ranges to include:")
			spew.Fdump(log.StandardLogger().Out, lr)
		}
		return types.NewLineRanges(lr), true
	}
	return types.LineRanges{}, false
}

// tagRanges parses the `tags` attribute if it exists in the given FileInclusion, and returns
// a corresponding `TagRanges` (or `false` if parsing failed to invalid input)
func tagRanges(incl types.FileInclusion, config configuration.Configuration) (types.TagRanges, bool) {
	tagRanges, exists, err := incl.Attributes.GetAsString(types.AttrTagRanges)
	if err != nil {
		log.Errorf("Unresolved directive in %s - %s", config.Filename, incl.RawText)
		return types.TagRanges{}, false
	}
	if exists {
		log.Debugf("tag ranges to include: %v", spew.Sdump(tagRanges))
		tr, err := Parse("", []byte(tagRanges), Entrypoint("TagRanges"))
		if err != nil {
			log.Errorf("Unresolved directive in %s - %s", config.Filename, incl.RawText)
			return types.TagRanges{}, false
		}
		return types.NewTagRanges(tr), true
	}
	return types.TagRanges{}, false
}

// TODO: instead of reading and parsing afterwards, simply parse the lines immediatly? ie: `readWithinLines` -> `parseWithinLines`
// (also, use a specific entrypoint if the doc is not a .adoc)
func readWithinLines(scanner *bufio.Scanner, content *bytes.Buffer, lineRanges types.LineRanges) error {
	line := 0
	for scanner.Scan() {
		line++
		// parse the line in search for the `tag::<tag>[]` or `end:<tag>[]` macros
		l, err := Parse("", scanner.Bytes(), Entrypoint("IncludedFileLine"))
		if err != nil {
			return err
		}
		fl, ok := l.(types.IncludedFileLine)
		if !ok {
			return errors.Errorf("unexpected type of parsed line in file to include: %T", l)
		}
		// skip if the line has tags
		if fl.HasTag() {
			continue
		}
		// TODO: stop reading if current line above highest range
		if lineRanges.Match(line) {
			_, err := content.Write(scanner.Bytes())
			if err != nil {
				return err
			}
			_, err = content.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func readWithinTags(path string, scanner *bufio.Scanner, content *bytes.Buffer, expectedRanges types.TagRanges) error {
	// log.Debugf("limiting to tag ranges: %v", expectedRanges)
	currentRanges := make(map[string]*types.CurrentTagRange, len(expectedRanges)) // ensure capacity
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Bytes()
		// parse the line in search for the `tag::<tag>[]` or `end:<tag>[]` macros
		l, err := Parse("", line, Entrypoint("IncludedFileLine"))
		if err != nil {
			return err
		}
		fl, ok := l.(types.IncludedFileLine)
		if !ok {
			return errors.Errorf("unexpected type of parsed line in file to include: %T", l)
		}
		// check if a start or end tag was found in the line
		if startTag, ok := fl.GetStartTag(); ok {
			currentRanges[startTag.Value] = &types.CurrentTagRange{
				StartLine: lineNumber,
				EndLine:   -1,
			}
		}
		if endTag, ok := fl.GetEndTag(); ok {
			currentRanges[endTag.Value].EndLine = lineNumber
		}
		if expectedRanges.Match(lineNumber, currentRanges) && !fl.HasTag() {
			_, err := content.Write(scanner.Bytes())
			if err != nil {
				return err
			}
			_, err = content.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}
	// after the file has been processed, let's check if all tags were "found"
	for _, tag := range expectedRanges {
		// log.Debugf("checking if tag '%s' was found...", tag.Name)
		switch tag.Name {
		case "*", "**":
			continue
		default:
			if tr, found := currentRanges[tag.Name]; !found {
				return fmt.Errorf("tag '%s' not found in include file: %s", tag.Name, path)
			} else if tr.EndLine == -1 {
				log.Warnf("detected unclosed tag '%s' starting at line %d of include file: %s", tag.Name, tr.StartLine, path)
			}
		}
	}
	return nil
}

func readAll(scanner *bufio.Scanner, content *bytes.Buffer) error {
	for scanner.Scan() {
		// parse the line in search for the `tag::<tag>[]` or `end:<tag>[]` macros
		l, err := Parse("", scanner.Bytes(), Entrypoint("IncludedFileLine"))
		if err != nil {
			return err
		}
		fl, ok := l.(types.IncludedFileLine)
		if !ok {
			return errors.Errorf("unexpected type of parsed line in file to include: %T", l)
		}
		// skip if the line has tags
		if fl.HasTag() {
			continue
		}
		_, err = content.Write(scanner.Bytes())
		if err != nil {
			return err
		}
		_, err = content.WriteString("\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func open(path string) (*os.File, string, func(), error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, "", func() {}, err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, "", func() {
			// log.Debugf("restoring current working dir to: %s", wd)
			if err := os.Chdir(wd); err != nil { // restore the previous working directory
				log.WithError(err).Error("failed to restore previous working directory")
			}
		}, err
	}
	dir := filepath.Dir(absPath)
	err = os.Chdir(dir)
	if err != nil {
		return nil, "", func() {
			// log.Debugf("restoring current working dir to: %s", wd)
			if err := os.Chdir(wd); err != nil { // restore the previous working directory
				log.WithError(err).Error("failed to restore previous working directory")
			}
		}, err
	}
	// read the file per-se
	// log.Debugf("opening '%s'", absPath)
	f, err := os.Open(absPath)
	if err != nil {
		return nil, absPath, func() {
			// log.Debugf("restoring current working dir to: %s", wd)
			if err := os.Chdir(wd); err != nil { // restore the previous working directory
				log.WithError(err).Error("failed to restore previous working directory")
			}
		}, err
	}
	return f, absPath, func() {
		// log.Debugf("restoring current working dir to: %s", wd)
		if err := os.Chdir(wd); err != nil { // restore the previous working directory
			log.WithError(err).Error("failed to restore previous working directory")
		}
		if err := f.Close(); err != nil {
			log.WithError(err).Errorf("failed to close file '%s'", absPath)
		}
	}, nil
}

// IsAsciidoc returns true if the file to include is an asciidoc file (based on the file location extension)
func IsAsciidoc(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".asciidoc" || ext == ".adoc" || ext == ".ad" || ext == ".asc" || ext == ".txt"
}
