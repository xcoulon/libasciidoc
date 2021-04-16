package parser

import (
	"bytes"
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

	blocks := []interface{}{}
	ctx := newProcessContext()
	pipeline := processFragments(ctx, done, assembleFragments(done, ParseDocumentFragments(r, done, options...)))
	for f := range pipeline {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debugf("received document fragment: %s", spew.Sdump(f))
		}

		if err := f.Error; err != nil {
			done <- struct{}{} // halt routines
			return types.Document{}, err
		}
		blocks = append(blocks, f.Content...)
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
		Elements: blocks,
	}, nil
}

type processContext struct {
	attributes types.Attributes
}

func newProcessContext() *processContext {
	return &processContext{
		attributes: types.Attributes{},
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

func assembleFragments(done <-chan interface{}, fragmentStream <-chan types.DocumentFragment) <-chan types.DocumentFragment {
	assembledFragmentStream := make(chan types.DocumentFragment)
	go func() {
		defer close(assembledFragmentStream)
		for f := range fragmentStream {
			// process the fragment (ie, set attributes and perform paragraph and block substitutions, etc.)
			assembleFragmentElements(done, assembledFragmentStream, f) // TODO: pass `assembledFragmentStream` as arg
		}
	}()
	return assembledFragmentStream
}

func assembleFragmentElements(done <-chan interface{}, assembledFragmentStream chan<- types.DocumentFragment, f types.DocumentFragment) {
	// if the fragment contains an error, then send it as-is downstream
	if err := f.Error; err != nil {
		assembledFragmentStream <- f
		return
	}
	var attributes types.Attributes
	var block types.RawBlock // a block with raw lines
	for i, e := range f.Content {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debugf("fragment content of type '%T'", e)
		}
		switch e := e.(type) {
		case types.Attributes:
			if attributes == nil {
				attributes = types.Attributes{}
			}
			attributes.SetAll(e)
		case types.AttributeDeclaration:
			send(done, assembledFragmentStream, types.NewDocumentFragment(f.LineOffset, e))
		case types.Section:
			e.Attributes = attributes
			attributes = types.Attributes{} // reset
			// send the section block downstream
			send(done, assembledFragmentStream, types.NewDocumentFragment(f.LineOffset, e))
		case types.BlockDelimiter:
			if block == nil {
				block = types.NewRawDelimitedBlock(e.Kind, attributes)
				attributes = nil // reset
				continue
			}
			send(done, assembledFragmentStream, types.NewDocumentFragment(f.LineOffset, block))
			block = nil // reset, in case there is more content afterwards
		case types.RawLine:
			if block == nil {
				block = types.NewRawParagraph(attributes)
				attributes = nil // reset
			}
			block = block.AddLine(e)
		default:
			// unknow type fragment element: set an error on the fragment and send it downstream
			f.Error = errors.Errorf("unexpected type of element on line %d: '%T'", f.LineOffset+i, e)
			assembledFragmentStream <- f
		}
	}
	// end of iteration on
	if block != nil {
		send(done, assembledFragmentStream, types.NewDocumentFragment(f.LineOffset, block))
	}
}

func send(done <-chan interface{}, fragmentStream chan<- types.DocumentFragment, f types.DocumentFragment) {
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debug("sending fragment with following content:")
	// 	spew.Fdump(log.StandardLogger().Out, f.Content...)
	// }
	select {
	case <-done:
		log.Debug("received 'done' signal")
		return
	case fragmentStream <- f:
	}
}

func processFragments(ctx *processContext, done <-chan interface{}, fragmentStream <-chan types.DocumentFragment) chan types.DocumentFragment {
	processedFragmentStream := make(chan types.DocumentFragment)
	go func() {
		defer close(processedFragmentStream)
		for f := range fragmentStream {
			processFragment(ctx, done, processedFragmentStream, f)
		}

	}()
	return processedFragmentStream
}

func processFragment(ctx *processContext, done <-chan interface{}, processedFragmentStream chan<- types.DocumentFragment, f types.DocumentFragment) {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("processing fragment")
		spew.Fdump(log.StandardLogger().Out, f)
	}
	// if the fragment contains an error, then send it as-is downstream
	if err := f.Error; err != nil {
		processedFragmentStream <- f
		return
	}
	for _, e := range f.Content {
		switch e := e.(type) {
		case types.RawParagraph:
			// apply the attribute substitutions in the attributes
			var err error
			if e.Attributes, err = applyAttributeSubstitutionsOnAttributes(ctx, e.Attributes); err != nil {
				f.Error = err
				processedFragmentStream <- f
			}
			// parse the lines, using the default substitutions
			lines, err := Parse("", serializeLines(e.Lines), Entrypoint("NormalSubstitution"), GlobalStore(substitutionContextKey, newSubstitutionContext()))
			if err != nil {
				f.Error = err
				processedFragmentStream <- f
			}
			if lines, ok := lines.([]interface{}); ok {
				// if log.IsLevelEnabled(log.DebugLevel) {
				// 	log.Debug("paragraph lines")
				// 	spew.Fdump(log.StandardLogger().Out, lines)
				// }
				p, err := types.NewParagraph(lines, e.Attributes)
				if err != nil {
					processedFragmentStream <- types.DocumentFragment{ // TODO: constructor for "Error Fragment"?
						LineOffset: f.LineOffset,
						Error:      err,
					}
				}
				if log.IsLevelEnabled(log.DebugLevel) {
					log.Printf("new paragraph: %s", spew.Sdump(p))
				}
				processedFragmentStream <- types.NewDocumentFragment(f.LineOffset, p)
			}
		case types.AttributeDeclaration:
			ctx.addAttribute(e.Name, e.Value)
			processedFragmentStream <- types.NewDocumentFragment(f.LineOffset, e)
		default:
			log.Debugf("forwarding fragment content of type '%T' as-is", e)
			processedFragmentStream <- types.NewDocumentFragment(f.LineOffset, e)
		}
	}
}

// ----------------------------------------------------------------------------
// Attribute substitutions
// ----------------------------------------------------------------------------

func applyAttributeSubstitutionsOnElements(ctx *processContext, elements []interface{}) ([]interface{}, error) {
	result := make([]interface{}, len(elements)) // maximum capacity should exceed initial input
	for i, element := range elements {
		e, err := applyAttributeSubstitutionsOnElement(ctx, element)
		if err != nil {
			return nil, err
		}
		result[i] = e
	}
	return result, nil
}

func applyAttributeSubstitutionsOnAttributes(ctx *processContext, attributes types.Attributes) (types.Attributes, error) {
	for key, value := range attributes {
		switch key {
		case types.AttrRoles, types.AttrOptions: // multi-value attributes
			result := []interface{}{}
			if values, ok := value.([]interface{}); ok {
				for _, value := range values {
					switch value := value.(type) {
					case []interface{}:
						value, err := applyAttributeSubstitutionsOnElements(ctx, value)
						if err != nil {
							return nil, err
						}
						result = append(result, types.Reduce(value))
					default:
						result = append(result, value)
					}

				}
				attributes[key] = result
			}
		default: // single-value attributes
			if value, ok := value.([]interface{}); ok {
				value, err := applyAttributeSubstitutionsOnElements(ctx, value)
				if err != nil {
					return nil, err
				}
				attributes[key] = types.Reduce(value)
			}
		}
	}
	return attributes, nil
}

func applyAttributeSubstitutionsOnElement(ctx *processContext, element interface{}) (interface{}, error) {
	switch e := element.(type) {
	case types.AttributeDeclaration:
		ctx.attributes.Set(e.Name, e.Value)
		return e, nil
	case types.AttributeReset:
		delete(ctx.attributes, e.Name)
		return e, nil
	case types.AttributeSubstitution:
		return types.StringElement{
			Content: ctx.attributes.GetAsStringWithDefault(e.Name, "{"+e.Name+"}"),
		}, nil
	default:
		return e, nil
	}
	// case types.CounterSubstitution:
	// 	if element, err = applyCounterSubstitution(ctx, e); err != nil {
	// 		return nil, err
	// 	}
	// case types.WithElementsToSubstitute:
	// 	elmts, err := applyAttributeSubstitutionsOnElements(ctx, e.ElementsToSubstitute())
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	element = e.ReplaceElements(types.Merge(elmts))
	// case types.WithLineSubstitution:
	// 	lines, err := applyAttributeSubstitutionsOnLines(ctx, e.LinesToSubstitute())
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	element = e.SubstituteLines(lines)
	// case types.ContinuedListItemElement:
	// 	if e.Element, err = applyAttributeSubstitutionsOnElement(ctx, e.Element); err != nil {
	// 		return nil, err
	// 	}
	// }
	// also, retain the attribute declaration value (if applicable)
}

func serializeLines(lines []types.RawLine) []byte {
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debug("serializing lines")
	// 	spew.Fdump(log.StandardLogger().Out, lines)
	// }
	buf := bytes.Buffer{}
	for i, l := range lines {
		buf.WriteString(string(l))
		if i < len(lines)-1 {
			buf.WriteString("\n")
		}
	}
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debug("serialized lines")
	// 	spew.Fdump(log.StandardLogger().Out, buf.String())
	// }
	return buf.Bytes()
}

// func assemble(fragments chan types.DocumentFragment) ([]interface{}, error) {
// 	return nil, nil
// }
