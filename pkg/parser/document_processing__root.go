package parser

import (
	"io"

	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	log "github.com/sirupsen/logrus"
)

// ParseDocument parses the content of the reader identitied by the filename
func ParseDocument(r io.Reader, config configuration.Configuration, opts ...Option) (types.Document, error) {
	done := make(chan interface{})
	defer close(done)

	attributes := types.Attributes{}
	blocks := []interface{}{}
	inHeader := true
	ctx := newProcessContext()
	pipeline :=
		// TODO: AggregateSections(...) ?
		ApplySubstitutions(ctx, done,
			SplitElements(done,
				ArrangeLists(done,
					AssembleFragments(done,
						ScanDocument(r, done, opts...),
					),
				),
			),
		)
	for f := range pipeline {
		// if log.IsLevelEnabled(log.DebugLevel) {
		// 	log.WithField("_state", "document_parsing").Debugf("received document fragment: %s", spew.Sdump(f))
		// }
		if err := f.Error; err != nil {
			log.WithField("pipeline_stage", "root").WithError(err).Error("error occurred")
			return types.Document{}, err
		}
		for _, element := range f.Elements {
			switch b := element.(type) {
			case types.Section:
				if b.Level != 0 && inHeader { // do not even allow 2ndary section with level 0 as headers
					inHeader = false
				}
			case *types.AttributeDeclaration:
				if inHeader {
					attributes.Set(b.Name, b.Value)
				}
			case *types.AttributeReset:
				delete(attributes, b.Name)
			default:
				// anything else and we're not in the header anynore
				inHeader = false
			}
			blocks = append(blocks, element)
		}
	}
	// fragments := make(chan types.DocumentFragment)
	// err := ParseDocumentFragmentGroups							(r, fragments, opts...)
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
