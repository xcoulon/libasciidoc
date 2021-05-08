package parser

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func AssembleFragments(done <-chan interface{}, fragmentStream <-chan types.DocumentFragmentGroup) <-chan types.DocumentFragment {
	assembledFragmentStream := make(chan types.DocumentFragment)
	go func() {
		defer close(assembledFragmentStream)
		for fragments := range fragmentStream {
			for _, fragment := range AssembleFragmentElements(fragments) {
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

func AssembleFragmentElements(f types.DocumentFragmentGroup) []types.DocumentFragment {
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
	var block types.BlockWithNestedElements // a delimited block or a paragraph
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
		case types.AttributeDeclaration, types.AttributeReset:
			result = append(result, types.NewDocumentFragment(f.LineOffset, e))
		case *types.Section:
			e.Attributes = attributes
			attributes = types.Attributes{} // reset
			// send the section block downstream
			result = append(result, types.NewDocumentFragment(f.LineOffset, e))
		// case types.BlockDelimiter:
		// 	// TODO: support nested blocks with the help of a Stack object?
		// 	if block == nil {
		// 		block = types.NewDelimitedBlock(e.Kind, attributes)
		// 		attributes = nil // reset
		// 		result = append(result, types.NewDocumentFragment(f.LineOffset, block))
		// 	} else {
		// 		block = nil // reset, in case there is more content afterwards
		// 	}
		case types.RawLine:
			if block == nil {
				block, _ = types.NewParagraph([]interface{}{}, attributes)
				attributes = nil // reset
				result = append(result, types.NewDocumentFragment(f.LineOffset, block))
			}
			block.AddElement(e)
		case types.SingleLineComment, *types.ImageBlock:
			result = append(result, types.NewDocumentFragment(f.LineOffset+i, e))
		case types.BlankLine:
			// do nothing for now
		default:
			// unknow type fragment element: set an error on the fragment and send it downstream
			fr := types.NewDocumentFragment(f.LineOffset, block)
			fr.Error = errors.Errorf("unexpected type of element on line %d: '%T'", f.LineOffset+i, e)
			result = append(result, fr)
		}
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.WithField("state", "fragment_assembling").Debugf("assembled fragments: %s", spew.Sdump(result))
	}

	return result
}
