package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("element attributes", func() {

	Context("in final documents", func() {

		Context("element links", func() {

			Context("with valid syntax", func() {

				It("element link alone", func() {
					source := `[link=http://foo.bar]
a paragraph`
					expected := &types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									"link": "http://foo.bar",
								},
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "a paragraph",
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})
				It("spaces in link", func() {
					source := `[link= http://foo.bar  ]
a paragraph`
					expected := &types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									"link": "http://foo.bar",
								},
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "a paragraph",
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})
			})

			Context("with invalid syntax", func() {

				It("spaces before keyword", func() {
					source := `[ link=http://foo.bar]
a paragraph`
					expected := &types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "[ link=http://foo.bar]",
										},
									},
									{
										&types.StringElement{
											Content: "a paragraph",
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("unbalanced brackets", func() {
					source := `[link=http://foo.bar
a paragraph`
					expected := &types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "[link=http://foo.bar",
										},
									},
									{
										&types.StringElement{
											Content: "a paragraph",
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

		Context("element ids", func() {

			Context("with valid syntax", func() {

				It("normal syntax", func() {
					source := `[[img-foobar]]
a paragraph`
					expected := &types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrID: "img-foobar",
								},
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "a paragraph",
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("short-hand syntax", func() {
					source := `[#img-foobar]
a paragraph`
					expected := &types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrID: "img-foobar",
								},
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "a paragraph",
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})
			})

			Context("with invalid syntax", func() {

				It("extra spaces", func() {
					source := `[ #img-foobar ]
a paragraph`
					expected := &types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "[ #img-foobar ]",
										},
									},
									{
										&types.StringElement{
											Content: "a paragraph",
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("unbalanced brackets", func() {
					source := `[#img-foobar
a paragraph`
					expected := &types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "[#img-foobar",
										},
									},
									{
										&types.StringElement{
											Content: "a paragraph",
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

		Context("element title", func() {

			Context("with valid syntax", func() {

				It("valid element title", func() {
					source := `.a title
a paragraph`
					expected := &types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrTitle: "a title",
								},
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "a paragraph",
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})
			})

			Context("with invalid syntax", func() {

				It("extra space after dot", func() {
					source := `. a title
a list item!`
					expected := &types.Document{
						Elements: []interface{}{
							types.OrderedList{
								Items: []*types.OrderedListElement{
									{
										Style: types.Arabic,
										Elements: []interface{}{
											types.Paragraph{
												Lines: [][]interface{}{
													{
														&types.StringElement{
															Content: "a title",
														},
													},
													{
														&types.StringElement{
															Content: "a list item!",
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

				It("not a dot", func() {
					source := `!a title
a paragraph`

					expected := &types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "!a title",
										},
									},
									{
										&types.StringElement{
											Content: "a paragraph",
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

		Context("element roles", func() {

			Context("with valid syntax", func() {

				It("shortcut role element", func() {
					source := `[.a_role]
a paragraph`
					expected := &types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrRoles: []interface{}{"a_role"},
								},
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "a paragraph",
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("full role syntax", func() {
					source := `[role=a_role]
a paragraph`
					expected := &types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrRoles: []interface{}{"a_role"},
								},
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "a paragraph",
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})
			})

			It("blank line after role attribute", func() {
				source := `[.a_role]

a paragraph`
				expected := &types.Document{
					Elements: []interface{}{
						// attribute was attached to blankline, which was filtered out in the final doc
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{
									Content: "a paragraph",
								},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("blank lines after id, role and title attributes", func() {
				source := `[.a_role]
[[ID]]
.title


a paragraph`
				expected := &types.Document{
					Elements: []interface{}{
						// attributes were attached to blankline, which was filtered out in the final doc
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{
									Content: "a paragraph",
								},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})

		Context("standalone attributes", func() {

			It("single standalone attribute", func() {
				source := `[.a_role]
`
				expected := &types.Document{
					Elements: []interface{}{},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("multiple standalone attributes", func() {
				source := `[.a_role]
[[ID]]
.title`
				expected := &types.Document{
					Elements: []interface{}{},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("multiple standalone attributes after a paragraph", func() {
				source := `a paragraph
			
[.a_role]
[[ID]]
.title`
				expected := &types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									&types.StringElement{
										Content: "a paragraph",
									},
								},
							},
						},
						// everything after the paragraph (blankline and standalone attributes)
						// is ignored in the final doc
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})
	})
})
