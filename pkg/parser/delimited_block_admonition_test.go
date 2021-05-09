package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("admonition blocks", func() {

	Context("in raw documents", func() {

		Context("delimited blocks", func() {

			It("example block as admonition", func() {
				source := `[NOTE]
====
foo
====`
				expected := []types.DocumentFragmentGroup{
					{
						LineOffset: 1,
						Content: []interface{}{

							types.ExampleBlock{
								Attributes: types.Attributes{
									types.AttrStyle: types.Note,
								},
								Elements: []interface{}{
									types.Paragraph{
										Lines: [][]interface{}{
											{
												types.StringElement{
													Content: "foo",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragmentGroups(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("as admonition", func() {
				source := `[NOTE]
----
multiple

paragraphs
----
`
				expected := []types.DocumentFragmentGroup{
					{
						LineOffset: 1,
						Content: []interface{}{

							types.ListingBlock{
								Attributes: types.Attributes{
									types.AttrStyle: types.Note,
								},
								Elements: []interface{}{
									types.StringElement{
										Content: "multiple",
									},
									types.StringElement{
										Content: "paragraphs",
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

			It("example block as admonition", func() {
				source := `[NOTE]
====
foo
====`
				expected := types.Document{
					Elements: []interface{}{
						types.ExampleBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Note,
							},
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{
												Content: "foo",
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))

			})

			It("example block as admonition with multiple lines", func() {
				source := `[NOTE]
====
multiple

paragraphs
====
`
				expected := types.Document{
					Elements: []interface{}{
						types.ExampleBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Note,
							},
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{
												Content: "multiple",
											},
										},
									},
								},
								types.BlankLine{},
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{
												Content: "paragraphs",
											},
										},
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
