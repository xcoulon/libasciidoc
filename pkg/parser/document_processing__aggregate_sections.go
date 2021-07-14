package parser

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"

	log "github.com/sirupsen/logrus"
)

// ArrangeSections pipeline task which organizes the sections in hierarchy, and
// keeps track of their references.
// Also takes care of wrapping all blocks between header (section 0) and first child section
// into a `Preamble`
// TODO: Also take care of inserting the Table of Contents
// returns the whole document at once (or an error)
func AggregateSections(done <-chan interface{}, fragmentStream <-chan types.DocumentFragment) (*types.Document, error) {
	doc, toc, err := aggregateSections(done, fragmentStream)
	if err != nil {
		return nil, err
	}
	doc.InsertPreamble()
	doc.InsertTableOfContents(toc)
	return doc, nil
}

func aggregateSections(done <-chan interface{}, fragmentStream <-chan types.DocumentFragment) (*types.Document, *types.TableOfContents, error) {
	attrs := types.Attributes{}
	refs := types.ElementReferences{}
	root := &types.Document{}
	lvls := &levels{
		elements: []types.WithElementAddition{
			root,
		},
	}
	var toc *types.TableOfContents
	for f := range fragmentStream {
		if err := f.Error; err != nil {
			return nil, nil, err
		}
		for _, element := range f.Elements {
			switch element := element.(type) {
			case *types.AttributeDeclaration:
				attrs[element.Name] = element.Value
				if element.Name == types.AttrTableOfContents {
					toc = &types.TableOfContents{}
				}
			case *types.AttributeReset:
				delete(attrs, element.Name)
			case *types.BlankLine:
				// ignore
			case *types.Section:
				if err := element.ResolveID(attrs, refs); err != nil {
					return nil, nil, err
				}
				lvls.appendSection(element)
				if toc != nil {
					toc.Add(element)
				}
			default:
				lvls.appendElement(element)
			}
		}
	}

	log.WithField("pipeline_task", "arrange_sections").Debug("done processing upstream content")
	if len(attrs) > 0 {
		root.Attributes = attrs
	}
	if len(refs) > 0 {
		root.ElementReferences = refs
	}
	return root, toc, nil
}

type levels struct {
	elements []types.WithElementAddition
}

func (l *levels) appendSection(s *types.Section) {
	// note: section levels start at 0, but first level is root (doc)
	if idx, found := l.indexOfParent(s); found {
		l.elements = l.elements[:idx+1]
	}
	log.Debugf("adding section with level %d at position %d in levels", s.Level, len(l.elements))
	// append
	l.elements[len(l.elements)-1].AddElement(s)
	l.elements = append(l.elements, s)
}

// return the index of the parent element for the given section,
// taking account the given section's level, and also gaps in other
// sections (eg: `1,2,4` instead of `0,1,2`)
func (l *levels) indexOfParent(s *types.Section) (int, bool) {
	for i, e := range l.elements {
		if p, ok := e.(*types.Section); ok {
			if p.Level >= s.Level {
				log.Debugf("found parent at index %d for section with level %d", i-1, s.Level)
				return i - 1, true // return previous
			}
		}
	}
	//
	return -1, false
}

func (l levels) appendElement(e interface{}) {
	l.elements[len(l.elements)-1].AddElement(e)
}
