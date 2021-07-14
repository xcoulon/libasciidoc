package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("example blocks", func() {

	Context("in raw documents", func() {

		Context("as delimited blocks", func() {

			It("with single rich line", func() {
				source := `====
some *example* content
====`
				expected := []types.DocumentFragment{
					{
						LineOffset: 1,
						Elements: []interface{}{
							&types.BlockDelimiter{
								Kind: types.Example,
							},
							types.RawLine("some *example* content"),
							&types.BlockDelimiter{
								Kind: types.Example,
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("with single line starting with a dot", func() {
				source := `====
.foo
====`
				expected := []types.DocumentFragment{
					{
						LineOffset: 1,
						Elements: []interface{}{
							&types.BlockDelimiter{
								Kind: types.Example,
							},
							types.Attributes{
								types.AttrTitle: "foo",
							},
							&types.BlockDelimiter{
								Kind: types.Example,
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("with rich lines", func() {
				source := `====
.foo
some *example* content
with _italic content_

* and a list item
====`
				expected := []types.DocumentFragment{
					{
						LineOffset: 1,
						Elements: []interface{}{
							&types.BlockDelimiter{
								Kind: types.Example,
							},
							types.Attributes{
								types.AttrTitle: "foo",
							},
							types.RawLine("some *example* content"),
							types.RawLine("with _italic content_"),
							&types.BlankLine{},
							&types.UnorderedListElement{
								BulletStyle: types.OneAsterisk,
								CheckStyle:  types.NoCheck,
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.RawLine("and a list item"),
										},
									},
								},
							},
							&types.BlockDelimiter{
								Kind: types.Example,
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("with unclosed delimiter", func() {
				source := `====
End of doc here`
				expected := []types.DocumentFragment{
					{
						LineOffset: 1,
						Elements: []interface{}{
							&types.BlockDelimiter{
								Kind: types.Example,
							},
							types.RawLine("End of doc here"),
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("with title", func() {
				source := `.example block title
====
foo
====`
				expected := []types.DocumentFragment{
					{
						LineOffset: 1,
						Elements: []interface{}{
							types.Attributes{
								types.AttrTitle: "example block title",
							},
							&types.BlockDelimiter{
								Kind: types.Example,
							},
							types.RawLine("foo"),
							&types.BlockDelimiter{
								Kind: types.Example,
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("with caption", func() {
				source := `[caption="a caption "]
====
foo
====`
				expected := []types.DocumentFragment{
					{
						LineOffset: 1,
						Elements: []interface{}{
							types.Attributes{
								types.AttrCaption: "a caption ", // trailing space is retained
							},
							&types.BlockDelimiter{
								Kind: types.Example,
							},
							types.RawLine("foo"),
							&types.BlockDelimiter{
								Kind: types.Example,
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("example block starting delimiter only", func() {
				source := `====`
				expected := []types.DocumentFragment{
					{
						LineOffset: 1,
						Elements: []interface{}{
							&types.BlockDelimiter{
								Kind: types.Example,
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})
		})

		Context("as paragraph blocks", func() {

			It("with single rich line", func() {
				source := `[example]
some *example* content`
				expected := []types.DocumentFragment{
					{
						LineOffset: 1,
						Elements: []interface{}{
							&types.Paragraph{
								Attributes: types.Attributes{
									types.AttrStyle: types.Example,
								},
								Elements: []interface{}{
									&types.StringElement{
										Content: "some ",
									},
									&types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											&types.StringElement{
												Content: "example",
											},
										},
									},
									&types.StringElement{
										Content: " content",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})
		})
	})

	Context("in final documents", func() {

		Context("as delimited blocks", func() {

			It("with single rich line", func() {
				source := `====
some *example* content
====`
				expected := &types.Document{
					Elements: []interface{}{
						&types.DelimitedBlock{
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
									Elements: []interface{}{
										&types.StringElement{
											Content: "some ",
										},
										&types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												&types.StringElement{
													Content: "example",
												},
											},
										},
										&types.StringElement{
											Content: " content",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with single line starting with a dot", func() {
				source := `====
.foo
====`
				expected := &types.Document{
					Elements: []interface{}{
						&types.DelimitedBlock{
							Kind: types.Example,
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with multiple lines", func() {
				source := `====
.foo
some listing code
with *bold content*

* and a list item
====`
				expected := &types.Document{
					Elements: []interface{}{
						&types.DelimitedBlock{
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
									Attributes: types.Attributes{
										types.AttrTitle: "foo",
									},
									Elements: []interface{}{
										&types.StringElement{
											Content: "some listing code\nwith ",
										},
										&types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												&types.StringElement{
													Content: "bold content",
												},
											},
										},
									},
								},
								&types.BlankLine{},
								&types.GenericList{
									Kind: types.UnorderedListKind,
									Elements: []types.ListElement{
										&types.UnorderedListElement{
											BulletStyle: types.OneAsterisk,
											CheckStyle:  types.NoCheck,
											Elements: []interface{}{
												&types.Paragraph{
													Elements: []interface{}{
														&types.StringElement{
															Content: "and a list item",
														},
													},
												},
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

			It("with unclosed delimiter", func() {
				source := `====
End of file here`
				expected := &types.Document{
					Elements: []interface{}{
						&types.DelimitedBlock{
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
									Elements: []interface{}{
										&types.StringElement{
											Content: "End of file here",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with title", func() {
				source := `.example block title
====
foo
====`
				expected := &types.Document{
					Elements: []interface{}{
						&types.DelimitedBlock{
							Kind: types.Example,
							Attributes: types.Attributes{
								types.AttrTitle: "example block title",
							},
							Elements: []interface{}{
								&types.Paragraph{
									Elements: []interface{}{
										&types.StringElement{
											Content: "foo",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("example block starting delimiter only", func() {
				source := `====`
				expected := &types.Document{
					Elements: []interface{}{
						&types.DelimitedBlock{
							Kind: types.Example,
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})

		Context("as paragraph blocks", func() {

			It("with single rich line", func() {
				source := `[example]
some *example* content`
				expected := &types.Document{
					Elements: []interface{}{
						&types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle: types.Example,
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "some ",
								},
								&types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										&types.StringElement{
											Content: "example",
										},
									},
								},
								&types.StringElement{
									Content: " content",
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
	})
})
