package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("headers", func() {

	Context("in raw documents", func() {

		Context("valid cases", func() {

			It("header alone", func() {
				source := `= Document Title
			
This journey continues.`

				expected := types.DocumentFragments{
					types.Section{
						Level: 0,
						Title: []interface{}{
							types.StringElement{
								Content: "Document Title",
							},
						},
					},
					types.BlankLine{},
					types.InlineElements{
						types.StringElement{Content: "This journey continues."},
					},
				}
				result, err := ParseRawSource(source)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(MatchDocumentFragments(expected))
			})

			Context("document authors", func() {

				Context("single author", func() {

					It("all author data with extra spaces", func() {
						source := `= title
John  Foo    Doe  <johndoe@example.com>`
						expected := types.DocumentFragments{
							types.Section{
								Level: 0,
								Title: []interface{}{
									types.StringElement{
										Content: "title",
									},
								},
							},
							[]types.DocumentAuthor{
								{
									FullName: "John  Foo    Doe  ",
									Email:    "johndoe@example.com",
								},
							},
						}
						result, err := ParseRawSource(source)
						Expect(err).NotTo(HaveOccurred())
						Expect(result).To(MatchDocumentFragments(expected))
					})

					It("lastname with underscores", func() {
						source := `= title
Jane the_Doe <jane@example.com>`
						expected := types.DocumentFragments{
							types.Section{
								Level: 0,
								Title: []interface{}{
									types.StringElement{
										Content: "title",
									},
								},
							},
							[]types.DocumentAuthor{
								{
									FullName: "Jane the_Doe ",
									Email:    "jane@example.com",
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("with middlename and composed lastname", func() {
						source := `= title
Jane Foo the Doe <jane@example.com>`
						expected := types.DocumentFragments{
							types.Section{
								Level: 0,
								Title: []interface{}{
									types.StringElement{
										Content: "title",
									},
								},
							},
							[]types.DocumentAuthor{
								{
									FullName: "Jane Foo the Doe ",
									Email:    "jane@example.com",
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("firstname and lastname only", func() {
						source := `= title
John Doe`
						expected := types.DocumentFragments{
							types.Section{
								Level: 0,
								Title: []interface{}{
									types.StringElement{
										Content: "title",
									},
								},
							},
							[]types.DocumentAuthor{
								{
									FullName: "John Doe",
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("firstname only", func() {
						source := `= title
Doe`
						expected := types.DocumentFragments{
							types.Section{
								Level: 0,
								Title: []interface{}{
									types.StringElement{
										Content: "title",
									},
								},
							},
							[]types.DocumentAuthor{
								{
									FullName: "Doe",
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("alternate author input", func() {
						source := `= title
:author: John Foo Doe` // `:"email":` is processed as a regular attribute
						expected := types.DocumentFragments{
							types.Section{
								Level: 0,
								Title: []interface{}{
									types.StringElement{
										Content: "title",
									},
								},
							},
							types.AttributeDeclaration{
								Name:  "author",
								Value: "John Foo Doe",
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})
				})

				Context("multiple authors", func() {

					It("2 authors", func() {
						source := `= title
John  Foo Doe  <johndoe@example.com>; Jane the_Doe <jane@example.com>`
						expected := types.DocumentFragments{
							types.Section{
								Level: 0,
								Title: []interface{}{
									types.StringElement{
										Content: "title",
									},
								},
							},
							[]types.DocumentAuthor{
								{
									FullName: "John  Foo Doe  ",
									Email:    "johndoe@example.com",
								},
								{
									FullName: "Jane the_Doe ",
									Email:    "jane@example.com",
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})
				})

				Context("authors and comments", func() {

					It("authors commented out", func() {
						source := `= title
					// John  Foo Doe  <johndoe@example.com>; Jane the_Doe <jane@example.com>`
						expected := types.DocumentFragments{
							types.Section{
								Level: 0,
								Title: []interface{}{
									types.StringElement{
										Content: "title",
									},
								},
							},
							types.SingleLineComment{
								Content: ` John  Foo Doe  <johndoe@example.com>; Jane the_Doe <jane@example.com>`,
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("authors after a single comment line", func() {
						source := `= title
					// a comment
					John  Foo Doe  <johndoe@example.com>; Jane the_Doe <jane@example.com>`
						expected := types.DocumentFragments{
							types.Section{
								Level: 0,
								Title: []interface{}{
									types.StringElement{
										Content: "title",
									},
								},
							},
							types.SingleLineComment{
								Content: " a comment",
							},
							[]types.DocumentAuthor{
								{
									FullName: "John  Foo Doe  ",
									Email:    "johndoe@example.com",
								},
								{
									FullName: "Jane the_Doe ",
									Email:    "jane@example.com",
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("authors after a comment block", func() {
						source := `= title
//// 
a comment
////
John  Foo Doe  <johndoe@example.com>; Jane the_Doe <jane@example.com>`
						expected := types.DocumentFragments{
							types.Section{
								Level: 0,
								Title: []interface{}{
									types.StringElement{
										Content: "title",
									},
								},
							},
							types.BlockDelimiter{
								Kind: types.Comment,
							},
							types.RawLine("a comment"),
							types.BlockDelimiter{
								Kind: types.Comment,
							},
							[]types.DocumentAuthor{
								{
									FullName: "John  Foo Doe  ",
									Email:    "johndoe@example.com",
								},
								{
									FullName: "Jane the_Doe ",
									Email:    "jane@example.com",
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})
				})
			})

			Context("document revisions", func() {

				It("full document revision without any comment", func() {
					source := `= title
				John Doe
				v1.0, March 29, 2020: Updated revision`
					expected := types.DocumentFragments{
						types.Section{
							Level: 0,
							Title: []interface{}{
								types.StringElement{
									Content: "title",
								},
							},
						},
						[]types.DocumentAuthor{
							{
								FullName: "John Doe",
							},
						},
						types.DocumentRevision{
							Revnumber: "1.0",
							Revdate:   "March 29, 2020",
							Revremark: "Updated revision",
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("full document revision with a comment before author", func() {
					source := `= title
				// a comment
				John Doe
				v1.0, March 29, 2020: Updated revision`
					expected := types.DocumentFragments{
						types.Section{
							Level: 0,
							Title: []interface{}{
								types.StringElement{
									Content: "title",
								},
							},
						},
						[]types.DocumentAuthor{
							{
								FullName: "John Doe",
							},
						},
						types.DocumentRevision{
							Revnumber: "1.0",
							Revdate:   "March 29, 2020",
							Revremark: "Updated revision",
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("full document revision with a comment before revision", func() {
					source := `= title
				John Doe
				// a comment
				v1.0, March 29, 2020: Updated revision`
					expected := types.DocumentFragments{
						types.Section{
							Level: 0,
							Title: []interface{}{
								types.StringElement{
									Content: "title",
								},
							},
						},
						[]types.DocumentAuthor{
							{
								FullName: "John Doe",
							},
						},
						types.DocumentRevision{
							Revnumber: "1.0",
							Revdate:   "March 29, 2020",
							Revremark: "Updated revision",
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("revision with revnumber and revdate only", func() {
					source := `= title
				John Doe
				v1.0, March 29, 2020`
					expected := types.DocumentFragments{
						types.Section{
							Level: 0,
							Title: []interface{}{
								types.StringElement{
									Content: "title",
								},
							},
						},
						[]types.DocumentAuthor{
							{
								FullName: "John Doe",
							},
						},
						types.DocumentRevision{
							Revnumber: "1.0",
							Revdate:   "March 29, 2020",
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("revision with revnumber and revdate - with colon separator", func() {
					source := `= title
				John Doe
				1.0, March 29, 2020:`
					expected := types.DocumentFragments{
						types.Section{
							Level: 0,
							Title: []interface{}{
								types.StringElement{
									Content: "title",
								},
							},
						},
						[]types.DocumentAuthor{
							{
								FullName: "John Doe",
							},
						},
						types.DocumentRevision{
							Revnumber: "1.0",
							Revdate:   "March 29, 2020",
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})
				It("revision with revnumber only - comma suffix", func() {
					source := `= title
				John Doe
				1.0,`
					expected := types.DocumentFragments{
						types.Section{
							Level: 0,
							Title: []interface{}{
								types.StringElement{
									Content: "title",
								},
							},
						},
						[]types.DocumentAuthor{
							{
								FullName: "John Doe",
							},
						},
						types.DocumentRevision{
							Revnumber: "1.0",
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("revision with revdate as number - spaces and no prefix no suffix", func() {
					source := `= title
				John Doe
				1.0`
					expected := types.DocumentFragments{
						types.Section{
							Level: 0,
							Title: []interface{}{
								types.StringElement{
									Content: "title",
								},
							},
							Elements: []interface{}{},
						},
						[]types.DocumentAuthor{
							{
								FullName: "John Doe",
							},
						},
						types.DocumentRevision{
							Revdate: "1.0",
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("revision with revdate as alphanum - spaces and no prefix no suffix", func() {
					source := `= title
				John Doe
				1.0a`
					expected := types.DocumentFragments{
						types.Section{
							Level: 0,
							Title: []interface{}{
								types.StringElement{
									Content: "title",
								},
							},
							Elements: []interface{}{},
						},
						[]types.DocumentAuthor{
							{
								FullName: "John Doe",
							},
						},
						types.DocumentRevision{
							Revdate: "1.0a",
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("revision with revnumber only", func() {
					source := `= title
				John Doe
				v1.0:`
					expected := types.DocumentFragments{
						types.Section{
							Level: 0,
							Title: []interface{}{
								types.StringElement{
									Content: "title",
								},
							},
						},
						[]types.DocumentAuthor{
							{
								FullName: "John Doe",
							},
						},
						types.DocumentRevision{
							Revnumber: "1.0",
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("revision with spaces and capital revnumber ", func() {
					source := `= title
				John Doe
				V1.0:`
					expected := types.DocumentFragments{
						types.Section{
							Level: 0,
							Title: []interface{}{
								types.StringElement{
									Content: "title",
								},
							},
							Elements: []interface{}{},
						},
						[]types.DocumentAuthor{
							{
								FullName: "John Doe",
							},
						},
						types.DocumentRevision{
							Revnumber: "1.0",
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("revision only - with comma separator", func() {
					source := `= title
				John Doe
				v1.0,`
					expected := types.DocumentFragments{
						types.Section{
							Level: 0,
							Title: []interface{}{
								types.StringElement{
									Content: "title",
								},
							},
						},
						[]types.DocumentAuthor{
							{
								FullName: "John Doe",
							},
						},
						types.DocumentRevision{
							Revnumber: "1.0",
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("revision with revnumber plus comma and colon separators", func() {
					source := `= title
				John Doe
				v1.0,:`
					expected := types.DocumentFragments{
						types.Section{
							Level: 0,
							Title: []interface{}{
								types.StringElement{
									Content: "title",
								},
							},
						},
						[]types.DocumentAuthor{
							{
								FullName: "John Doe",
							},
						},
						types.DocumentRevision{
							Revnumber: "1.0",
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("revision with revnumber and empty revremark", func() {
					source := `= title
John Doe
v1.0:`
					expected := types.DocumentFragments{
						types.Section{
							Level: 0,
							Title: []interface{}{
								types.StringElement{
									Content: "title",
								},
							},
						},
						[]types.DocumentAuthor{
							{
								FullName: "John Doe",
							},
						},
						types.DocumentRevision{
							Revnumber: "1.0",
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

			})

			Context("document header attributes", func() {

				It("valid attribute names", func() {
					source := `:a:
:author: Xavier
:_auth0r: Xavier`
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "a",
							Value: nil,
						},
						types.AttributeDeclaration{
							Name:  "author",
							Value: "Xavier",
						},
						types.AttributeDeclaration{
							Name:  "_auth0r",
							Value: "Xavier",
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("attributes and paragraph without blank line in-between", func() {
					source := `:toc:
:date:  2017-01-01
:author: Xavier
:hardbreaks:
a paragraph`
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "toc",
							Value: nil,
						},
						types.AttributeDeclaration{
							Name:  "date",
							Value: "2017-01-01",
						},
						types.AttributeDeclaration{
							Name:  "author",
							Value: "Xavier",
						},
						types.AttributeDeclaration{
							Name:  "hardbreaks",
							Value: nil,
						},
						types.InlineElements{
							types.StringElement{
								Content: "a paragraph",
							},
						},
					}
					result, err := ParseRawSource(source)
					Expect(err).NotTo(HaveOccurred())
					Expect(result).To(MatchDocumentFragments(expected))
				})

				It("contiguous attributes and paragraph with blank line in-between", func() {
					source := `:toc:
:date: 2017-01-01
:author: Xavier
:hardbreaks:

a paragraph`
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "toc",
							Value: nil,
						},
						types.AttributeDeclaration{
							Name:  "date",
							Value: "2017-01-01",
						},
						types.AttributeDeclaration{
							Name:  "author",
							Value: "Xavier",
						},
						types.AttributeDeclaration{
							Name:  "hardbreaks",
							Value: nil,
						},
						types.BlankLine{},
						types.InlineElements{
							types.StringElement{
								Content: "a paragraph",
							},
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("splitted attributes and paragraph with blank line in-between", func() {
					source := `:toc:
:date: 2017-01-01

:author: Xavier

:hardbreaks:

a paragraph`
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "toc",
							Value: nil,
						},
						types.AttributeDeclaration{
							Name:  "date",
							Value: "2017-01-01",
						},
						types.BlankLine{},
						types.AttributeDeclaration{
							Name:  "author",
							Value: "Xavier",
						},
						types.BlankLine{},
						types.AttributeDeclaration{
							Name:  "hardbreaks",
							Value: nil,
						},
						types.BlankLine{},
						types.InlineElements{
							types.StringElement{
								Content: "a paragraph",
							},
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("no header and attributes in body", func() {
					source := `a paragraph
	
:toc:
:date: 2017-01-01
:author: Xavier
:hardbreaks:`
					expected := types.DocumentFragments{
						types.InlineElements{
							types.StringElement{
								Content: "a paragraph",
							},
						},
						types.BlankLine{},
						types.AttributeDeclaration{
							Name:  "toc",
							Value: nil,
						},
						types.AttributeDeclaration{
							Name:  "date",
							Value: "2017-01-01",
						},
						types.AttributeDeclaration{
							Name:  "author",
							Value: "Xavier",
						},
						types.AttributeDeclaration{
							Name:  "hardbreaks",
							Value: nil,
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})
			})
		})

		Context("invalid cases", func() {

			It("paragraph without blank line before attribute declarations", func() {
				source := `a paragraph
:toc:
:date: 2017-01-01
:author: Xavier`
				expected := types.DocumentFragments{
					types.InlineElements{
						types.StringElement{Content: "a paragraph"},
					},
					types.InlineElements{
						types.StringElement{Content: ":toc:"},
					},
					types.InlineElements{
						types.StringElement{Content: ":date: 2017-01-01"},
					},
					types.InlineElements{
						types.StringElement{Content: ":author: Xavier"},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("invalid attribute names", func() {
				source := `:@date: 2017-01-01
:{author}: Xavier`
				expected := types.DocumentFragments{
					types.InlineElements{
						types.StringElement{Content: ":@date: 2017-01-01"},
					},
					types.InlineElements{
						types.StringElement{Content: ":{author}: Xavier"},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})
		})
	})

	Context("in final documents", func() {

		Context("valid cases", func() {

			It("header alone", func() {
				source := `= Document Title
			
This journey continues.`

				title := []interface{}{
					types.StringElement{Content: "Document Title"},
				}
				expected := types.Document{
					ElementReferences: types.ElementReferences{
						"_document_title": title,
					},
					Elements: []interface{}{
						types.Section{
							Level: 0,
							Attributes: types.Attributes{
								types.AttrID: "_document_title",
							},
							Title: title,
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "This journey continues."},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			Context("document authors", func() {

				Context("single author", func() {

					It("all author data with extra spaces", func() {
						source := `= title
John  Foo    Doe  <johndoe@example.com>`
						title := []interface{}{
							types.StringElement{
								Content: "title",
							},
						}
						expected := types.Document{
							Attributes: types.Attributes{
								types.AttrAuthors: []types.DocumentAuthor{
									{
										FullName: "John Foo Doe",
										Email:    "johndoe@example.com",
									},
								},
								"firstname":      "John",
								"middlename":     "Foo",
								"lastname":       "Doe",
								"author":         "John Foo Doe",
								"authorinitials": "JFD",
								"email":          "johndoe@example.com",
							},
							ElementReferences: types.ElementReferences{
								"_title": title,
							},
							Elements: []interface{}{
								types.Section{
									Level: 0,
									Attributes: types.Attributes{
										types.AttrID: "_title",
									},
									Title:    title,
									Elements: []interface{}{},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("lastname with underscores", func() {
						source := `= title
Jane the_Doe <jane@example.com>`
						title := []interface{}{
							types.StringElement{
								Content: "title",
							},
						}
						expected := types.Document{
							Attributes: types.Attributes{
								types.AttrAuthors: []types.DocumentAuthor{
									{
										FullName: "Jane the Doe",
										Email:    "jane@example.com",
									},
								},
								"firstname":      "Jane",
								"lastname":       "the Doe",
								"author":         "Jane the Doe",
								"authorinitials": "Jt",
								"email":          "jane@example.com",
							},
							ElementReferences: types.ElementReferences{
								"_title": title,
							},
							Elements: []interface{}{
								types.Section{
									Level: 0,
									Attributes: types.Attributes{
										types.AttrID: "_title",
									},
									Title:    title,
									Elements: []interface{}{},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("with middlename and composed lastname", func() {
						source := `= title
Jane Foo the Doe <jane@example.com>`
						title := []interface{}{
							types.StringElement{
								Content: "title",
							},
						}
						expected := types.Document{
							Attributes: types.Attributes{
								types.AttrAuthors: []types.DocumentAuthor{
									{
										FullName: "Jane Foo the Doe",
										Email:    "jane@example.com",
									},
								},
								"firstname":      "Jane",
								"middlename":     "Foo",
								"lastname":       "the Doe",
								"author":         "Jane Foo the Doe",
								"authorinitials": "JFt",
								"email":          "jane@example.com",
							},
							ElementReferences: types.ElementReferences{
								"_title": title,
							},
							Elements: []interface{}{
								types.Section{
									Level: 0,
									Attributes: types.Attributes{
										types.AttrID: "_title",
									},
									Title:    title,
									Elements: []interface{}{},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("firstname and lastname only", func() {
						source := `= title
John Doe`
						title := []interface{}{
							types.StringElement{
								Content: "title",
							},
						}
						expected := types.Document{
							Attributes: types.Attributes{
								types.AttrAuthors: []types.DocumentAuthor{
									{
										FullName: "John Doe",
									},
								},
								"firstname":      "John",
								"lastname":       "Doe",
								"author":         "John Doe",
								"authorinitials": "JD",
							},
							ElementReferences: types.ElementReferences{
								"_title": title,
							},
							Elements: []interface{}{
								types.Section{
									Level: 0,
									Attributes: types.Attributes{
										types.AttrID: "_title",
									},
									Title:    title,
									Elements: []interface{}{},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("firstname only", func() {
						source := `= title
Doe`
						title := []interface{}{
							types.StringElement{
								Content: "title",
							},
						}
						expected := types.Document{
							Attributes: types.Attributes{
								types.AttrAuthors: []types.DocumentAuthor{
									{
										FullName: "Doe",
									},
								},
								"firstname":      "Doe",
								"author":         "Doe",
								"authorinitials": "D",
							},
							ElementReferences: types.ElementReferences{
								"_title": title,
							},
							Elements: []interface{}{
								types.Section{
									Level: 0,
									Attributes: types.Attributes{
										types.AttrID: "_title",
									},
									Title:    title,
									Elements: []interface{}{},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("alternate author input", func() {
						source := `= title
:author: John Foo Doe` // `:"email":` is processed as a regular attribute
						title := []interface{}{
							types.StringElement{
								Content: "title",
							},
						}
						expected := types.Document{
							Attributes: types.Attributes{
								types.AttrAuthors: []types.DocumentAuthor{
									{
										FullName: "John Foo Doe",
									},
								},
								"firstname":      "John",
								"middlename":     "Foo",
								"lastname":       "Doe",
								"author":         "John Foo Doe",
								"authorinitials": "JFD",
							},
							ElementReferences: types.ElementReferences{
								"_title": title,
							},
							Elements: []interface{}{
								types.Section{
									Level: 0,
									Attributes: types.Attributes{
										types.AttrID: "_title",
									},
									Title:    title,
									Elements: []interface{}{},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})
				})

				Context("multiple authors", func() {

					It("2 authors", func() {
						source := `= title
John  Foo Doe  <johndoe@example.com>; Jane the_Doe <jane@example.com>`
						title := []interface{}{
							types.StringElement{
								Content: "title",
							},
						}
						expected := types.Document{
							Attributes: types.Attributes{
								types.AttrAuthors: []types.DocumentAuthor{
									{
										FullName: "John Foo Doe",
										Email:    "johndoe@example.com",
									},
									{
										FullName: "Jane the Doe",
										Email:    "jane@example.com",
									},
								},
								"firstname":        "John",
								"middlename":       "Foo",
								"lastname":         "Doe",
								"author":           "John Foo Doe",
								"authorinitials":   "JFD",
								"email":            "johndoe@example.com",
								"firstname_2":      "Jane",
								"lastname_2":       "the Doe",
								"author_2":         "Jane the Doe",
								"authorinitials_2": "Jt",
								"email_2":          "jane@example.com",
							},
							ElementReferences: types.ElementReferences{
								"_title": title,
							},
							Elements: []interface{}{
								types.Section{
									Level: 0,
									Attributes: types.Attributes{
										types.AttrID: "_title",
									},
									Title:    title,
									Elements: []interface{}{},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})
				})

				Context("authors and comments", func() {

					It("authors commented out", func() {
						source := `= title
					// John  Foo Doe  <johndoe@example.com>; Jane the_Doe <jane@example.com>`
						title := []interface{}{
							types.StringElement{
								Content: "title",
							},
						}
						expected := types.Document{
							ElementReferences: types.ElementReferences{
								"_title": title,
							},
							Elements: []interface{}{
								types.Section{
									Level: 0,
									Attributes: types.Attributes{
										types.AttrID: "_title",
									},
									Title:    title,
									Elements: []interface{}{},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("authors after a single comment line", func() {
						source := `= title
					// a comment
					John  Foo Doe  <johndoe@example.com>; Jane the_Doe <jane@example.com>`
						title := []interface{}{
							types.StringElement{
								Content: "title",
							},
						}
						expected := types.Document{
							Attributes: types.Attributes{
								types.AttrAuthors: []types.DocumentAuthor{
									{
										FullName: "John Foo Doe",
										Email:    "johndoe@example.com",
									},
									{
										FullName: "Jane the Doe",
										Email:    "jane@example.com",
									},
								},
								"firstname":        "John",
								"middlename":       "Foo",
								"lastname":         "Doe",
								"author":           "John Foo Doe",
								"authorinitials":   "JFD",
								"email":            "johndoe@example.com",
								"firstname_2":      "Jane",
								"lastname_2":       "the Doe",
								"author_2":         "Jane the Doe",
								"authorinitials_2": "Jt",
								"email_2":          "jane@example.com",
							},
							ElementReferences: types.ElementReferences{
								"_title": title,
							},
							Elements: []interface{}{
								types.Section{
									Level: 0,
									Attributes: types.Attributes{
										types.AttrID: "_title",
									},
									Title:    title,
									Elements: []interface{}{},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("authors after a comment block", func() {
						source := `= title
//// 
a comment
////
John  Foo Doe  <johndoe@example.com>; Jane the_Doe <jane@example.com>`
						title := []interface{}{
							types.StringElement{
								Content: "title",
							},
						}
						expected := types.Document{
							Attributes: types.Attributes{
								types.AttrAuthors: []types.DocumentAuthor{
									{
										FullName: "John Foo Doe",
										Email:    "johndoe@example.com",
									},
									{
										FullName: "Jane the Doe",
										Email:    "jane@example.com",
									},
								},
								"firstname":        "John",
								"middlename":       "Foo",
								"lastname":         "Doe",
								"author":           "John Foo Doe",
								"authorinitials":   "JFD",
								"email":            "johndoe@example.com",
								"firstname_2":      "Jane",
								"lastname_2":       "the Doe",
								"author_2":         "Jane the Doe",
								"authorinitials_2": "Jt",
								"email_2":          "jane@example.com",
							},
							ElementReferences: types.ElementReferences{
								"_title": title,
							},
							Elements: []interface{}{
								types.Section{
									Level: 0,
									Attributes: types.Attributes{
										types.AttrID: "_title",
									},
									Title:    title,
									Elements: []interface{}{},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})
				})
			})

			Context("document revisions", func() {

				It("full document revision", func() {
					source := `= title
				John Doe
				v1.0, March 29, 2020: Updated revision`
					title := []interface{}{
						types.StringElement{
							Content: "title",
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrAuthors: []types.DocumentAuthor{
								{
									FullName: "John Doe",
								},
							},
							"firstname":      "John",
							"lastname":       "Doe",
							"author":         "John Doe",
							"authorinitials": "JD",
							types.AttrRevision: types.DocumentRevision{
								Revnumber: "1.0",
								Revdate:   "March 29, 2020",
								Revremark: "Updated revision",
							},
							"revnumber": "1.0",
							"revdate":   "March 29, 2020",
							"revremark": "Updated revision",
						},
						ElementReferences: types.ElementReferences{
							"_title": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 0,
								Attributes: types.Attributes{
									types.AttrID: "_title",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("full document revision with a comment before author", func() {
					source := `= title
				// a comment
				John Doe
				v1.0, March 29, 2020: Updated revision`
					title := []interface{}{
						types.StringElement{
							Content: "title",
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrAuthors: []types.DocumentAuthor{
								{
									FullName: "John Doe",
								},
							},
							"firstname":      "John",
							"lastname":       "Doe",
							"author":         "John Doe",
							"authorinitials": "JD",
							types.AttrRevision: types.DocumentRevision{
								Revnumber: "1.0",
								Revdate:   "March 29, 2020",
								Revremark: "Updated revision",
							},
							"revnumber": "1.0",
							"revdate":   "March 29, 2020",
							"revremark": "Updated revision",
						},
						ElementReferences: types.ElementReferences{
							"_title": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 0,
								Attributes: types.Attributes{
									types.AttrID: "_title",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("full document revision with a comment before revision", func() {
					source := `= title
				John Doe
				// a comment
				v1.0, March 29, 2020: Updated revision`
					title := []interface{}{
						types.StringElement{
							Content: "title",
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrAuthors: []types.DocumentAuthor{
								{
									FullName: "John Doe",
								},
							},
							"firstname":      "John",
							"lastname":       "Doe",
							"author":         "John Doe",
							"authorinitials": "JD",
							types.AttrRevision: types.DocumentRevision{
								Revnumber: "1.0",
								Revdate:   "March 29, 2020",
								Revremark: "Updated revision",
							},
							"revnumber": "1.0",
							"revdate":   "March 29, 2020",
							"revremark": "Updated revision",
						},
						ElementReferences: types.ElementReferences{
							"_title": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 0,
								Attributes: types.Attributes{
									types.AttrID: "_title",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("revision with revnumber and revdate only", func() {
					source := `= title
				John Doe
				v1.0, March 29, 2020`
					title := []interface{}{
						types.StringElement{
							Content: "title",
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrAuthors: []types.DocumentAuthor{
								{
									FullName: "John Doe",
								},
							},
							"firstname":      "John",
							"lastname":       "Doe",
							"author":         "John Doe",
							"authorinitials": "JD",
							types.AttrRevision: types.DocumentRevision{
								Revnumber: "1.0",
								Revdate:   "March 29, 2020",
							},
							"revnumber": "1.0",
							"revdate":   "March 29, 2020",
						},
						ElementReferences: types.ElementReferences{
							"_title": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 0,
								Attributes: types.Attributes{
									types.AttrID: "_title",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("revision with revnumber and revdate - with colon separator", func() {
					source := `= title
				John Doe
				1.0, March 29, 2020:`
					title := []interface{}{
						types.StringElement{
							Content: "title",
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrAuthors: []types.DocumentAuthor{
								{
									FullName: "John Doe",
								},
							},
							"firstname":      "John",
							"lastname":       "Doe",
							"author":         "John Doe",
							"authorinitials": "JD",
							types.AttrRevision: types.DocumentRevision{
								Revnumber: "1.0",
								Revdate:   "March 29, 2020",
							},
							"revnumber": "1.0",
							"revdate":   "March 29, 2020",
						},
						ElementReferences: types.ElementReferences{
							"_title": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 0,
								Attributes: types.Attributes{
									types.AttrID: "_title",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("revision with revnumber only - comma suffix", func() {
					source := `= title
				John Doe
				1.0,`
					title := []interface{}{
						types.StringElement{
							Content: "title",
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrAuthors: []types.DocumentAuthor{
								{
									FullName: "John Doe",
								},
							},
							"firstname":      "John",
							"lastname":       "Doe",
							"author":         "John Doe",
							"authorinitials": "JD",
							types.AttrRevision: types.DocumentRevision{
								Revnumber: "1.0",
							},
							"revnumber": "1.0",
						},
						ElementReferences: types.ElementReferences{
							"_title": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 0,
								Attributes: types.Attributes{
									types.AttrID: "_title",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("revision with revdate as number - spaces and no prefix no suffix", func() {
					source := `= title
				John Doe
				1.0`
					title := []interface{}{
						types.StringElement{
							Content: "title",
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrAuthors: []types.DocumentAuthor{
								{
									FullName: "John Doe",
								},
							},
							"firstname":      "John",
							"lastname":       "Doe",
							"author":         "John Doe",
							"authorinitials": "JD",
							types.AttrRevision: types.DocumentRevision{
								Revdate: "1.0",
							},
							"revdate": "1.0",
						},
						ElementReferences: types.ElementReferences{
							"_title": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 0,
								Attributes: types.Attributes{
									types.AttrID: "_title",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("revision with revdate as alphanum - spaces and no prefix no suffix", func() {
					source := `= title
				John Doe
				1.0a`
					title := []interface{}{
						types.StringElement{
							Content: "title",
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrAuthors: []types.DocumentAuthor{
								{
									FullName: "John Doe",
								},
							},
							"firstname":      "John",
							"lastname":       "Doe",
							"author":         "John Doe",
							"authorinitials": "JD",
							types.AttrRevision: types.DocumentRevision{
								Revdate: "1.0a",
							},
							"revdate": "1.0a",
						},
						ElementReferences: types.ElementReferences{
							"_title": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 0,
								Attributes: types.Attributes{
									types.AttrID: "_title",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("revision with revnumber only", func() {
					source := `= title
				John Doe
				v1.0:`
					title := []interface{}{
						types.StringElement{
							Content: "title",
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrAuthors: []types.DocumentAuthor{
								{
									FullName: "John Doe",
								},
							},
							"firstname":      "John",
							"lastname":       "Doe",
							"author":         "John Doe",
							"authorinitials": "JD",
							types.AttrRevision: types.DocumentRevision{
								Revnumber: "1.0",
							},
							"revnumber": "1.0",
						},
						ElementReferences: types.ElementReferences{
							"_title": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 0,
								Attributes: types.Attributes{
									types.AttrID: "_title",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("revision with spaces and capital revnumber ", func() {
					source := `= title
				John Doe
				V1.0:`
					title := []interface{}{
						types.StringElement{
							Content: "title",
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrAuthors: []types.DocumentAuthor{
								{
									FullName: "John Doe",
								},
							},
							"firstname":      "John",
							"lastname":       "Doe",
							"author":         "John Doe",
							"authorinitials": "JD",
							types.AttrRevision: types.DocumentRevision{
								Revnumber: "1.0",
							},
							"revnumber": "1.0",
						},
						ElementReferences: types.ElementReferences{
							"_title": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 0,
								Attributes: types.Attributes{
									types.AttrID: "_title",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("revision only - with comma separator", func() {
					source := `= title
				John Doe
				v1.0,`
					title := []interface{}{
						types.StringElement{
							Content: "title",
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrAuthors: []types.DocumentAuthor{
								{
									FullName: "John Doe",
								},
							},
							"firstname":      "John",
							"lastname":       "Doe",
							"author":         "John Doe",
							"authorinitials": "JD",
							types.AttrRevision: types.DocumentRevision{
								Revnumber: "1.0",
							},
							"revnumber": "1.0",
						},
						ElementReferences: types.ElementReferences{
							"_title": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 0,
								Attributes: types.Attributes{
									types.AttrID: "_title",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("revision with revnumber plus comma and colon separators", func() {
					source := `= title
				John Doe
				v1.0,:`
					title := []interface{}{
						types.StringElement{
							Content: "title",
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrAuthors: []types.DocumentAuthor{
								{
									FullName: "John Doe",
								},
							},
							"firstname":      "John",
							"lastname":       "Doe",
							"author":         "John Doe",
							"authorinitials": "JD",
							types.AttrRevision: types.DocumentRevision{
								Revnumber: "1.0",
							},
							"revnumber": "1.0",
						},
						ElementReferences: types.ElementReferences{
							"_title": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 0,
								Attributes: types.Attributes{
									types.AttrID: "_title",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("revision with revnumber and empty revremark", func() {
					source := `= title
John Doe
v1.0:`
					title := []interface{}{
						types.StringElement{
							Content: "title",
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrAuthors: []types.DocumentAuthor{
								{
									FullName: "John Doe",
								},
							},
							"firstname":      "John",
							"lastname":       "Doe",
							"author":         "John Doe",
							"authorinitials": "JD",
							types.AttrRevision: types.DocumentRevision{
								Revnumber: "1.0",
							},
							"revnumber": "1.0",
							// "revremark": "", // found but is empty
						},
						ElementReferences: types.ElementReferences{
							"_title": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 0,
								Attributes: types.Attributes{
									types.AttrID: "_title",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

			})

			Context("document header attributes", func() {

				It("valid attribute names", func() {
					source := `:a:
:author: Xavier
:_author: Xavier`
					expected := types.Document{
						Attributes: types.Attributes{
							"a":       nil,
							"author":  "Xavier",
							"_author": "Xavier",
						},
						Elements: []interface{}{},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("attributes and paragraph without blank line in-between", func() {
					source := `:toc:
:date:  2017-01-01
:author: Xavier
:hardbreaks:
a paragraph`
					expected := types.Document{
						Attributes: types.Attributes{
							"toc":                        nil,
							"date":                       "2017-01-01",
							"author":                     "Xavier",
							types.DocumentAttrHardBreaks: nil,
						},
						Elements: []interface{}{
							types.TableOfContentsPlaceHolder{},
							types.Paragraph{
								Lines: [][]interface{}{
									{types.StringElement{Content: "a paragraph"}},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("contiguous attributes and paragraph with blank line in-between", func() {
					source := `:toc:
:date: 2017-01-01
:author: Xavier

a paragraph`
					expected := types.Document{
						Attributes: types.Attributes{
							"toc":    nil,
							"date":   "2017-01-01",
							"author": "Xavier",
						},
						Elements: []interface{}{
							types.TableOfContentsPlaceHolder{},
							types.Paragraph{
								Lines: [][]interface{}{
									{types.StringElement{Content: "a paragraph"}},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("splitted attributes and paragraph with blank line in-between", func() {
					source := `:toc:
:date: 2017-01-01

:author: Xavier

:hardbreaks:

a paragraph`
					expected := types.Document{
						Attributes: types.Attributes{
							"toc":        nil,
							"date":       "2017-01-01",
							"author":     "Xavier",
							"hardbreaks": nil,
						},
						Elements: []interface{}{
							types.TableOfContentsPlaceHolder{},
							types.Paragraph{
								Lines: [][]interface{}{
									{types.StringElement{Content: "a paragraph"}},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("no header and attributes in body", func() {
					source := `a paragraph
	
:toc:
:date: 2017-01-01
:author: Xavier`
					expected := types.Document{
						Attributes: types.Attributes{
							"toc":    nil,
							"date":   "2017-01-01",
							"author": "Xavier",
						},
						Elements: []interface{}{
							types.TableOfContentsPlaceHolder{},
							types.Paragraph{
								Lines: [][]interface{}{
									{types.StringElement{Content: "a paragraph"}},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})
			})

		})

		Context("invalid cases", func() {

			It("paragraph without blank line before attribute declarations", func() {
				source := `a paragraph
:toc:
:date: 2017-01-01
:author: Xavier`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a paragraph"},
								},
								{
									types.StringElement{Content: ":toc:"},
								},
								{
									types.StringElement{Content: ":date: 2017-01-01"},
								},
								{
									types.StringElement{Content: ":author: Xavier"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("invalid attribute names", func() {
				source := `:@date: 2017-01-01
:{author}: Xavier`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: ":@date: 2017-01-01"},
								},
								{
									types.StringElement{Content: ":{author}: Xavier"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})

		Context("with overrides", func() {

			It("custom icon attribute", func() {
				// given
				attrs := map[string]string{
					"icons":              "font",
					"source-highlighter": "pygments",
				}
				source := `{icons}`
				expected := types.Document{
					Attributes: types.Attributes{
						"icons":              "font",
						"source-highlighter": "pygments",
					},
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "font"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source, configuration.WithAttributes(attrs))).To(Equal(expected))
			})
		})
	})
})
