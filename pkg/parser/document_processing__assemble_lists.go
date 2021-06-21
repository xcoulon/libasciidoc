package parser

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ArrangeLists pipeline task which consists in grouping list elements into lists (with continutions, etc)
func ArrangeLists(done <-chan interface{}, fragmentGroupStream <-chan types.DocumentFragmentGroup) <-chan types.DocumentFragment {
	assembledFragmentStream := make(chan types.DocumentFragment)
	go func() {
		defer close(assembledFragmentStream)
		for group := range fragmentGroupStream {
			for _, fragment := range arrangeLists(group) {
				select {
				case <-done:
					log.WithField("pipeline_task", "assemble_lists").Debug("received 'done' signal")
					return
				case assembledFragmentStream <- fragment:
				}
			}
		}
		log.WithField("pipeline_task", "assemble_lists").Debug("done processing upstream content")
	}()
	return assembledFragmentStream
}

func arrangeLists(f types.DocumentFragmentGroup) []types.DocumentFragment {
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
	blocks := newBlockStack() // so we can support delimited blocks in list elements, etc.
content:
	for i, block := range f.Content {

		// TODO: refactor the whole approach: "bufferize" the blocks while being built,
		// and only add/append to the parent block once the current block is complete
		// (for example, when a delimited block is closed, when a blankline is found after a paragraph, etc.)

		if log.IsLevelEnabled(log.DebugLevel) {
			log.WithField("pipeline_task", "assemble_lists").Debugf("assembling fragment content of type '%T'", block)
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
			if parentBlock.CanAddElement(block) {
				if err := parentBlock.AddElement(block); err != nil {
					result = append(result, types.NewErrorFragment(f.LineOffset, errors.Wrap(err, "unable to assemble fragments")))
				}
				continue content
			}
			log.Debugf("couldn't add element of type '%T' to block of type '%T'", block, parentBlock)
			blocks.pop()
		}

		switch e := block.(type) {
		// no need to put such an element on top of the stack
		case *types.Section:
			blocks.push(e)
			log.Debug("adding a new fragment with a section")
			result = append(result, types.NewDocumentFragment(f.LineOffset, e))
		case *types.DelimitedBlock:
			log.Debug("adding a new fragment with a delimitedBlock")
			result = append(result, types.NewDocumentFragment(f.LineOffset, e)) // TODO: move this into the `push()` method?
			blocks.push(e)
		case *types.BlankLine, *types.SingleLineComment:
			// end of current block
			blocks.pop()
			log.Debug("adding a new fragment with a blankline or singleline comment")
			result = append(result, types.NewDocumentFragment(f.LineOffset+i, e))
		case types.ListElement:
			list, err := types.NewList(e)
			if err != nil {
				result = append(result, types.NewErrorFragment(f.LineOffset, errors.Wrap(err, "unable to assemble fragments")))
				continue
			}
			blocks.push(list)
			log.Debug("adding a new fragment with a list")
			result = append(result, types.NewDocumentFragment(f.LineOffset, list))
		case *types.ListElementContinuation:
			// fall back for a list element continuation which is not following a list element
			block, _ := types.NewParagraph([]interface{}{
				types.StringElement{
					Content: "+\n",
				},
			}, nil)
			blocks.push(block)
			log.Debug("adding a new fragment with a paragraph")
			result = append(result, types.NewDocumentFragment(f.LineOffset, block))
		case *types.ImageBlock, *types.ThematicBreak, *types.Paragraph:
			log.Debugf("adding a new fragment with an element of type '%T'", e)
			result = append(result, types.NewDocumentFragment(f.LineOffset+i, e))
		default:
			// unknow type fragment element: set an error on the fragment and send it downstream
			result = append(result, types.NewErrorFragment(f.LineOffset, errors.Errorf("unable to assemble fragments: unexpected type of element on line %d: '%T'", f.LineOffset+i, e)))
		}
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.WithField("pipeline_task", "assemble_lists").Debugf("assembled fragments/lists: %s", spew.Sdump(result))
	}

	return result
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
