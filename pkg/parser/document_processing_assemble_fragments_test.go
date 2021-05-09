package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("document fragment assembling", func() {

	It("should assemble 1 paragraph with single line", func() {
		source := `a line`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Content: &types.Paragraph{
					Elements: []interface{}{
						types.RawLine("a line"),
					},
				},
			},
		}
		Expect(AssembleDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should assemble 2 pagragaphs with single line each", func() {
		source := `a line
		
another line`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Content: &types.Paragraph{
					Elements: []interface{}{
						types.RawLine("a line"),
					},
				},
			},
			{
				LineOffset: 3,
				Content: &types.Paragraph{
					Elements: []interface{}{
						types.RawLine("another line"),
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
				Content: &types.DelimitedBlock{
					Kind: types.Listing,
					Elements: []interface{}{
						types.RawLine("a line"),
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
				Content: &types.DelimitedBlock{
					Kind: types.Listing,
					Elements: []interface{}{
						types.RawLine("a line"),
						types.BlankLine{},
						types.RawLine("****"),
						types.RawLine("not a sidebar block"),
						types.RawLine("****"),
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
				Content: &types.DelimitedBlock{
					Kind: types.Listing,
					Elements: []interface{}{
						types.RawLine("a line"),
						types.BlankLine{},
						types.RawLine("another line"),
					},
				},
			},
			{
				LineOffset: 8,
				Content: &types.Paragraph{
					Elements: []interface{}{
						types.RawLine("a paragraph"),
						types.RawLine("on"),
						types.RawLine("3 lines."),
					},
				},
			},
		}
		Expect(AssembleDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})
})
