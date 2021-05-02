package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

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

// func parseFileToInclude(ctx *parserContext, incl types.FileInclusion, opts ...Option) ([]interface{}, error) {
// 	incl.Location.Path = substituteAttributes(incl.Location.Path, ctx.attributes)
// 	path := incl.Location.Stringify()
// 	currentDir := filepath.Dir(ctx.config.Filename)
// 	f, absPath, done, err := open(filepath.Join(currentDir, path))
// 	defer done()
// 	if err != nil {
// 		return nil, errors.Errorf("Unresolved directive in %s - %s", ctx.config.Filename, incl.RawText)
// 	}
// 	content := bytes.NewBuffer(nil)
// 	scanner := bufio.NewScanner(bufio.NewReader(f))
// 	if log.IsLevelEnabled(log.DebugLevel) {
// 		log.Debugf("parsing file to %s", incl.RawText)
// 	}
// 	if lr, ok, err := lineRanges(incl, ctx.config); err != nil {
// 		log.WithError(err).Error("error occurred while checking if file inclusion has line ranges")
// 		return nil, errors.Errorf("Unresolved directive in %s - %s", ctx.config.Filename, incl.RawText)
// 	} else if ok {
// 		if err := readWithinLines(scanner, content, lr); err != nil {
// 			log.WithError(err).Error("error occurred while reading file within line ranges")
// 			return nil, errors.Errorf("Unresolved directive in %s - %s", ctx.config.Filename, incl.RawText)
// 		}
// 	} else if tr, ok, err := tagRanges(incl, ctx.config); err != nil {
// 		log.WithError(err).Error("error occurred while checking if file inclusion has tag ranges")
// 		return nil, errors.Errorf("Unresolved directive in %s - %s", ctx.config.Filename, incl.RawText)
// 	} else if ok {
// 		if err := readWithinTags(path, scanner, content, tr); err != nil {
// 			log.WithError(err).Error("error occurred while reading file within tag ranges")
// 			return nil, err // keep the underlying error here
// 		}
// 	} else {
// 		if err := readAll(scanner, content); err != nil {
// 			log.WithError(err).Error("error occurred while reading file")
// 			return nil, errors.Errorf("Unresolved directive in %s - %s", ctx.config.Filename, incl.RawText)
// 		}
// 	}
// 	if err := scanner.Err(); err != nil {
// 		return nil, errors.Errorf("Unresolved directive in %s - %s", ctx.config.Filename, incl.RawText)
// 	}
// 	// if the file to include is not an Asciidoc document, just return the content as "raw lines"
// 	if !IsAsciidoc(absPath) {
// 		return []interface{}{
// 			types.RawLine(content.Bytes()),
// 		}, nil
// 	}
// 	// parse the content, and returns the corresponding elements
// 	lvl, found, err := incl.Attributes.GetAsString(types.AttrLevelOffset)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if found {
// 		offset, err := strconv.Atoi(lvl)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "unable to read file to include")
// 		}
// 		if strings.HasPrefix(lvl, "+") || strings.HasPrefix(lvl, "-") {
// 			ctx.levelOffsets = append(ctx.levelOffsets, relativeOffset(offset))
// 		} else {
// 			ctx.levelOffsets = []levelOffset{absoluteOffset(offset)}

// 		}
// 	}
// 	actualFilename := ctx.config.Filename
// 	defer func() {
// 		// restore actual filename after exiting
// 		ctx.config.Filename = actualFilename
// 	}()
// 	ctx.config.Filename = absPath
// 	// now, let's parse this content and process nested file inclusions (recursively)
// 	return parseDocumentFragments(ctx, content, opts...)
// }

// lineRanges parses the `lines` attribute if it exists in the given FileInclusion, and returns
// a corresponding `LineRanges` (or `false` if parsing failed to invalid input)
func lineRanges(incl types.FileInclusion, config configuration.Configuration) (types.LineRanges, bool, error) {
	lineRanges, exists, err := incl.Attributes.GetAsString(types.AttrLineRanges)
	if err != nil {
		return types.LineRanges{}, false, err
	}
	if exists {
		lr, err := Parse("", []byte(lineRanges), Entrypoint("LineRanges"))
		if err != nil {
			return types.LineRanges{}, false, err
		}
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debugf("line ranges to include: %s", spew.Sdump(lr))
		}
		return types.NewLineRanges(lr), true, nil
	}
	return types.LineRanges{}, false, nil
}

// tagRanges parses the `tags` attribute if it exists in the given FileInclusion, and returns
// a corresponding `TagRanges` (or `false` if parsing failed to invalid input)
func tagRanges(incl types.FileInclusion, config configuration.Configuration) (types.TagRanges, bool, error) {
	tagRanges, exists, err := incl.Attributes.GetAsString(types.AttrTagRanges)
	if err != nil {
		return types.TagRanges{}, false, err
	}
	if exists {
		log.Debugf("tag ranges to include: %v", spew.Sdump(tagRanges))
		tr, err := Parse("", []byte(tagRanges), Entrypoint("TagRanges"))
		if err != nil {
			return types.TagRanges{}, false, err
		}
		return types.NewTagRanges(tr), true, nil
	}
	return types.TagRanges{}, false, nil
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
