package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("front-matters", func() {

	Context("in raw documents", func() {

		Context("alone", func() {

			It("with simple attributes and no blanklines", func() {
				source := `---
title: a title
author: Xavier
---
`
				expected := types.DocumentFragments{
					types.BlockDelimiter{
						Kind: types.FrontMatter,
					},
					types.StringElement{
						Content: "title: a title",
					},
					types.StringElement{
						Content: "author: Xavier",
					},
					types.BlockDelimiter{
						Kind: types.FrontMatter,
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
			})

			It("with simple attributes and blanklines", func() {
				source := `---

title: a title

author: Xavier

---
`
				expected := types.DocumentFragments{
					types.BlockDelimiter{
						Kind: types.FrontMatter,
					},
					types.BlankLine{},
					types.StringElement{
						Content: "title: a title",
					},
					types.BlankLine{},
					types.StringElement{
						Content: "author: Xavier",
					},
					types.BlankLine{},
					types.BlockDelimiter{
						Kind: types.FrontMatter,
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
			})

			It("empty front-matter", func() {
				source := `---
---`
				expected := types.DocumentFragments{
					types.BlockDelimiter{
						Kind: types.FrontMatter,
					},
					types.BlockDelimiter{
						Kind: types.FrontMatter,
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
			})

			It("no front-matter", func() {
				source := `some content`
				expected := types.DocumentFragments{
					types.StringElement{
						Content: "some content",
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
			})
		})

		Context("with content afterwards", func() {

			It("with document header", func() {
				source := `---
title: a title
author: Xavier
---
= A Title
`
				expected := types.DocumentFragments{
					types.BlockDelimiter{
						Kind: types.FrontMatter,
					},
					types.StringElement{
						Content: "title: a title",
					},
					types.StringElement{
						Content: "author: Xavier",
					},
					types.BlockDelimiter{
						Kind: types.FrontMatter,
					},
					types.Section{
						Level: 0,
						Title: []interface{}{
							types.StringElement{
								Content: "A Title",
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
			})
		})
	})

	Context("in final documents", func() {

		Context("yaml front-matter", func() {

			It("with simple attributes", func() {
				source := `---
title: a title
author: Xavier
---

first paragraph`
				expected := types.Document{
					Attributes: types.Attributes{
						"title":  "a title", // TODO: convert `title` attribute from front-matter into `doctitle` here ?
						"author": "Xavier",
					},
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "first paragraph"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("empty front-matter", func() {
				source := `---
---

first paragraph`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "first paragraph"},
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
