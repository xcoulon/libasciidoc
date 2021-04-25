package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("paragraphs", func() {

	Context("in final documents", func() {

		Context("thematic breaks", func() {

			It("thematic break form1 by itself", func() {
				source := "***"
				expected := types.Document{
					Elements: []interface{}{
						types.ThematicBreak{},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("thematic break form2 by itself", func() {
				source := "* * *"
				expected := types.Document{
					Elements: []interface{}{
						types.ThematicBreak{},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("thematic break form3 by itself", func() {
				source := "---"
				expected := types.Document{
					Elements: []interface{}{
						types.ThematicBreak{},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("thematic break form4 by itself", func() {
				source := "- - -"
				expected := types.Document{
					Elements: []interface{}{
						types.ThematicBreak{},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("thematic break form5 by itself", func() {
				source := "___"
				expected := types.Document{
					Elements: []interface{}{
						types.ThematicBreak{},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("thematic break form4 by itself", func() {
				source := "_ _ _"
				expected := types.Document{
					Elements: []interface{}{
						types.ThematicBreak{},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("thematic break with leading text", func() {
				source := "text ***"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "text ***"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			// NB: three asterisks gets confused with bullets if with trailing text
			It("thematic break with trailing text", func() {
				source := "* * * text"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "* * * text"},
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
