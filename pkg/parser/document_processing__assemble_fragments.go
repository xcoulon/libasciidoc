package parser

import (
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
func AssembleFragments(done <-chan interface{}, fragmentGroupStream <-chan types.DocumentFragmentGroup) <-chan types.DocumentFragmentGroup {
	assembledElementStream := make(chan types.DocumentFragmentGroup)
	go func() {
		defer close(assembledElementStream)
		for group := range fragmentGroupStream {
			select {
			case <-done:
				log.WithField("pipeline_task", "assemble_fragments").Debug("received 'done' signal")
				return
			case assembledElementStream <- assembleFragments(group):
			}
		}
		log.WithField("pipeline_task", "assemble_fragments").Debug("done processing upstream content")
	}()
	return assembledElementStream
}

func assembleFragments(f types.DocumentFragmentGroup) types.DocumentFragmentGroup {
	// if the fragment contains an error, then send it as-is downstream
	if f.Error != nil {
		return f
	}
	attributes := newAttributeStack()
	content := make([]interface{}, 0, len(f.Content))
	var block types.WithElements // here, a delimited block or a paragraph
	for _, e := range f.Content {
		switch e := e.(type) {
		case types.Attributes:
			attributes.push(e)
		case *types.BlockDelimiter:
			// opening
			if _, ok := block.(*types.DelimitedBlock); !ok {
				block = types.NewDelimitedBlock(e.Kind, attributes.pop())
				content = append(content, block)
				continue
			}
			// closing (ie, reset block)
			block = nil
		case types.ListElement:
			e.SetAttributes(attributes.pop())
			block = e
			content = append(content, block)
		case *types.BlankLine:
			switch block.(type) {
			case *types.Paragraph, types.ListElement:
				// end of paragraph
				block = nil
			}
			content = append(content, e)
		case types.RawLine:
			if block == nil {
				block, _ = types.NewParagraph([]interface{}{}, attributes.pop())
				content = append(content, block)
			}
			if err := block.AddElement(e); err != nil {
				return types.NewErrorFragmentGroup(f.LineOffset, errors.Wrap(err, "unable to assemble fragments"))
			}
		case *types.AdmonitionLine:
			block = types.NewAdminitionParagraph(e, attributes.pop())
			content = append(content, block)
			log.Debug("adding a new fragment with an admonition paragraph")
		case *types.SingleLineComment:
			// add into block (if applicable)
			if block != nil {
				if err := block.AddElement(e); err != nil {
					return types.NewErrorFragmentGroup(f.LineOffset, errors.Wrap(err, "unable to assemble fragments"))
				}
				continue
			}
			content = append(content, e)
		case *types.ListElementContinuation:
			content = append(content, e)
			// what's coming next shall not be attach to the current block (a list element)
			block = nil
		case types.WithAttributes:
			// set attributes on the target element
			e.SetAttributes(attributes.pop())
			content = append(content, e)
		default:
			content = append(content, e)
		}
	}

	result := types.DocumentFragmentGroup{
		LineOffset: f.LineOffset,
		Content:    content,
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.WithField("pipeline_task", "assemble_fragments").Debugf("assembled fragment group: %s", spew.Sdump(result))
	}
	return result
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
