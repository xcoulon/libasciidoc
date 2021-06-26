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
			expected := []types.DocumentFragment{}
			Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
		})

		It("should parse header without empty first line", func() {
			source := `= My title
Garrett D'Amore
1.0, July 4, 2020
`
			expected := []types.DocumentFragment{
				{
					LineOffset: 1,
					Elements: []interface{}{
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
					},
				},
			}
			Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))

		})

		It("should parse header with empty first line", func() {
			source := `
= My title
Garrett D'Amore
1.0, July 4, 2020`
			expected := []types.DocumentFragment{
				{
					LineOffset: 1,
					Elements: []interface{}{
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
					},
				},
			}
			Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))

		})
	})

	Context("in final documents", func() {

		It("should parse empty document", func() {
			source := ``
			expected := types.Document{
				Elements: []interface{}{},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("should parse basic document", func() {
			source := `== Lorem Ipsum
			
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, 
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, 
sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. 
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. 
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, 
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, 
sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. 
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit *amet*.`

			expected := types.Document{
				Elements: []interface{}{
					&types.Section{
						Level: 1,
						Title: []interface{}{
							types.StringElement{
								Content: "Lorem Ipsum",
							},
						},
					},
					&types.BlankLine{},
					&types.Paragraph{
						Elements: []interface{}{
							types.StringElement{
								Content: `Lorem ipsum dolor sit amet, consetetur sadipscing elitr, 
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, 
sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. 
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. 
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, 
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, 
sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. 
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit `,
							},
							&types.QuotedText{
								Kind: types.SingleQuoteBold,
								Elements: []interface{}{
									types.StringElement{
										Content: "amet",
									},
								},
							},
							types.StringElement{
								Content: ".",
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
	})
})
