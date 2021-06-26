package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("markdown-style quote blocks", func() {

	Context("in final documents", func() {

		Context("as delimited blocks", func() {

			It("with single marker without author", func() {
				source := `> some text
on *multiple lines*`

				expected := types.Document{
					Elements: []interface{}{
						types.MarkdownQuoteBlock{
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "some text",
									},
								},
								{
									types.StringElement{
										Content: "on ",
									},
									&types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{
												Content: "multiple lines",
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

			It("with marker on each line without author", func() {
				source := `> some text
> on *multiple lines*`
				expected := types.Document{
					Elements: []interface{}{
						types.MarkdownQuoteBlock{
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "some text",
									},
								},
								{
									types.StringElement{
										Content: "on ",
									},
									&types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{
												Content: "multiple lines",
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

			It("with marker on each line with author only", func() {
				source := `> some text
> on *multiple lines*
> -- John Doe`
				expected := types.Document{
					Elements: []interface{}{
						types.MarkdownQuoteBlock{
							Attributes: types.Attributes{
								types.AttrQuoteAuthor: "John Doe",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "some text",
									},
								},
								{
									types.StringElement{
										Content: "on ",
									},
									&types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{
												Content: "multiple lines",
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

			It("with marker on each line with author and title", func() {
				source := `.title
> some text
> on *multiple lines*
> -- John Doe`
				expected := types.Document{
					Elements: []interface{}{
						types.MarkdownQuoteBlock{
							Attributes: types.Attributes{
								types.AttrTitle:       "title",
								types.AttrQuoteAuthor: "John Doe",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "some text",
									},
								},
								{
									types.StringElement{
										Content: "on ",
									},
									&types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{
												Content: "multiple lines",
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

			It("with with author only", func() {
				source := `> -- John Doe`
				expected := types.Document{
					Elements: []interface{}{
						types.MarkdownQuoteBlock{
							Attributes: types.Attributes{
								types.AttrQuoteAuthor: "John Doe",
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})
	})
})
