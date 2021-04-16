package parser

import (
	"bufio"
	"io"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	log "github.com/sirupsen/logrus"
)

// ParseDocumentFragments parses a document's content and applies the preprocessing directives (file inclusions)
// Returns the document fragments (to be assembled) or an error
func ParseDocumentFragments(r io.Reader, done <-chan interface{}, options ...Option) <-chan types.DocumentFragment {
	return parseDocumentFragments(r, done, options...)
}

// parseDocumentFragments reads the given reader's content line by line, then parses each line using the appropriate
// grammar rule (depending on the context)
func parseDocumentFragments(source io.Reader, done <-chan interface{}, options ...Option) <-chan types.DocumentFragment {
	fragmentStream := make(chan types.DocumentFragment)
	// errStream := make(chan error)
	go func() {
		defer close(fragmentStream)
		// defer close(errStream)
		scanner := NewDocumentFragmentScanner(source, options...)
		for scanner.Scan() {
			select {
			case <-done:
				log.Info("exiting the document parsing routine")
				return // stops/exits the go routine
			case fragmentStream <- scanner.Fragment():
			}
		}
		// if err := scanner.Err(); err != nil {
		// 	errStream <- err
		// }
	}()
	// return fragmentStream, errStream
	return fragmentStream
	// // front-matter
	// frontmatterFragments, done := parseDocumentFrontMatter(ctx, scanner, options...)
	// fragments = append(fragments, frontmatterFragments...)
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debug("front-matter:")
	// 	spew.Fdump(log.StandardLogger().Out, fragments)
	// }
	// if done {
	// 	return fragments, nil
	// }
	// // document header
	// headerFragments, done := parseDocumentHeader(ctx, scanner, options...)
	// fragments = append(fragments, headerFragments...)
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debug("document header:")
	// 	spew.Fdump(log.StandardLogger().Out, fragments)
	// }
	// if done {
	// 	return fragments, nil
	// }
	// // document body
	// bodyFragments, err := parseDocumentBody(ctx, scanner, options...)
	// if err != nil {
	// 	return nil, err
	// }
	// fragments = append(fragments, bodyFragments...)

	// make sure the scanner reached the end of the source
	// return scanner.Err()
}

type DocumentFragmentScanner struct {
	options    []Option
	scanner    *bufio.Scanner
	lineNumber int
	err        error // sticky error
	fragment   types.DocumentFragment
}

func NewDocumentFragmentScanner(source io.Reader, options ...Option) *DocumentFragmentScanner {
	return &DocumentFragmentScanner{
		scanner: bufio.NewScanner(source),
		options: append(options, Entrypoint("DocumentFragmentElement")),
	}
}

// Scan retrieves the next document fragment
// (ie, a group of lines with essentially raw lines and optionally block delimiters)
// and returns `true` if a block was found, `false` if the end of the doc was reached,
// or if an error occurred (see `DocumentFragmentScanner.Err()`)
func (s *DocumentFragmentScanner) Scan() bool {
	if s.err != nil {
		// error was found in the previous call, so let's stop now.
		return false
	}
	elements := []interface{}{}
	s.fragment = types.DocumentFragment{
		LineOffset: s.lineNumber + 1, // next fragment will begin at the next line read by the underlying scanner
	}
	blockDelimiters := newStack()
scan:
	for s.scanner.Scan() {
		s.lineNumber++
		element, err := Parse("", s.scanner.Bytes(), s.options...)
		if err != nil {
			log.WithError(err).Error("failed to parse the content")
			s.err = err // will cause next call to `Scan()` to return false
			s.fragment.Error = err
			return true
		}
		switch element := element.(type) {
		case types.BlankLine:
			// blanklines outside of a delimited block causes the scanner to stop (for this call)
			if blockDelimiters.empty() {
				break scan
			}
			// if inside a delimited blocl, retain the blankline
			elements = append(elements, types.RawLine(""))
		case types.BlockDelimiter:
			if blockDelimiters.get() == element { // match starting/ending delimiters
				blockDelimiters.pop() // remove element from the top of the stack
			} else {
				blockDelimiters.push(element)
			}
			elements = append(elements, element)
		default:
			elements = append(elements, element)
		}
	}
	if err := s.scanner.Err(); err != nil {
		log.WithError(err).Error("failed to read the content")
		s.err = err // will cause next call to `Scan()` to return false
		s.fragment.Error = err
		return true // this fragment needs to be processed upstream
	}
	if len(elements) > 0 {
		s.fragment.Content = elements
		return true // also return the underlying scanner's error (if something wrong happened)
	}
	// reached the end of source
	return false // also return the underlying scanner's error (if something wrong happened)
}

// // Err returns the error that way have occurred when the underlying scanner
// // was reading the source, or when parsing the content to find the next block
// func (s *DocumentFragmentScanner) Err() error {
// 	if s.err != nil {
// 		return s.err
// 	}
// 	return s.scanner.Err()
// }

// Fragment returns document fragment that was read by the last call to `Next`.
// Multiple calls will return the same value, until `Next` is called again
func (s *DocumentFragmentScanner) Fragment() types.DocumentFragment {
	return s.fragment
}

// // parseDocumentFrontMatter attempts to read the front-matter if it exists.
// // scans line by line, exit after the front-matter delimiter (if applicable)
// // return the document fragments along with a bool flag to indicate if the scanner reached the end of the document
// func parseDocumentFrontMatter(ctx *parserContext, scanner *bufio.Scanner, options ...Option) (types.DocumentFragments, bool) {
// 	log.Debug("parsing front-matter...")
// 	fragments := make([]interface{}, 0)
// 	options = append(options, Entrypoint("FrontMatterFragment"))
// 	withinBlock := false
// 	for scanner.Scan() {
// 		fragment, err := Parse("", scanner.Bytes(), options...)
// 		if err != nil { // no match
// 			return fragments, false
// 		}
// 		switch f := fragment.(type) {
// 		case types.BlockDelimiter:
// 			fragments = append(fragments, f)
// 			if withinBlock {
// 				// it's time to exit the parsing of the front-matter
// 				return fragments, false
// 			}
// 			withinBlock = true
// 		default:
// 			if !withinBlock {
// 				// unexpected content
// 				return fragments, false
// 			}
// 			fragments = append(fragments, fragment)
// 		}
// 	}
// 	return fragments, true
// }

// // parseDocumentHeader attempts to read the document header (title and metadata) if it exists.
// // scans line by line, exit after a blankline is found (if applicable)
// // return the document fragments along with a bool flag to indicate if the scanner reached the end of the document
// func parseDocumentHeader(ctx *parserContext, scanner *bufio.Scanner, options ...Option) (types.DocumentFragments, bool) {
// 	log.Debug("parsing document header...")
// 	fragments := make([]interface{}, 0)
// 	// check if there is a title
// 	options = append(options, Entrypoint("DocumentTitle"))
// 	title, found := doParse(ctx, scanner.Bytes(), options...)
// 	if !found {
// 		// if there is no title, then there is no header at all
// 		return fragments, false
// 	}
// 	fragments = append(fragments, title)
// 	// check if there are authors
// 	options = append(options, Entrypoint("DocumentAuthorsMetadata"))
// 	for scanner.Scan() {
// 		fragment, found := doParse(ctx, scanner.Bytes(), options...)
// 		if found { // no match
// 			fragments = append(fragments, fragment)
// 		}
// 	}
// 	// check if there is a revision
// 	options = append(options, Entrypoint("DocumentRevisionMetadata"))
// 	for scanner.Scan() {
// 		fragment, found := doParse(ctx, scanner.Bytes(), options...)
// 		if !found { // no match
// 			return fragments, false
// 		}
// 		fragments = append(fragments, fragment)
// 	}
// 	return fragments, true
// }

// func doParse(ctx *parserContext, content []byte, options ...Option) (interface{}, bool) {
// 	fragment, err := Parse("", content, options...)
// 	if err != nil {
// 		// no match
// 		return nil, false
// 	}
// 	return fragment, true
// }

// // parseDocumentBody attempts to read the document body if it exists.
// func parseDocumentBody(ctx *parserContext, scanner *bufio.Scanner, options ...Option) (types.DocumentFragments, error) {
// 	log.Debug("parsing document body...")
// 	fragments := make([]interface{}, 0)
// 	options = append(options, Entrypoint("DocumentBodyFragment"))
// 	fragment, err := doParseDocumentBody(ctx, scanner.Bytes(), options...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fragments = append(fragments, fragment)
// 	for scanner.Scan() {
// 		fragment, err := doParseDocumentBody(ctx, scanner.Bytes(), options...)
// 		if err != nil {
// 			return nil, err
// 		}
// 		switch f := fragment.(type) {
// 		case []interface{}:
// 			// happens when a `FileInclusion` element was replaced by the target file's content (ie, fragments)
// 			fragments = append(fragments, f...)
// 		default:
// 			fragments = append(fragments, fragment)
// 		}
// 	}
// 	return fragments, nil
// }

// func doParseDocumentBody(ctx *parserContext, content []byte, options ...Option) (interface{}, error) {
// 	fragment, err := Parse("", content, options...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	switch fragment := fragment.(type) {
// 	case types.FileInclusion:
// 		return parseFileToInclude(ctx.clone(), fragment, options...) // clone the context so that further level offset in this "child" doc are not applied to the rest of this "parent" doc
// 	case types.AttributeDeclaration:
// 		// immediatly process the attribute substitutions in the value (if there is any)
// 		if err := ctx.onAttributeDeclaration(fragment); err != nil {
// 			return nil, err
// 		}
// 		return fragment, nil
// 	case types.Section:
// 		// apply the level offsets on the sections of the doc to include (by default, there is no offset)
// 		oldLevel := fragment.Level
// 		for _, offset := range ctx.levelOffsets {
// 			oldLevel := offset.apply(&fragment)
// 			// replace the absolute level offset with a relative one,
// 			// taking into account the offset that was applied on this section,
// 			// so that following sections in the document have the same relative offset
// 			// applied, but keep the same hierarchy
// 			if offset.absolute {
// 				// replace current offset and all previous ones with a single relative offset
// 				ctx.levelOffsets = []levelOffset{relativeOffset(fragment.Level - oldLevel)}
// 			}
// 		}
// 		if log.IsLevelEnabled(log.DebugLevel) {
// 			log.Debugf("applied level offset on section: level %d -> %d", oldLevel, fragment.Level)
// 			// log.Debug("level offsets:")
// 			// spew.Fdump(log.StandardLogger().Out, ctx.levelOffsets)
// 		}
// 		return fragment, nil
// 	case types.BlockDelimiter:
// 		ctx.onBlockDelimiter(fragment)
// 		return fragment, nil
// 	case types.InlineElements:
// 		return onInlineElements(ctx, fragment)
// 	default:
// 		return fragment, nil
// 	}
// }

// func parseDocumentFragments(ctx *parserContext, r io.Reader, options ...Option) (types.DocumentFragments, error) {
// 	content, err := ParseReader(ctx.config.Filename, r, append(options, Entrypoint("DocumentFragments"), GlobalStore(parseContextKey, ctx))...)
// 	if err != nil {
// 		log.Errorf("failed to parse raw document: %s", err)
// 		return nil, err
// 	}
// 	fragments, ok := content.(types.DocumentFragments)
// 	if !ok {
// 		return nil, fmt.Errorf("unexpected type of raw lines: '%T'", content)
// 	}
// 	// look-up file inclusions and replace with the file content
// 	result := make([]interface{}, 0, len(fragments)) // assume the initial capacity as if there is no file inclusion to process
// 	for _, fragment := range fragments {
// 		switch fragment := fragment.(type) {
// 		case types.FileInclusion:
// 			fileContent, err := parseFileToInclude(ctx.clone(), fragment, options...) // clone the context so that further level offset in this "child" doc are not applied to the rest of this "parent" doc
// 			if err != nil {
// 				return nil, err
// 			}
// 			result = append(result, fileContent...)
// 		case types.AttributeDeclaration:
// 			// immediatly process the attribute substitutions in the value (if there is any)
// 			value := substituteAttributes(fragment.Value, ctx.attributes)
// 			switch value := value.(type) {
// 			case types.StringElement:
// 				ctx.attributes[fragment.Name] = value.Content
// 			case string:
// 				ctx.attributes[fragment.Name] = value
// 			default:
// 				return nil, fmt.Errorf("unexpected type of value after substituing attributes: '%T'", value)
// 			}
// 			// add the attribute to the current content
// 			result = append(result, fragment)

// 		case types.Section:
// 			// apply the level offsets on the sections of the doc to include (by default, there is no offset)
// 			oldLevel := fragment.Level
// 			for _, offset := range ctx.levelOffsets {
// 				oldLevel := offset.apply(&fragment)
// 				// replace the absolute level offset with a relative one,
// 				// taking into account the offset that was applied on this section,
// 				// so that following sections in the document have the same relative offset
// 				// applied, but keep the same hierarchy
// 				if offset.absolute {
// 					// replace current offset and all previous ones with a single relative offset
// 					ctx.levelOffsets = []levelOffset{relativeOffset(fragment.Level - oldLevel)}
// 				}
// 			}
// 			if log.IsLevelEnabled(log.DebugLevel) {
// 				log.Debugf("applied level offset on section: level %d -> %d", oldLevel, fragment.Level)
// 				// log.Debug("level offsets:")
// 				// spew.Fdump(log.StandardLogger().Out, ctx.levelOffsets)
// 			}
// 			result = append(result, fragment)
// 		case types.BlockDelimiter:
// 			ctx.onBlockDelimiter(fragment)
// 			result = append(result, fragment)
// 		case types.InlineElements:
// 			fragment, err := onInlineElements(ctx, fragment)
// 			if err != nil {
// 				return nil, err
// 			}
// 			result = append(result, fragment)
// 		default:
// 			result = append(result, fragment)
// 		}
// 	}
// 	// if log.IsLevelEnabled(log.DebugLevel) {
// 	// 	log.Debug("parsed file to include")
// 	// 	spew.Fdump(log.StandardLogger().Out, result)
// 	// }
// 	return result, nil
// }
