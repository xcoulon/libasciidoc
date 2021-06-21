package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("quote blocks", func() {

	Context("in final documents", func() {

		Context("delimited blocks", func() {

			It("single-line quote block with author and title", func() {
				source := `[quote, john doe, quote title]
____
some *quote* content
____`
				expected := types.Document{
					Elements: []interface{}{
						types.QuoteBlock{
							Attributes: types.Attributes{
								types.AttrStyle:       types.Quote,
								types.AttrQuoteAuthor: "john doe",
								types.AttrQuoteTitle:  "quote title",
							},
							Elements: []interface{}{
								&types.Paragraph{
									Elements: []interface{}{
										types.StringElement{
											Content: "some ",
										},
										&types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{
													Content: "quote",
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
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("multi-line quote with author only", func() {
				source := `[quote, john doe,   ]
____
- some 
- quote 
- content 
____
`
				expected := types.Document{
					Elements: []interface{}{
						types.QuoteBlock{
							Attributes: types.Attributes{
								types.AttrStyle:       types.Quote,
								types.AttrQuoteAuthor: "john doe",
							},
							Elements: []interface{}{
								&types.GenericList{
									Kind: types.UnorderedListKind,
									Elements: []types.ListElement{
										&types.UnorderedListElement{
											BulletStyle: types.Dash,
											CheckStyle:  types.NoCheck,
											Elements: []interface{}{
												&types.Paragraph{
													Elements: []interface{}{
														types.StringElement{
															Content: "some ",
														},
													},
												},
											},
										},
										&types.UnorderedListElement{
											BulletStyle: types.Dash,
											CheckStyle:  types.NoCheck,
											Elements: []interface{}{
												&types.Paragraph{
													Elements: []interface{}{
														types.StringElement{
															Content: "quote ",
														},
													},
												},
											},
										},
										&types.UnorderedListElement{
											BulletStyle: types.Dash,
											CheckStyle:  types.NoCheck,
											Elements: []interface{}{
												&types.Paragraph{
													Elements: []interface{}{
														types.StringElement{
															Content: "content ",
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

			It("single-line quote with title only", func() {
				source := `[quote, ,quote title]
____
some quote content 
____
`
				expected := types.Document{
					Elements: []interface{}{
						types.QuoteBlock{
							Attributes: types.Attributes{
								types.AttrStyle:      types.Quote,
								types.AttrQuoteTitle: "quote title",
							},
							Elements: []interface{}{
								&types.Paragraph{
									Elements: []interface{}{
										types.StringElement{
											Content: "some quote content ",
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
				source := `[quote]
____
.foo
____`
				expected := types.Document{
					Elements: []interface{}{
						types.QuoteBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Quote,
							},
							Elements: []interface{}{},
						},
					},
				}
				result, err := ParseDocument(source)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(MatchDocument(expected))
			})

			It("multi-line quote with rendered lists and block and without author and title", func() {
				source := `[quote]
____
* some
----
* quote 
----
* content
____`
				expected := types.Document{
					Elements: []interface{}{
						types.QuoteBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Quote,
							},
							Elements: []interface{}{
								&types.GenericList{
									Kind: types.UnorderedListKind,
									Elements: []types.ListElement{
										&types.UnorderedListElement{
											BulletStyle: types.OneAsterisk,
											CheckStyle:  types.NoCheck,
											Elements: []interface{}{
												&types.Paragraph{
													Elements: []interface{}{
														types.StringElement{
															Content: "some",
														},
													},
												},
											},
										},
									},
								},
								types.ListingBlock{
									Elements: []interface{}{
										types.StringElement{
											Content: "* quote ",
										},
									},
								},
								&types.GenericList{
									Kind: types.UnorderedListKind,
									Elements: []types.ListElement{
										&types.UnorderedListElement{
											BulletStyle: types.OneAsterisk,
											CheckStyle:  types.NoCheck,
											Elements: []interface{}{
												&types.Paragraph{
													Elements: []interface{}{
														types.StringElement{
															Content: "content",
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

			It("multi-line quote with rendered list and without author and title", func() {
				source := `[quote]
____
* some


* quote 


* content
____`
				expected := types.Document{
					Elements: []interface{}{
						types.QuoteBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Quote,
							},
							Elements: []interface{}{
								&types.GenericList{
									Kind: types.UnorderedListKind,
									Elements: []types.ListElement{
										&types.UnorderedListElement{
											BulletStyle: types.OneAsterisk,
											CheckStyle:  types.NoCheck,
											Elements: []interface{}{
												&types.Paragraph{
													Elements: []interface{}{
														types.StringElement{
															Content: "some",
														},
													},
												},
											},
										},
										&types.UnorderedListElement{
											BulletStyle: types.OneAsterisk,
											CheckStyle:  types.NoCheck,
											Elements: []interface{}{
												&types.Paragraph{
													Elements: []interface{}{
														types.StringElement{
															Content: "quote ",
														},
													},
												},
											},
										},
										&types.UnorderedListElement{
											BulletStyle: types.OneAsterisk,
											CheckStyle:  types.NoCheck,
											Elements: []interface{}{
												&types.Paragraph{
													Elements: []interface{}{
														types.StringElement{
															Content: "content",
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

			It("empty quote without author and title", func() {
				source := `[quote]
____
____`
				expected := types.Document{
					Elements: []interface{}{
						types.QuoteBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Quote,
							},
							Elements: []interface{}{},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("unclosed quote without author and title", func() {
				source := `[quote]
____
foo
`
				expected := types.Document{
					Elements: []interface{}{
						types.QuoteBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Quote,
							},
							Elements: []interface{}{
								&types.Paragraph{
									Elements: []interface{}{
										types.StringElement{
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
		})
	})
})
