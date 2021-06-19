package parser

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func AssembleFragments(done <-chan interface{}, fragmentGroupStream <-chan types.DocumentFragmentGroup) <-chan types.DocumentFragment {
	assembledFragmentStream := make(chan types.DocumentFragment)
	go func() {
		defer close(assembledFragmentStream)
		for group := range fragmentGroupStream {
			for _, fragment := range assembleFragments(group) {
				select {
				case <-done:
					log.WithField("pipeline_stage", "fragment_assembling").Debug("received 'done' signal")
					return
				case assembledFragmentStream <- fragment:
				}
			}
		}
		log.Debug("end of fragment assembly")
	}()
	return assembledFragmentStream
}

func assembleFragments(f types.DocumentFragmentGroup) []types.DocumentFragment {
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
	var attributes types.Attributes // TODO: use a special kind of stack where `pop()` returns all the attributes (merged) and resets the stack?
	blocks := newBlockStack()
	for i, e := range f.Content {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.WithField("pipeline_stage", "fragment_assembling").Debugf("assembling fragment content of type '%T'", e)
		}
		switch e := e.(type) {
		case types.Attributes:
			if attributes == nil {
				attributes = types.Attributes{}
			}
			attributes.SetAll(e)
		case types.AttributeDeclaration, types.AttributeReset:
			log.Debug("adding a new fragment with an attribute declaration/reset")
			result = append(result, types.NewDocumentFragment(f.LineOffset, e))
		case *types.Section:
			e.Attributes = attributes
			attributes = types.Attributes{} // reset
			// send the section block downstream
			log.Debug("adding a new fragment with a section")
			result = append(result, types.NewDocumentFragment(f.LineOffset, e))
		case types.BlockDelimiter:
			switch existing := blocks.get().(type) {
			case nil:
				delimitedBlock := types.NewDelimitedBlock(e.Kind, attributes)
				attributes = nil // reset
				log.Debug("adding a new fragment with a delimitedBlock (no parent block)")
				result = append(result, types.NewDocumentFragment(f.LineOffset, delimitedBlock)) // TODO: move this into the `push()` method?
				blocks.push(delimitedBlock)
			case *types.DelimitedBlock:
				if existing.Kind == e.Kind {
					// "closing" the delimited block
					blocks.pop()
				}
				// TODO: handle cases of nested delimited blocks here
			case *types.GenericList:
				delimitedBlock := types.NewDelimitedBlock(e.Kind, attributes)
				attributes = nil // reset
				// add to existing only if its last element is a `ListElementContinuation`
				if existing.ListElementContinuation() {
					if err := existing.AddElement(delimitedBlock); err != nil {
						result = append(result, types.NewErrorFragment(f.LineOffset, errors.Wrap(err, "unable to assemble fragments")))
					}
				} else {
					log.Debug("adding a new fragment with a delimitedBlock (no list element continuation)")
					result = append(result, types.NewDocumentFragment(f.LineOffset, delimitedBlock)) // TODO: move this into the `push()` method?
					// also, remove existing from stack
					blocks.pop()
				}
				blocks.push(delimitedBlock)
				log.Debugf("added delimited block of king '%s' into parent block of type '%T'", e.Kind, existing)
			default:
				result = append(result, types.NewErrorFragment(f.LineOffset, errors.Errorf("unable to assemble delimited block of kind '%s'", e.Kind)))
			}
		case types.RawLine:
			block := blocks.get()
			if block == nil {
				// by default, append to a paragraph
				block, _ = types.NewParagraph([]interface{}{}, attributes)
				blocks.push(block)
				attributes = nil // reset
				log.Debug("adding a new fragment with a paragraph")
				result = append(result, types.NewDocumentFragment(f.LineOffset, block))
			} else if l, ok := block.(*types.GenericList); ok && !l.CanAddElement() {
				block, _ = types.NewParagraph([]interface{}{}, attributes)
				attributes = nil // reset
				blocks.push(block)
				log.Debug("adding a new fragment with a paragraph")
				result = append(result, types.NewDocumentFragment(f.LineOffset, block))
			}
			if err := block.AddElement(e); err != nil {
				result = append(result, types.NewErrorFragment(f.LineOffset, errors.Wrap(err, "unable to assemble fragments")))
			}
			log.Debugf("added rawline into parent block of type '%T'", block)
		case *types.AdmonitionLine:
			block := blocks.get()
			if block == nil {
				// by default, append to a paragraph
				attributes = attributes.Set(types.AttrStyle, e.Kind)
				block, _ = types.NewParagraph([]interface{}{}, attributes)
				blocks.push(block)
				attributes = nil // reset
				log.Debug("adding a new fragment with a paragraph")
				result = append(result, types.NewDocumentFragment(f.LineOffset, block))
			} else if l, ok := block.(*types.GenericList); ok && !l.CanAddElement() {
				attributes = attributes.Set(types.AttrStyle, e.Kind)
				block, _ = types.NewParagraph([]interface{}{}, attributes)
				attributes = nil // reset
				blocks.push(block)
				log.Debug("adding a new fragment with a paragraph")
				result = append(result, types.NewDocumentFragment(f.LineOffset, block))
			}
			if err := block.AddElement(e); err != nil {
				result = append(result, types.NewErrorFragment(f.LineOffset, errors.Wrap(err, "unable to assemble fragments")))
			}
			log.Debugf("added admonition line into parent block of type '%T'", block)
		case types.BlankLine:
			// behaviour depends on the current block which is being processed
			switch b := blocks.get().(type) {
			case *types.DelimitedBlock:
				if err := b.AddElement(e); err != nil {
					result = append(result, types.NewErrorFragment(f.LineOffset, errors.Wrap(err, "unable to assemble fragments")))
				}
			case *types.GenericList:
				if err := b.AddElement(e); err != nil {
					result = append(result, types.NewErrorFragment(f.LineOffset, errors.Wrap(err, "unable to assemble fragments")))
				}
			case *types.Paragraph:
				blocks.pop() // end of paragraph
			default:
			}
		case types.ListElement:
			e.SetAttributes(attributes)
			attributes = nil // reset
			// if there is no "root" list yet, create one
			switch list := blocks.get().(type) {
			case nil:
				var err error
				list, err = types.NewList(e)
				if err != nil {
					result = append(result, types.NewErrorFragment(f.LineOffset, errors.Wrap(err, "unable to assemble fragments")))
					continue
				}
				log.Debug("adding a new fragment with a list")
				result = append(result, types.NewDocumentFragment(f.LineOffset, list))
				blocks.push(list)
			case *types.GenericList:
				// add the element to the list
				if err := list.AddElement(e); err != nil {
					result = append(result, types.NewErrorFragment(f.LineOffset, errors.Wrap(err, "unable to assemble fragments")))
				}
			}
		case *types.ListElementContinuation:
			list, ok := blocks.get().(*types.GenericList)
			if !ok {
				// TODO: report an error
				log.Error("found a list continuation but there was no list element before")
				// ignore
				continue
			}
			if err := list.AddElement(e); err != nil {
				result = append(result, types.NewErrorFragment(f.LineOffset, errors.Wrap(err, "unable to assemble fragments")))
			}
		case types.SingleLineComment, *types.ImageBlock, *types.ThematicBreak:
			log.Debug("adding a new fragment with a singleline comment or an imageblock")
			result = append(result, types.NewDocumentFragment(f.LineOffset+i, e))
		default:
			// unknow type fragment element: set an error on the fragment and send it downstream
			result = append(result, types.NewErrorFragment(f.LineOffset, errors.Errorf("unable to assemble fragments: unexpected type of element on line %d: '%T'", f.LineOffset+i, e)))
		}
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.WithField("pipeline_stage", "fragment_assembling").Debugf("assembled fragments: %s", spew.Sdump(result))
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

// func (s *scanningScopes) size() int {
// 	return s.index + 1
// }

// func (s *scanningScopes) empty() bool {
// 	return s.index == -1
// }

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
