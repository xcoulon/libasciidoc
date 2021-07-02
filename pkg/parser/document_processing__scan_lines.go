package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

// ScanDocument scans a document's content and applies the preprocessing directives (file inclusions)
// Returns the document fragments (to be assembled) or an error
func ScanDocument(source io.Reader, done <-chan interface{}, opts ...Option) <-chan types.DocumentFragment {
	fragmentGroupStream := make(chan types.DocumentFragment)
	// errStream := make(chan error)
	go func() {
		defer close(fragmentGroupStream)
		// defer close(errStream)
		scanner := NewDocumentScanner(source, opts...)
		for scanner.Scan() {
			select {
			case <-done:
				log.Info("exiting the document parsing routine")
				return // stops/exits the go routine
			case fragmentGroupStream <- scanner.Fragments():
			}
		}
		log.WithField("pipeline_task", "fragment_parsing").Debug("end of fragment parsing")
	}()
	return fragmentGroupStream
}

type DocumentFragmentScanner struct {
	opts         []Option
	scanner      *bufio.Scanner
	lineNumber   int
	currentGroup types.DocumentFragment
	scopes       *scanScopeStack
	err          error // sticky error
}

type scanScope struct {
	name string
	kind string
}

func (s scanScope) String() string {
	result := &strings.Builder{}
	result.WriteString(s.name)
	if s.kind != "" {
		result.WriteString("/")
		result.WriteString(s.kind)
	}
	return result.String()
}

var (
	defaultScope       = &scanScope{name: "default"}
	withinParagraph    = &scanScope{name: "within_paragraph"}
	withinListingBlock = &scanScope{name: "within_delimited_block", kind: types.Listing}
	withinExampleBlock = &scanScope{name: "within_delimited_block", kind: types.Example}
	withinFencedBlock  = &scanScope{name: "within_delimited_block", kind: types.Fenced}
	withinLiteralBlock = &scanScope{name: "within_delimited_block", kind: types.Literal}
)

func NewDocumentScanner(source io.Reader, opts ...Option) *DocumentFragmentScanner {
	opts = append([]Option{Entrypoint("DefaultDocumentFragmentElement")}, opts...) // set the default entrypoint to parse the content
	return &DocumentFragmentScanner{
		scanner: bufio.NewScanner(source),
		opts:    opts,
		scopes:  newScanningScopeStack(),
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
	s.currentGroup = types.DocumentFragment{
		LineOffset: s.lineNumber + 1, // next fragment will begin at the next line read by the underlying scanner
	}
scan:
	for s.scanner.Scan() {
		s.lineNumber++
		element, err := Parse("", s.scanner.Bytes(), append(s.opts, s.entrypoint()...)...) // TODO: only parse for Blankline and SingleLineComments if within a paragraph
		if err != nil {
			// assume that the content is just a RawLine
			element = types.RawLine(s.scanner.Text())
		}
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debugf("parsed '%s':\n%s", s.scanner.Text(), spew.Sdump(element))
			if err != nil {
				log.Debugf("err=%v", err)
			}
		}
		switch element := element.(type) {
		case *types.BlankLine:
			elements = append(elements, element)
			// blanklines outside of a delimited block causes the scanner to stop (for this call)
			if s.scopes.get() == withinParagraph || s.scopes.get() == defaultScope {
				s.scopes.pop()
				break scan // end of fragment
			}
		case *types.OrderedListElement, *types.UnorderedListElement, *types.LabeledListElement, *types.CalloutListElement:
			elements = append(elements, element)
		case *types.BlockDelimiter:
			currentScope := s.scopes.get()
			if currentScope.kind == element.Kind {
				// done with current delimited block
				s.scopes.pop()
				if s.scopes.empty() {
					break scan // end of fragment
				}
			} else {
				switch element.Kind {
				case types.Listing:
					s.scopes.push(withinListingBlock)
				case types.Example:
					s.scopes.push(withinExampleBlock)
				case types.Fenced:
					s.scopes.push(withinFencedBlock)
				case types.Literal:
					s.scopes.push(withinLiteralBlock)
				default:
					s.err = fmt.Errorf("unsupported kind of block delimiter: %v", element.Kind)
					break scan // end of fragment group
				}
			}
			// log.Debugf("updated scanner scope: %s", s.scopes.get())
			elements = append(elements, element)
		case types.RawLine:
			if s.scopes.get() == defaultScope {
				// we're now within a paragraph
				s.scopes.push(withinParagraph)
			}
			elements = append(elements, element)
		case *types.ListElementContinuation:
			// s.scopes.push(withinListContinuation)
			elements = append(elements, element)
		default:
			elements = append(elements, element)
		}
	}
	if err := s.scanner.Err(); err != nil {
		log.WithError(err).Error("failed to read the content")
		s.err = err                   // will cause next call to `Scan()` to return false
		s.currentGroup.Elements = nil // reset
		s.currentGroup.Error = err
		return true // this fragment needs to be processed upstream
	}
	if len(elements) > 0 {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debugf("parsed fragment elements:\n%s", spew.Sdump(elements))
		}

		s.currentGroup.Elements = elements
		return true // (and will also return the underlying scanner's error if something wrong happened)
	}
	// reached the end of source
	return false // (and will also return the underlying scanner's error if something wrong happened)
}

func (c *current) isValidBlockDelimiter(d *types.BlockDelimiter) (bool, error) {
	if k, ok := c.globalStore[validDelimitedBlockKind].(string); ok {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debugf("checking if delimiter matches kind '%v': %s", k, spew.Sdump(d))
		}
		return d.Kind == k, nil
	}
	log.Debug("no valid delimiter registered in the GlobalStore")
	return false, nil
}

const validDelimitedBlockKind = "valid_delimited_block_kind"

// returns the EntryPoint and GlobalStore if within a delimited block,
// so we know which delimiter can be accepted
func (s *DocumentFragmentScanner) entrypoint() []Option {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("current parsing scope: %v", s.scopes.get())
	}
	currentScope := s.scopes.get()
	switch currentScope {
	case withinParagraph:
		log.Debugf("using '%s' entrypoint", "DocumentFragmentElementWithinParagraph")
		return []Option{
			Entrypoint("DocumentFragmentElementWithinParagraph"),
		}
	case withinListingBlock, withinFencedBlock, withinLiteralBlock:
		log.Debugf("using '%s' entrypoint", "DocumentFragmentElementWithinVerbatimBlock")
		return []Option{
			Entrypoint("DocumentFragmentElementWithinVerbatimBlock"),
			GlobalStore(validDelimitedBlockKind, currentScope.kind),
		}
	case withinExampleBlock:
		log.Debugf("using '%s' entrypoint", "DocumentFragmentElementWithinNormalBlock")
		return []Option{
			Entrypoint("DocumentFragmentElementWithinNormalBlock"),
			GlobalStore(validDelimitedBlockKind, currentScope.kind),
		}
	default:
		log.Debugf("using default entrypoint")
		return nil
	}
}

// Fragments returns document fragments that were read by the last call to `Next`.
// Multiple calls will return the same value, until `Next` is called again
func (s *DocumentFragmentScanner) Fragments() types.DocumentFragment {
	log.Debugf("returning fragment group with %d elements", len(s.currentGroup.Elements))
	return s.currentGroup
}

type scanScopeStack struct {
	scopes []*scanScope
}

func newScanningScopeStack() *scanScopeStack {
	scopes := make([]*scanScope, 0, 10)
	scopes = append(scopes, defaultScope)
	return &scanScopeStack{
		scopes: scopes,
	}
}

// func (s *scanningScopes) size() int {
// 	return s.index + 1
// }

// func (s *scanningScopes) empty() bool {
// 	return s.index == -1
// }

func (s *scanScopeStack) push(a *scanScope) {
	s.scopes = append(s.scopes, a)
}

func (s *scanScopeStack) pop() *scanScope {
	if len(s.scopes) <= 0 {
		return defaultScope
	}
	scope := s.scopes[len(s.scopes)-1]
	s.scopes = s.scopes[:len(s.scopes)-1] // shrink
	return scope
}

func (s *scanScopeStack) get() *scanScope {
	if len(s.scopes) <= 0 {
		return defaultScope
	}
	return s.scopes[len(s.scopes)-1]
}

func (s *scanScopeStack) empty() bool {
	return len(s.scopes) == 0
}

// // parseDocumentFrontMatter attempts to read the front-matter if it exists.
// // scans line by line, exit after the front-matter delimiter (if applicable)
// // return the document fragments along with a bool flag to indicate if the scanner reached the end of the document
// func parseDocumentFrontMatter(ctx *parserContext, scanner *bufio.Scanner, opts ...Option) (types.DocumentFragments, bool) {
// 	log.WithField("pipeline_task", "fragment_parsing").Debug("parsing front-matter...")
// 	fragments := make([]interface{}, 0)
// 	opts = append(opts, Entrypoint("FrontMatterFragment"))
// 	withinBlock := false
// 	for scanner.Scan() {
// 		fragment, err := Parse("", scanner.Bytes(), opts...)
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
// func parseDocumentHeader(ctx *parserContext, scanner *bufio.Scanner, opts ...Option) (types.DocumentFragments, bool) {
// 	log.WithField("pipeline_task", "fragment_parsing").Debug("parsing document header...")
// 	fragments := make([]interface{}, 0)
// 	// check if there is a title
// 	opts = append(opts, Entrypoint("DocumentTitle"))
// 	title, found := doParse(ctx, scanner.Bytes(), opts...)
// 	if !found {
// 		// if there is no title, then there is no header at all
// 		return fragments, false
// 	}
// 	fragments = append(fragments, title)
// 	// check if there are authors
// 	opts = append(opts, Entrypoint("DocumentAuthorsMetadata"))
// 	for scanner.Scan() {
// 		fragment, found := doParse(ctx, scanner.Bytes(), opts...)
// 		if found { // no match
// 			fragments = append(fragments, fragment)
// 		}
// 	}
// 	// check if there is a revision
// 	opts = append(opts, Entrypoint("DocumentRevisionMetadata"))
// 	for scanner.Scan() {
// 		fragment, found := doParse(ctx, scanner.Bytes(), opts...)
// 		if !found { // no match
// 			return fragments, false
// 		}
// 		fragments = append(fragments, fragment)
// 	}
// 	return fragments, true
// }

// func doParse(ctx *parserContext, content []byte, opts ...Option) (interface{}, bool) {
// 	fragment, err := Parse("", content, opts...)
// 	if err != nil {
// 		// no match
// 		return nil, false
// 	}
// 	return fragment, true
// }

// // parseDocumentBody attempts to read the document body if it exists.
// func parseDocumentBody(ctx *parserContext, scanner *bufio.Scanner, opts ...Option) (types.DocumentFragments, error) {
// 	log.WithField("pipeline_task", "fragment_parsing").Debug("parsing document body...")
// 	fragments := make([]interface{}, 0)
// 	opts = append(opts, Entrypoint("DocumentBodyFragment"))
// 	fragment, err := doParseDocumentBody(ctx, scanner.Bytes(), opts...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fragments = append(fragments, fragment)
// 	for scanner.Scan() {
// 		fragment, err := doParseDocumentBody(ctx, scanner.Bytes(), opts...)
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

// func doParseDocumentBody(ctx *parserContext, content []byte, opts ...Option) (interface{}, error) {
// 	fragment, err := Parse("", content, opts...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	switch fragment := fragment.(type) {
// 	case types.FileInclusion:
// 		return parseFileToInclude(ctx.clone(), fragment, opts...) // clone the context so that further level offset in this "child" doc are not applied to the rest of this "parent" doc
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
// 		if log.IsLevelEnabled(log.WithField("pipeline_task", "fragment_parsing").DebugLevel) {
// 			log.WithField("pipeline_task", "fragment_parsing").Debugf("applied level offset on section: level %d -> %d", oldLevel, fragment.Level)
// 			// log.WithField("pipeline_task", "fragment_parsing").Debug("level offsets:")
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

// func parseDocumentFragments(ctx *parserContext, r io.Reader, opts ...Option) (types.DocumentFragments, error) {
// 	content, err := ParseReader(ctx.config.Filename, r, append(opts, Entrypoint("DocumentFragments"), GlobalStore(parseContextKey, ctx))...)
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
// 			fileContent, err := parseFileToInclude(ctx.clone(), fragment, opts...) // clone the context so that further level offset in this "child" doc are not applied to the rest of this "parent" doc
// 			if err != nil {
// 				return nil, err
// 			}
// 			result = append(result, fileContent...)
// 		case types.AttributeDeclaration:
// 			// immediatly process the attribute substitutions in the value (if there is any)
// 			value := substituteAttributes(fragment.Value, ctx.attributes)
// 			switch value := value.(type) {
// 			case *types.StringElement:
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
// 			if log.IsLevelEnabled(log.WithField("pipeline_task", "fragment_parsing").DebugLevel) {
// 				log.WithField("pipeline_task", "fragment_parsing").Debugf("applied level offset on section: level %d -> %d", oldLevel, fragment.Level)
// 				// log.WithField("pipeline_task", "fragment_parsing").Debug("level offsets:")
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
// 	// if log.IsLevelEnabled(log.WithField("pipeline_task", "fragment_parsing").DebugLevel) {
// 	// 	log.WithField("pipeline_task", "fragment_parsing").Debug("parsed file to include")
// 	// 	spew.Fdump(log.StandardLogger().Out, result)
// 	// }
// 	return result, nil
// }
