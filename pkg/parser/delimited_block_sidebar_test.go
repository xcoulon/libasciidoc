package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("sidebar blocks", func() {

	Context("in final documents", func() {

		Context("delimited blocks", func() {

			It("with rich content", func() {
				source := `****
some *verse* content
****`
				expected := types.Document{
					Elements: []interface{}{
						types.SidebarBlock{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{
												Content: "some ",
											},
											&types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{
														Content: "verse",
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
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with title, paragraph and sourcecode block", func() {
				source := `.a title
****
some *verse* content

----
foo
bar
----
****`
				expected := types.Document{
					Elements: []interface{}{
						types.SidebarBlock{
							Attributes: types.Attributes{
								types.AttrTitle: "a title",
							},
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{
												Content: "some ",
											},
											&types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{
														Content: "verse",
													},
												},
											},
											types.StringElement{
												Content: " content",
											},
										},
									},
								},
								&types.BlankLine{}, // blankline is required between paragraph and the next block
								types.ListingBlock{
									Elements: []interface{}{
										types.StringElement{
											Content: "foo",
										},
										types.StringElement{
											Content: "bar",
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
				source := `
****
.foo
****`
				expected := types.Document{
					Elements: []interface{}{
						types.SidebarBlock{
							Elements: []interface{}{},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})
	})
	Context("in final documents", func() {

		Context("sidebar blocks", func() {

			It("with paragraph", func() {
				source := `****
some *verse* content
****`
				expected := types.Document{
					Elements: []interface{}{
						types.SidebarBlock{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{
												Content: "some ",
											},
											&types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{
														Content: "verse",
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
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with title, paragraph and sourcecode block", func() {
				source := `.a title
****
some *verse* content

----
foo
bar
----
****`
				expected := types.Document{
					Elements: []interface{}{
						types.SidebarBlock{
							Attributes: types.Attributes{
								types.AttrTitle: "a title",
							},
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{
												Content: "some ",
											},
											&types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{
														Content: "verse",
													},
												},
											},
											types.StringElement{
												Content: " content",
											},
										},
									},
								},
								&types.BlankLine{}, // blankline is required between paragraph and the next block
								types.ListingBlock{
									Elements: []interface{}{
										types.StringElement{
											Content: "foo",
										},
										types.StringElement{
											Content: "bar",
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
				source := `
****
.foo
****`
				expected := types.Document{
					Elements: []interface{}{
						types.SidebarBlock{
							Elements: []interface{}{},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

		})
	})
})
