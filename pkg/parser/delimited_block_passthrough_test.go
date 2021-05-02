package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("passthrough blocks", func() {

	Context("in raw documents", func() {

		Context("delimited blocks", func() {

			It("with title", func() {
				source := `.a title
++++
_foo_

*bar*
++++`
				expected := []types.DocumentFragmentGroup{
					{
						LineOffset: 1,
						Content: []interface{}{

							types.PassthroughBlock{
								Attributes: types.Attributes{
									types.AttrTitle: "a title",
								},
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "_foo_",
										},
									},
									{},
									{
										types.StringElement{
											Content: "*bar*",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragmentGroups(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("with special characters", func() {
				source := `++++
<input>

<input>
++++`
				expected := []types.DocumentFragmentGroup{
					{
						LineOffset: 1,
						Content: []interface{}{

							types.PassthroughBlock{
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "<input>",
										},
									},
									{},
									{
										types.StringElement{
											Content: "<input>",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragmentGroups(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("with inline link", func() {
				source := `++++
http://example.com[]
++++`
				expected := []types.DocumentFragmentGroup{
					{
						LineOffset: 1,
						Content: []interface{}{

							types.PassthroughBlock{
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "http://example.com[]",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragmentGroups(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("with inline pass", func() {
				source := `++++
pass:[foo]
++++`
				expected := []types.DocumentFragmentGroup{
					{
						LineOffset: 1,
						Content: []interface{}{

							types.PassthroughBlock{
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "pass:[foo]",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragmentGroups(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("with quoted text", func() {
				source := `++++
*foo*
++++`
				expected := []types.DocumentFragmentGroup{
					{
						LineOffset: 1,
						Content: []interface{}{

							types.PassthroughBlock{
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "*foo*",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragmentGroups(source)).To(MatchDocumentFragmentGroups(expected))
			})
		})

		Context("paragraph with attribute", func() {

			It("2-line paragraph followed by another paragraph", func() {
				source := `[pass]
_foo_
*bar*

another paragraph`
				expected := []types.DocumentFragmentGroup{
					{
						LineOffset: 1,
						Content: []interface{}{

							types.PassthroughBlock{
								Attributes: types.Attributes{
									types.AttrStyle: types.Passthrough,
								},
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "_foo_",
										},
									},
									{
										types.StringElement{
											Content: "*bar*",
										},
									},
								},
							},
							types.BlankLine{},
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "another paragraph",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragmentGroups(source)).To(MatchDocumentFragmentGroups(expected))
			})
		})

	})

	Context("in final documents", func() {

		Context("delimited blocks", func() {

			It("with title", func() {
				source := `.a title
++++
_foo_

*bar*
++++`
				expected := types.Document{
					Elements: []interface{}{
						types.PassthroughBlock{
							Attributes: types.Attributes{
								types.AttrTitle: "a title",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "_foo_",
									},
								},
								{},
								{
									types.StringElement{
										Content: "*bar*",
									},
								},
							},
						},
					},
				}
				result, err := ParseDocument(source)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(MatchDocument(expected))
			})
		})

		Context("paragraph with attribute", func() {

			It("2-line paragraph followed by another paragraph", func() {
				source := `[pass]
_foo_
*bar*

another paragraph`
				expected := types.Document{
					Elements: []interface{}{
						types.PassthroughBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Passthrough,
							},
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "_foo_",
									},
								},
								{
									types.StringElement{
										Content: "*bar*",
									},
								},
							},
						},
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "another paragraph",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})
	})
})
