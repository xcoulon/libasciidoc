package parser

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ArrangeLists pipeline task which consists join list elements into lists
func ArrangeLists(done <-chan interface{}, fragmentStream <-chan types.DocumentFragment) <-chan types.DocumentFragment {
	arrangedStream := make(chan types.DocumentFragment)
	go func() {
		defer close(arrangedStream)
		for fragment := range fragmentStream {
			select {
			case <-done:
				log.WithField("pipeline_task", "arrange_lists").Debug("received 'done' signal")
				return
			case arrangedStream <- arrangeLists(fragment):
			}
		}
		log.WithField("pipeline_task", "arrange_lists").Debug("done processing upstream content")
	}()
	return arrangedStream
}

func arrangeLists(f types.DocumentFragment) types.DocumentFragment {
	// if the fragment contains an error, then send it as-is downstream
	if err := f.Error; err != nil {
		log.Debugf("skipping list elements arrangement: %v", f.Error)
		return f
	}
	log.Debug("arranging list elements")
	elements, err := arrangeListElements(f.Elements)
	if e, ok := err.(fragmentError); ok {
		return types.NewErrorFragment(f.LineOffset+e.lineOffset, err)
	} else if err != nil {
		return types.NewErrorFragment(f.LineOffset, err)
	}
	result := types.DocumentFragment{
		LineOffset: f.LineOffset,
		Elements:   elements,
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.WithField("pipeline_task", "arrange_lists").Debugf("arranged lists: %s", spew.Sdump(result))
	}

	return result
}

func arrangeListElements(elements []interface{}) ([]interface{}, error) {
	result := make([]interface{}, 0, len(elements))
	blocks := newBlockStack() // so we can support delimited blocks in list elements, etc.

content:
	for i, elements := range elements {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.WithField("pipeline_task", "arrange_lists").Debugf("arranging element of type '%T'", elements)
		}
		// case types.AttributeDeclaration, types.AttributeReset:
		// 	log.Debug("adding a new fragment with an attribute declaration/reset")
		// 	result = append(result, types.NewDocumentFragment(f.LineOffset, e))

		// lookup the parent block that can add the given element
		for {
			parentBlock := blocks.get()
			if parentBlock == nil {
				break
			}
			if parentBlock.CanAddElement(elements) {
				if err := parentBlock.AddElement(elements); err != nil {
					return nil, fragmentError{
						err:        errors.Wrap(err, "unable to assemble fragments"),
						lineOffset: i,
					}
				}
				continue content
			}
			log.Debugf("couldn't add element of type '%T' to block of type '%T'", elements, parentBlock)
			blocks.pop()
		}

		switch e := elements.(type) {
		// case *types.Section:
		// 	blocks.push(e)
		// 	log.Debug("adding a new fragment with a section")
		// 	result = append(result, types.NewDocumentFragment(f.LineOffset, e))
		case *types.DelimitedBlock:
			log.Debug("checking list elements in delimitedblock")
			var err error
			if e.Elements, err = arrangeListElements(e.Elements); err != nil {
				return nil, err
			}
			blocks.push(e)
			result = append(result, e)
		case *types.BlankLine, *types.SingleLineComment:
			// end of current block
			blocks.pop()
			log.Debug("adding a new fragment with a blankline or singleline comment")
			result = append(result, e)
		case types.ListElement:
			list, err := types.NewList(e)
			if err != nil {
				return nil, fragmentError{
					err:        errors.Wrap(err, "unable to assemble fragments"),
					lineOffset: i,
				}
			}
			blocks.push(list)
			log.Debug("adding a new list element")
			result = append(result, list)
		case *types.ListElementContinuation:
			// fall back for a list element continuation which is not following a list element
			block, _ := types.NewParagraph([]interface{}{
				&types.StringElement{
					Content: "+\n",
				},
			}, nil)
			blocks.push(block)
			log.Debug("adding a new fragment with a paragraph")
			result = append(result, e)
		default:
			log.Debugf("adding a new fragment with an element of type '%T' (default case)", e)
			result = append(result, e)
		}
	}
	if len(result) == 0 {
		result = nil
	}
	return result, nil
}

type fragmentError struct {
	err        error
	lineOffset int
}

var _ error = fragmentError{}

func (e fragmentError) Error() string {
	return e.err.Error()
}

type blockStack struct {
	index int
	stack []types.WithElements
}

func newBlockStack() *blockStack {
	return &blockStack{
		stack: make([]types.WithElements, 10),
		index: -1,
	}
}

func (s *blockStack) push(a types.WithElements) {
	s.index++
	s.stack[s.index] = a
}

func (s *blockStack) pop() types.WithElements {
	if s.index < 0 {
		return nil
	}
	a := s.stack[s.index]
	s.index--
	return a
}

func (s *blockStack) get() types.WithElements {
	if s.index < 0 {
		return nil
	}
	return s.stack[s.index]
}
