package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("example blocks", func() {

	Context("in raw documents", func() {

		Context("delimited blocks", func() {

			It("with single rich line", func() {
				source := `====
some *example* content
====`
				expected := types.DocumentFragments{
					types.ExampleBlock{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "some ",
										},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{
													Content: "example",
												},
											},
										},
										types.StringElement{
											Content: " content",
										},
									},
								},
							},
						},
					},
				}

				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with single line starting with a dot", func() {
				source := `====
.foo
====`
				expected := types.DocumentFragments{
					types.ExampleBlock{
						Elements: []interface{}{
							types.StandaloneAttributes{
								types.AttrTitle: "foo",
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with rich lines", func() {
				source := `====
.foo
some listing *bold code*
with _italic content_

* and a list item
====`
				expected := types.DocumentFragments{
					types.ExampleBlock{
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrTitle: "foo",
								},
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "some listing ",
										},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{
													Content: "bold code",
												},
											},
										},
									},
									{
										types.StringElement{
											Content: "with ",
										},
										types.QuotedText{
											Kind: types.SingleQuoteItalic,
											Elements: []interface{}{
												types.StringElement{
													Content: "italic content",
												},
											},
										},
									},
								},
							},
							types.BlankLine{},
							types.UnorderedListItem{
								Level:       1,
								BulletStyle: types.OneAsterisk,
								CheckStyle:  types.NoCheck,
								Elements: []interface{}{
									types.Paragraph{
										Lines: [][]interface{}{
											{
												types.StringElement{
													Content: "and a list item",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with unclosed delimiter", func() {
				source := `====
End of doc here`
				expected := types.DocumentFragments{
					types.ExampleBlock{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "End of doc here",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with title", func() {
				source := `.example block title
====
foo
====`
				expected := types.DocumentFragments{
					types.ExampleBlock{
						Attributes: types.Attributes{
							types.AttrTitle: "example block title",
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
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with caption", func() {
				source := `[caption="a caption "]
====
foo
====`
				expected := types.DocumentFragments{
					types.ExampleBlock{
						Attributes: types.Attributes{
							types.AttrCaption: "a caption ", // trailing space is retained
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
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("example block starting delimiter only", func() {
				source := `====`
				expected := types.DocumentFragments{
					types.ExampleBlock{},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})
		})

		Context("paragraph blocks", func() {

			It("with single rich line", func() {
				source := `[example]
some *example* content`
				expected := types.DocumentFragments{
					types.Paragraph{
						Attributes: types.Attributes{
							types.AttrStyle: types.Example,
						},
						Lines: [][]interface{}{
							{
								types.StringElement{
									Content: "some ",
								},
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{
											Content: "example",
										},
									},
								},
								types.StringElement{
									Content: " content",
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})
		})
	})

	Context("final documents", func() {

		Context("delimited blocks", func() {

			It("with single rich line", func() {
				source := `====
some *example* content
====`
				expected := types.Document{
					Elements: []interface{}{
						types.ExampleBlock{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{
												Content: "some ",
											},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{
														Content: "example",
													},
												},
											},
											types.StringElement{
												Content: " content",
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

			It("with single line starting with a dot", func() {
				source := `====
.foo
====`
				expected := types.Document{
					Elements: []interface{}{
						types.ExampleBlock{
							Elements: []interface{}{},
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
				expected := types.Document{
					Elements: []interface{}{
						types.ExampleBlock{
							Elements: []interface{}{
								types.Paragraph{
									Attributes: types.Attributes{
										types.AttrTitle: "foo",
									},
									Lines: [][]interface{}{
										{
											types.StringElement{
												Content: "some listing code",
											},
										},
										{
											types.StringElement{
												Content: "with ",
											},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{
														Content: "bold content",
													},
												},
											},
										},
									},
								},
								types.BlankLine{},
								types.UnorderedList{
									Items: []types.UnorderedListItem{
										{
											Level:       1,
											BulletStyle: types.OneAsterisk,
											CheckStyle:  types.NoCheck,
											Elements: []interface{}{
												types.Paragraph{
													Lines: [][]interface{}{
														{
															types.StringElement{
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
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with unclosed delimiter", func() {
				source := `====
End of file here`
				expected := types.Document{
					Elements: []interface{}{
						types.ExampleBlock{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{
												Content: "End of file here",
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

			It("with title", func() {
				source := `.example block title
====
foo
====`
				expected := types.Document{
					Elements: []interface{}{
						types.ExampleBlock{
							Attributes: types.Attributes{
								types.AttrTitle: "example block title",
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

			It("example block starting delimiter only", func() {
				source := `====`
				expected := types.Document{
					Elements: []interface{}{
						types.ExampleBlock{
							Elements: []interface{}{},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})

		Context("paragraph blocks", func() {

			It("with single rich line", func() {
				source := `[example]
some *example* content`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle: types.Example,
							},
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "some ",
									},
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{
												Content: "example",
											},
										},
									},
									types.StringElement{
										Content: " content",
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
	})
})
