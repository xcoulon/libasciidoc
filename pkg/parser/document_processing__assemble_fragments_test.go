package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("document element assembling", func() {

	It("should assemble 1 paragraph with single line", func() {
		source := `a line`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							types.RawLine("a line"),
						},
					},
				},
			},
		}
		Expect(AssembleDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should assemble 2 paragraphs with single line each", func() {
		source := `a line
		
another line`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							types.RawLine("a line"),
						},
					},
					&types.BlankLine{},
				},
			},
			{
				LineOffset: 3,
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							types.RawLine("another line"),
						},
					},
				},
			},
		}
		Expect(AssembleDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should assemble 1 delimited block with single rawline", func() {
		source := `----
a line
----`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					&types.DelimitedBlock{
						Kind: types.Listing,
						Elements: []interface{}{
							types.RawLine("a line"),
						},
					},
				},
			},
		}
		Expect(AssembleDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should collect 1 delimited block with multiple rawlines only", func() {
		source := `----
a line

****
not a sidebar block
****
----
`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					&types.DelimitedBlock{
						Kind: types.Listing,
						Elements: []interface{}{
							types.RawLine("a line"),
							types.RawLine(""),
							types.RawLine("****"),
							types.RawLine("not a sidebar block"),
							types.RawLine("****"),
						},
					},
				},
			},
		}
		Expect(AssembleDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should collect 1 delimited block with multiple rawlines and content afterwards", func() {
		source := `----
a line

another line
----


a paragraph
on
3 lines.

`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					&types.DelimitedBlock{
						Kind: types.Listing,
						Elements: []interface{}{
							types.RawLine("a line"),
							types.RawLine(""),
							types.RawLine("another line"),
						},
					},
					&types.BlankLine{},
				},
			},
			{
				LineOffset: 7,
				Elements: []interface{}{
					&types.BlankLine{},
				},
			},
			{
				LineOffset: 8,
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							types.RawLine("a paragraph"),
							types.RawLine("on"),
							types.RawLine("3 lines."),
						},
					},
					&types.BlankLine{},
				},
			},
		}
		Expect(AssembleDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should collect 1 section and content afterwards", func() {
		source := `== section title


a paragraph
on
3 lines.

`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					&types.Section{
						Level: 1,
						Title: []interface{}{
							types.StringElement{
								Content: "section title",
							},
						},
					},
					&types.BlankLine{},
				},
			},
			{
				LineOffset: 3,
				Elements: []interface{}{
					&types.BlankLine{},
				},
			},
			{
				LineOffset: 4,
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							types.RawLine("a paragraph"),
							types.RawLine("on"),
							types.RawLine("3 lines."),
						},
					},
					&types.BlankLine{},
				},
			},
		}
		Expect(AssembleDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should assemble callout list items without blankline in-between", func() {
		source := `<1> first item
<2> second item
`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					&types.CalloutListElement{
						Ref: 1,
						Elements: []interface{}{
							&types.Paragraph{
								Elements: []interface{}{
									types.RawLine("first item"),
								},
							},
						},
					},
					&types.CalloutListElement{
						Ref: 2,
						Elements: []interface{}{
							&types.Paragraph{
								Elements: []interface{}{
									types.RawLine("second item"),
								},
							},
						},
					},
				},
			},
		}
		Expect(AssembleDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should assemble callout list items with blanklines in-between", func() {
		source := `<1> first item

		
<2> second item
`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					&types.CalloutListElement{
						Ref: 1,
						Elements: []interface{}{
							&types.Paragraph{
								Elements: []interface{}{
									types.RawLine("first item"),
								},
							},
						},
					},
					&types.BlankLine{},
					&types.BlankLine{},
					&types.CalloutListElement{
						Ref: 2,
						Elements: []interface{}{
							&types.Paragraph{
								Elements: []interface{}{
									types.RawLine("second item"),
								},
							},
						},
					},
				},
			},
		}
		Expect(AssembleDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should assemble ordered list items without blanklines in-between", func() {
		source := `. first item
. second item
`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					&types.OrderedListElement{
						Style: types.Arabic,
						Elements: []interface{}{
							&types.Paragraph{
								Elements: []interface{}{
									types.RawLine("first item"),
								},
							},
						},
					},
					&types.OrderedListElement{
						Style: types.Arabic,
						Elements: []interface{}{
							&types.Paragraph{
								Elements: []interface{}{
									types.RawLine("second item"),
								},
							},
						},
					},
				},
			},
		}
		Expect(AssembleDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should assemble ordered list items with blanklines in-between", func() {
		source := `. first item


. second item
`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					&types.OrderedListElement{
						Style: types.Arabic,
						Elements: []interface{}{
							&types.Paragraph{
								Elements: []interface{}{
									types.RawLine("first item"),
								},
							},
						},
					},
					&types.BlankLine{},
					&types.BlankLine{},
					&types.OrderedListElement{
						Style: types.Arabic,
						Elements: []interface{}{
							&types.Paragraph{
								Elements: []interface{}{
									types.RawLine("second item"),
								},
							},
						},
					},
				},
			},
		}
		Expect(AssembleDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})
})
