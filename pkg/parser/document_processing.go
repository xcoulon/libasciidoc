package parser

import (
	"io"

	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ParseDocument parses the content of the reader identitied by the filename
func ParseDocument(r io.Reader, config configuration.Configuration, options ...Option) (types.Document, error) {
	done := make(chan interface{})
	defer close(done)

	attributes := types.Attributes{}
	blocks := []interface{}{}
	inHeader := true
	ctx := newProcessContext()
	pipeline := processFragments(ctx, done, assembleFragments(done, ParseDocumentFragments(r, done, options...)))
	for f := range pipeline {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.WithField("_state", "document_parsing").Debugf("received document fragment: %s", spew.Sdump(f))
		}
		if err := f.Error; err != nil {
			log.WithField("_state", "document_parsing").WithError(err).Error("error occurred")
			return types.Document{}, err
		}
		switch b := f.Content.(type) {
		case types.Section:
			if b.Level != 0 && inHeader { // do not even allow 2ndary section with level 0 as headers
				inHeader = false
			}
		case types.AttributeDeclaration:
			if inHeader {
				attributes.Set(b.Name, b.Value)
			}
		default:
			// anything else and we're not in the header anynore
			inHeader = false
		}
		blocks = append(blocks, f.Content)
	}
	// fragments := make(chan types.DocumentFragment)
	// err := ParseDocumentFragments(r, fragments, options...)
	// if err != nil {
	// 	return types.Document{}, err
	// }

	// elements, err := assemble(fragments)

	// // draftDoc, err := ApplySubstitutions(rawDoc, config)
	// // if err != nil {
	// // 	return types.Document{}, err
	// // }
	// // if log.IsLevelEnabled(log.DebugLevel) {
	// // 	log.Debug("draft doc:")
	// // 	spew.Fdump(log.StandardLogger().Out, draftDoc)
	// // }

	// // now, merge list items into proper lists
	// blocks, err := rearrangeListItems(elements, false)
	// if err != nil {
	// 	return types.Document{}, err
	// }
	// // filter out blocks not needed in the final doc
	// blocks = filter(blocks, allMatchers...)

	// blocks, footnotes := processFootnotes(blocks)
	// // now, rearrange elements in a hierarchical manner
	// doc, err := rearrangeSections(blocks)
	// if err != nil {
	// 	return types.Document{}, err
	// }
	// // also, set the footnotes
	// doc.Footnotes = footnotes
	// // insert the preamble at the right location
	// doc = includePreamble(doc)
	// // doc.Attributes = doc.Attributes.SetAll(draftDoc.Attributes)
	// // also insert the table of contents
	// doc = includeTableOfContentsPlaceHolder(doc)
	// // finally
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debug("final document:")
	// 	spew.Fdump(log.StandardLogger().Out, doc)
	// }
	// return doc, nil
	return types.Document{
		Attributes: attributes.NilIfEmpty(),
		Elements:   blocks,
	}, nil
}

type processContext struct {
	attributes types.Attributes
	counters   map[string]interface{}
}

func newProcessContext() *processContext {
	return &processContext{
		attributes: types.Attributes{},
		counters:   map[string]interface{}{},
	}
}
func (ctx *processContext) addAttribute(name string, value interface{}) {
	ctx.attributes[name] = value
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

func assembleFragments(done <-chan interface{}, fragmentStream <-chan types.DocumentFragmentGroup) <-chan types.DocumentFragment {
	assembledFragmentStream := make(chan types.DocumentFragment)
	go func() {
		defer close(assembledFragmentStream)
		for fragments := range fragmentStream {
			for _, fragment := range assembleFragmentElements(fragments) {
				select {
				case <-done:
					log.WithField("stage", "fragment_assembling").Debug("received 'done' signal")
					return
				case assembledFragmentStream <- fragment:
				}
			}
		}
		log.Debug("end of fragment assembly")
	}()
	return assembledFragmentStream
}

func assembleFragmentElements(f types.DocumentFragmentGroup) []types.DocumentFragment {
	// if the fragment contains an error, then send it as-is downstream
	if err := f.Error; err != nil {
		return []types.DocumentFragment{
			{
				LineOffset: f.LineOffset,
				Error:      f.Error,
			},
		}
	}
	result := make([]types.DocumentFragment, 0, len(f.Content))
	var attributes types.Attributes
	var block types.RawBlock // a block with raw lines
	for i, e := range f.Content {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.WithField("state", "fragment_assembling").Debugf("assembling fragment content of type '%T'", e)
		}
		switch e := e.(type) {
		case types.Attributes:
			if attributes == nil {
				attributes = types.Attributes{}
			}
			attributes.SetAll(e)
		case types.AttributeDeclaration:
			result = append(result, types.NewDocumentFragment(f.LineOffset, e))
		case types.Section:
			e.Attributes = attributes
			attributes = types.Attributes{} // reset
			// send the section block downstream
			result = append(result, types.NewDocumentFragment(f.LineOffset, e))
		case types.BlockDelimiter:
			if block == nil {
				block = types.NewRawDelimitedBlock(e.Kind, attributes)
				attributes = nil // reset
				result = append(result, types.NewDocumentFragment(f.LineOffset, block))
			} else {
				block = nil // reset, in case there is more content afterwards
			}
		case types.RawLine:
			if block == nil {
				block = types.NewRawParagraph(attributes)
				attributes = nil // reset
				result = append(result, types.NewDocumentFragment(f.LineOffset, block))
			}
			block.AddLine(e)
		case types.SingleLineComment:
			result = append(result, types.NewDocumentFragment(f.LineOffset+i, e))
		default:
			// unknow type fragment element: set an error on the fragment and send it downstream
			fr := types.NewDocumentFragment(f.LineOffset, block)
			fr.Error = errors.Errorf("unexpected type of element on line %d: '%T'", f.LineOffset+i, e)
			result = append(result, fr)
		}
	}
	return result
}

// func send(done <-chan interface{}, fragmentStream chan<- types.DocumentFragment, f types.DocumentFragment) {
// 	// if log.IsLevelEnabled(log.DebugLevel) {
// 	// 	log.Debug("sending fragment with following content:")
// 	// 	spew.Fdump(log.StandardLogger().Out, f.Content...)
// 	// }
// 	select {
// 	case <-done:
// 		log.Debug("received 'done' signal")
// 		return
// 	case fragmentStream <- f:
// 	}
// }

func processFragments(ctx *processContext, done <-chan interface{}, fragmentStream <-chan types.DocumentFragment) chan types.DocumentFragment {
	processedFragmentStream := make(chan types.DocumentFragment)
	go func() {
		defer close(processedFragmentStream)
		for f := range fragmentStream {
			select {
			case <-done:
				log.WithField("stage", "fragment_processing").Debug("received 'done' signal")
				return
			case processedFragmentStream <- processFragment(ctx, f):
			}
		}
		log.Debug("end of fragment processing")
	}()
	return processedFragmentStream
}

func processFragment(ctx *processContext, f types.DocumentFragment) types.DocumentFragment {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.WithField("stage", "fragment_processing").Debugf("incoming fragment:\n%s", spew.Sdump(f))
	}
	// if the fragment contains an error, then send it as-is downstream
	if err := f.Error; err != nil {
		return f
	}
	switch e := f.Content.(type) {
	case *types.RawParagraph:
		p, err := applySubstitutionsOnParagraph(ctx, e)
		return types.DocumentFragment{
			LineOffset: f.LineOffset,
			Content:    p,
			Error:      err,
		}
	case types.AttributeDeclaration:
		ctx.addAttribute(e.Name, e.Value)
		return types.NewDocumentFragment(f.LineOffset, e)
	default:
		log.WithField("stage", "fragment_processing").Debugf("forwarding fragment content of type '%T' as-is", e)
		return types.NewDocumentFragment(f.LineOffset, e)
	}
}
