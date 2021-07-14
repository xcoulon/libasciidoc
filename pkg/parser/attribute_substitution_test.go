package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("attribute substitutions", func() {

	Context("in final documents", func() {

		It("paragraph with attribute substitution", func() {
			source := `:author: Xavier

a paragraph written by {author}.`
			expected := &types.Document{
				Attributes: types.Attributes{
					"author": "Xavier",
				},
				Elements: []interface{}{
					types.Paragraph{
						Lines: [][]interface{}{
							{types.StringElement{Content: "a paragraph written by Xavier."}},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("paragraph with attribute resets", func() {
			source := `:author: Xavier
				
:!author1:
:author2!:
a paragraph written by {author}.`
			expected := &types.Document{
				Attributes: types.Attributes{
					"author":  "Xavier",
					"author1": nil,
					"author2": nil,
				},
				Elements: []interface{}{
					types.Paragraph{
						Lines: [][]interface{}{
							{types.StringElement{Content: "a paragraph written by Xavier."}},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("header with 2 authors, revision and attributes", func() {
			source := `= Document Title
John Foo Doe <johndoe@example.com>; Jane the_Doe <jane@example.com>
v1.0, March 29, 2020: Updated revision
:toc:
:keywords: documentation, team, obstacles, journey, victory

This journey continues`
			title := []interface{}{
				&types.StringElement{Content: "Document Title"},
			}
			expected := &types.Document{
				Attributes: types.Attributes{
					"toc":      nil,
					"keywords": "documentation, team, obstacles, journey, victory",
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
							types.TableOfContentsPlaceHolder{},
							types.Paragraph{
								Lines: [][]interface{}{
									{
										&types.StringElement{Content: "This journey continues"},
									},
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

		It("paragraph with attribute substitution from front-matter", func() {
			source := `---
author: Xavier
---

a paragraph written by {author}.`
			expected := &types.Document{
				Attributes: types.Attributes{
					"author": "Xavier",
				},
				Elements: []interface{}{
					types.Paragraph{
						Lines: [][]interface{}{
							{types.StringElement{Content: "a paragraph written by Xavier."}},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
	})
})
