package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("documents", func() {

	Context("raw documents", func() {

		It("should parse empty document", func() {
			source := ``
			expected := types.DocumentFragments{}
			Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
		})

		It("should parse header without empty first line", func() {
			source := `= My title
Garrett D'Amore
1.0, July 4, 2020
`
			expected := types.DocumentFragments{
				types.Section{
					Level: 0,
					Attributes: types.Attributes{
						"id": "_my_title",
					},
					Title: []interface{}{
						types.StringElement{
							Content: "My title",
						},
					},
				},
				types.DocumentAuthor{
					FullName: "Garrett D'Amore",
					Email:    "",
				},
				types.DocumentRevision{
					Revnumber: "1.0",
					Revdate:   "July 4, 2020",
				},
			}
			Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))

		})

		It("should parse header with empty first line", func() {
			source := `
= My title
Garrett D'Amore
1.0, July 4, 2020`
			expected := types.DocumentFragments{
				types.Section{
					Level: 0,
					Attributes: types.Attributes{
						"id": "_my_title",
					},
					Title: []interface{}{
						types.StringElement{
							Content: "My title",
						},
					},
				},
				types.DocumentAuthor{
					FullName: "Garrett D'Amore",
					Email:    "",
				},
				types.DocumentRevision{
					Revnumber: "1.0",
					Revdate:   "July 4, 2020",
					Revremark: "",
				},
			}
			Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))

		})
	})

	Context("in final documents", func() {

		It("should parse empty document", func() {
			source := ``
			expected := types.Document{
				Elements: []interface{}{},
			}
			Expect(ParseDocument(source)).To(Equal(expected))
		})
	})
})
