package parser

import (
	"fmt"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/pkg/errors"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

// AssembleFragments pipeline task which consist in grouping some elements into blocks. For example:
// - delimited blocks using start/end delimiters and rawlines in-between
// - paragraphs using rawlines
//
// also, this tasks takes care of "attaching" attributes to their parent block/element
func AssembleFragments(done <-chan interface{}, fragmentStream <-chan types.DocumentFragment) <-chan types.DocumentFragment {
	assembledFragmentStream := make(chan types.DocumentFragment)
	go func() {
		defer close(assembledFragmentStream)
		for fragment := range fragmentStream {
			select {
			case <-done:
				log.WithField("pipeline_task", "assemble_fragment_elements").Debug("received 'done' signal")
				return
			case assembledFragmentStream <- assembleFragments(fragment):
			}
		}
		log.WithField("pipeline_task", "assemble_fragment_elements").Debug("done processing upstream content")
	}()
	return assembledFragmentStream
}

func assembleFragments(f types.DocumentFragment) types.DocumentFragment {
	// if the fragment contains an error, then send it as-is downstream
	if f.Error != nil {
		return f
	}
	elements, err := assembleFragmentElements(f.Elements)
	if err != nil {
		return types.NewErrorFragment(f.LineOffset, errors.Wrap(err, "unable to assemble fragments"))
	}
	return types.DocumentFragment{
		LineOffset: f.LineOffset,
		Elements:   elements,
	}
}

func assembleFragmentElements(elements []interface{}) ([]interface{}, error) {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.WithField("pipeline_task", "assemble_fragment_elements").Debugf("assembling %d elements", len(elements))
	}
	attributes := newAttributeStack()
	result := make([]interface{}, 0, len(elements))
	// parentBlocks := newBlockStack()
	var parentBlock types.WithElements // here, a delimited block or a paragraph
	for _, e := range elements {
		log.Debugf("assembling element of type '%T' (current parent block type: '%T')", e, parentBlock)
		switch e := e.(type) {
		case types.Attributes:
			// within a delimited block, we need to add this attributes as-is
			if b, ok := parentBlock.(*types.DelimitedBlock); ok {
				log.Debugf("adding attributes as an element of a delimited block of kind '%v'", b.Kind)
				if err := b.AddElement(e); err != nil {
					return nil, err
				}
				continue
			}
			// outside of a block, we keep the attributes in a stack
			attributes.push(e)
		case *types.BlockDelimiter:
			// closing if parent block is a delimited block of the same kind
			b, ok := parentBlock.(*types.DelimitedBlock)
			if ok && b.Kind == e.Kind {
				// closing the parent block
				continue
			}
			// nested delimited block (will be assembled later)
			if ok && b.Kind != e.Kind {
				attrs := attributes.pop()
				if err := b.AddElement(attrs); err != nil {
					return nil, err
				}
				if err := b.AddElement(e); err != nil {
					return nil, err
				}
				continue
			}
			// otherwise, let's init a new delimited block
			parentBlock = types.NewDelimitedBlock(e.Kind, attributes.pop())
			result = append(result, parentBlock)
		case types.ListElement:
			e.SetAttributes(attributes.pop())
			// if the current block can take this list element, then let's add it
			if parentBlock != nil && parentBlock.CanAddElement(e) {
				if err := parentBlock.AddElement(e); err != nil {
					return nil, err
				}
			} else {
				// parentBlocks.push(e)
				parentBlock = e
				result = append(result, e)
			}
		case *types.BlankLine:
			switch b := parentBlock.(type) {
			case *types.Paragraph, types.ListElement:
				// end of paragraph
				// parentBlocks.pop()
				result = append(result, e)
				parentBlock = nil
			case *types.DelimitedBlock:
				if b.CanAddElement(e) {
					if err := b.AddElement(e); err != nil {
						return nil, err
					}
				} else {
					result = append(result, e)
				}
			default:
				result = append(result, e)
			}
		case types.RawLine:
			if parentBlock != nil && parentBlock.CanAddElement(e) {
				if err := parentBlock.AddElement(e); err != nil {
					return nil, err
				}
				continue
			}
			parentBlock, _ = types.NewParagraph([]interface{}{e}, attributes.pop())
			// parentBlocks.push(parentBlock)
			result = append(result, parentBlock)
		case *types.AdmonitionLine:
			parentBlock = types.NewAdminitionParagraph(e, attributes.pop())
			// parentBlocks.push(parentBlock)
			result = append(result, parentBlock)
			log.Debug("adding a new fragment with an admonition paragraph")
		case *types.SingleLineComment:
			// add into block (if applicable)
			if parentBlock != nil {
				if err := parentBlock.AddElement(e); err != nil {
					return nil, err
				}
				continue
			}
			result = append(result, e)
		case *types.ListElementContinuation:
			// at this stage, we just append the element to the result set, and remove the parent block from the stack
			result = append(result, e)
			parentBlock = nil
		case types.WithAttributes:
			// set attributes on the target element
			e.SetAttributes(attributes.pop())
			result = append(result, e)
		default:
			result = append(result, e)
		}
	}
	// for delimited blocks (even unclosed) with normal content, redo the assembly with their own content
	for _, block := range result {
		if b, ok := block.(*types.DelimitedBlock); ok {
			switch b.Kind {
			case types.Listing:
				// verbatim content: do nothing
			case types.Example, types.Quote:
				var err error
				if b.Elements, err = assembleFragmentElements(b.Elements); err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("unsupported kind of delimited block: '%s'", b.Kind)
			}
		}
	}
	// no need to return empty slices
	if len(result) == 0 {
		result = nil
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.WithField("pipeline_task", "assemble_fragment_elements").Debugf("assembled fragment elements: %s", spew.Sdump(result))
	}
	return result, nil
}

type attributeStack struct {
	attrs types.Attributes
}

func newAttributeStack() *attributeStack {
	return &attributeStack{}
}

func (s *attributeStack) push(attrs types.Attributes) {
	if s.attrs == nil {
		s.attrs = types.Attributes{}
	}
	s.attrs.SetAll(attrs)
}

func (s *attributeStack) pop() types.Attributes {
	attrs := s.attrs
	s.attrs = nil // reset
	return attrs
}
