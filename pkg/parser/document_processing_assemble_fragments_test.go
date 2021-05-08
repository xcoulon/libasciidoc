package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("document fragment assembling", func() {

	Context("raw lines", func() {

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
	})

	Context("delimited blocks lines", func() {
		It("should assemble 1 delimited block with single line", func() {
			source := `----
a line
----`
			expected := []types.DocumentFragment{
				{
					LineOffset: 1,
					Content: &types.ListingBlock{

						Elements: []interface{}{
							types.RawLine("a line"),
						},
					},
				},
			}
			Expect(AssembleDocumentFragments(source)).To(MatchDocumentFragments(expected))
		})
	})
})
