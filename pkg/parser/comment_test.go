package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("comments", func() {

	Context("in raw documents", func() {

		Context("single line comments", func() {

			It("single line comment alone", func() {
				source := `// A single-line comment.`
				expected := []types.DocumentFragment{
					{
						LineOffset: 1,
						Elements: []interface{}{
							&types.SingleLineComment{
								Content: " A single-line comment.",
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("single line comment at end of line", func() {
				source := `foo // A single-line comment.`
				expected := []types.DocumentFragment{
					{
						LineOffset: 1,
						Elements: []interface{}{

							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "foo // A single-line comment.",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("single line comment within a paragraph", func() {
				source := `a first line
// A single-line comment.
another line // not a comment`
				expected := []types.DocumentFragment{
					{
						LineOffset: 1,
						Elements: []interface{}{

							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "a first line",
										},
									},
									{
										&types.SingleLineComment{
											Content: " A single-line comment.",
										},
									},
									{
										types.StringElement{
											Content: "another line // not a comment",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			Context("invalid", func() {

				It("single line comment with prefixing spaces alone", func() {
					source := `  // A single-line comment.`
					expected := []types.DocumentFragment{
						{
							LineOffset: 1,
							Elements: []interface{}{

								types.LiteralBlock{
									Attributes: types.Attributes{
										types.AttrStyle:            types.Literal,
										types.AttrLiteralBlockType: types.LiteralBlockWithSpacesOnFirstLine,
									},
									Lines: [][]interface{}{
										{
											types.StringElement{
												Content: "  // A single-line comment.",
											},
										},
									},
								},
							},
						},
					}
					result, err := AssembleDocumentFragments(source)
					Expect(err).NotTo(HaveOccurred())
					Expect(result).To(MatchDocumentFragmentGroups(expected))
				})

				It("single line comment with prefixing tabs alone", func() {
					source := "\t\t// A single-line comment."
					expected := []types.DocumentFragment{
						{
							LineOffset: 1,
							Elements: []interface{}{

								types.LiteralBlock{
									Attributes: types.Attributes{
										types.AttrStyle:            types.Literal,
										types.AttrLiteralBlockType: types.LiteralBlockWithSpacesOnFirstLine,
									},
									Lines: [][]interface{}{
										{
											types.StringElement{
												Content: "\t\t// A single-line comment.",
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
				})

				It("single line comment within a paragraph with tab", func() {
					source := `a first line
	// A single-line comment.
another line`
					expected := []types.DocumentFragment{
						{
							LineOffset: 1,
							Elements: []interface{}{

								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{
												Content: "a first line",
											},
										},
										{
											types.StringElement{
												Content: "\t// A single-line comment.",
											},
										},
										{
											types.StringElement{
												Content: "another line",
											},
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

		Context("comment blocks", func() {

			It("comment block alone", func() {
				source := `//// 
a *comment* block
with multiple lines
////`
				expected := []types.DocumentFragment{
					{
						LineOffset: 1,
						Elements: []interface{}{

							types.CommentBlock{
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "a *comment* block",
										},
									},
									{
										types.StringElement{
											Content: "with multiple lines",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("comment block with paragraphs around", func() {
				source := `a first paragraph

//// 
a *comment* block
with multiple lines
////
a second paragraph`
				expected := []types.DocumentFragment{
					{
						LineOffset: 1,
						Elements: []interface{}{

							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "a first paragraph",
										},
									},
								},
							},
							&types.BlankLine{}, // blankline is required between a paragraph and the next block
							types.CommentBlock{
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "a *comment* block",
										},
									},
									{
										types.StringElement{
											Content: "with multiple lines",
										},
									},
								},
							},
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "a second paragraph",
										},
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

		Context("single line comments", func() {

			It("single line comment alone", func() {
				source := `// A single-line comment.`
				expected := types.Document{
					Elements: []interface{}{},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("single line comment with prefixing spaces alone", func() {
				source := `  // A single-line comment.`
				expected := types.Document{
					Elements: []interface{}{
						types.LiteralBlock{
							Attributes: types.Attributes{
								types.AttrStyle:            types.Literal,
								types.AttrLiteralBlockType: types.LiteralBlockWithSpacesOnFirstLine,
							},
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "  // A single-line comment.",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("single line comment with prefixing tabs alone", func() {
				source := "\t\t// A single-line comment."
				expected := types.Document{
					Elements: []interface{}{
						types.LiteralBlock{
							Attributes: types.Attributes{
								types.AttrStyle:            types.Literal,
								types.AttrLiteralBlockType: types.LiteralBlockWithSpacesOnFirstLine,
							},
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "\t\t// A single-line comment.",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("single line comment at end of line", func() {
				source := `foo // A single-line comment.`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "foo // A single-line comment."},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("single line comment within a paragraph", func() {
				source := `a first line
// A single-line comment.
another line`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a first line"},
								},
								{
									types.StringElement{Content: "another line"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			Context("invalid", func() {

				It("single line comment within a paragraph with tab", func() {
					source := `a first line
	// A single-line comment.
another line`
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a first line"},
									},
									{
										types.StringElement{Content: "\t// A single-line comment."},
									},
									{
										types.StringElement{Content: "another line"},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})
			})
		})

		Context("comment blocks", func() {

			It("comment block alone", func() {
				source := `//// 
a *comment* block
with multiple lines
////`
				expected := types.Document{
					Elements: []interface{}{},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("comment block with paragraphs around", func() {
				source := `a first paragraph
			
//// 
a *comment* block
with multiple lines
////
a second paragraph`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a first paragraph"},
								},
							},
						},
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a second paragraph"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})

		It("comment in section", func() {
			source := `== section 1

a first paragraph

//// 
a *comment* block
with multiple lines
////
a second paragraph`
			section1Title := []interface{}{
				types.StringElement{
					Content: "section 1",
				},
			}
			expected := types.Document{
				ElementReferences: types.ElementReferences{
					"_section_1": section1Title,
				},
				Elements: []interface{}{
					types.Section{
						Attributes: types.Attributes{
							types.AttrID: "_section_1",
						},
						Level: 1,
						Title: section1Title,
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a first paragraph"},
									},
								},
							},
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a second paragraph"},
									},
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("comment in preamble", func() {
			source := `= section 0

//// 
a *comment* block
with multiple lines
////

== section 1

a first paragraph

a second paragraph`
			section0Title := []interface{}{
				types.StringElement{
					Content: "section 0",
				},
			}
			section1Title := []interface{}{
				types.StringElement{
					Content: "section 1",
				},
			}
			expected := types.Document{
				ElementReferences: types.ElementReferences{
					"_section_0": section0Title,
					"_section_1": section1Title,
				},
				Elements: []interface{}{
					types.Section{
						Attributes: types.Attributes{
							types.AttrID: "_section_0",
						},
						Level: 0,
						Title: section0Title,
						Elements: []interface{}{
							types.Section{
								Attributes: types.Attributes{
									types.AttrID: "_section_1",
								},
								Level: 1,
								Title: section1Title,
								Elements: []interface{}{
									types.Paragraph{
										Lines: [][]interface{}{
											{
												types.StringElement{Content: "a first paragraph"},
											},
										},
									},
									types.Paragraph{
										Lines: [][]interface{}{
											{
												types.StringElement{Content: "a second paragraph"},
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
	})

})
