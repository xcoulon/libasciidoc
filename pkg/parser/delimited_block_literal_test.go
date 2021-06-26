package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("literal blocks", func() {

	Context("in final documents", func() {

		Context("with spaces indentation", func() {

			It("from 1-line paragraph with single space", func() {
				source := ` some literal content`
				expected := types.Document{
					Elements: []interface{}{
						&types.DelimitedBlock{
							Kind: types.Literal,
							Attributes: types.Attributes{
								types.AttrStyle:            types.Literal,
								types.AttrLiteralBlockType: types.LiteralBlockWithSpacesOnFirstLine,
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: " some literal content",
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("from paragraph with single space on first line", func() {
				source := ` some literal content
on 3
lines.`
				expected := types.Document{
					Elements: []interface{}{
						&types.DelimitedBlock{
							Kind: types.Literal,
							Attributes: types.Attributes{
								types.AttrStyle:            types.Literal,
								types.AttrLiteralBlockType: types.LiteralBlockWithSpacesOnFirstLine,
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: " some literal content",
								},
								&types.StringElement{
									Content: "on 3",
								},
								&types.StringElement{
									Content: "lines.",
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("mixing literal block with attributes followed by a paragraph ", func() {
				source := `.title
[#ID]
  some literal content

a normal paragraph.`
				expected := types.Document{
					Elements: []interface{}{
						&types.DelimitedBlock{
							Kind: types.Literal,
							Attributes: types.Attributes{
								types.AttrStyle:            types.Literal,
								types.AttrLiteralBlockType: types.LiteralBlockWithSpacesOnFirstLine,
								types.AttrID:               "ID",
								// types.AttrCustomID:         true,
								types.AttrTitle: "title",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "  some literal content",
								},
							},
						},
						&types.BlankLine{},
						&types.Paragraph{
							Elements: []interface{}{
								&types.StringElement{
									Content: "a normal paragraph.",
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})

		Context("with block delimiter", func() {

			It("with empty blank line", func() {
				source := `....

some content
....`
				expected := types.Document{
					Elements: []interface{}{
						&types.DelimitedBlock{
							Kind: types.Literal,
							Elements: []interface{}{
								&types.StringElement{
									Content: "\nsome content",
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with delimited and attributes followed by 1-line paragraph", func() {
				source := `[#ID]
.title
....
some literal content
....
a normal paragraph.`
				expected := types.Document{
					Elements: []interface{}{
						&types.DelimitedBlock{
							Kind: types.Literal,
							Attributes: types.Attributes{
								types.AttrID: "ID",
								// types.AttrCustomID:         true,
								types.AttrTitle: "title",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "some literal content",
								},
							},
						},
						&types.Paragraph{
							Elements: []interface{}{
								&types.StringElement{
									Content: "a normal paragraph.",
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})

		Context("with attribute", func() {

			It("from 1-line paragraph with attribute", func() {
				source := `[literal]   
some literal content

a normal paragraph.`
				expected := types.Document{
					Elements: []interface{}{
						&types.DelimitedBlock{
							Kind: types.Literal,
							Attributes: types.Attributes{
								types.AttrStyle:            types.Literal,
								types.AttrLiteralBlockType: types.LiteralBlockWithAttribute,
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "some literal content",
								},
							},
						},
						&types.BlankLine{},
						&types.Paragraph{
							Elements: []interface{}{
								&types.StringElement{
									Content: "a normal paragraph.",
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("from 2-lines paragraph with attribute", func() {
				source := `[#ID]
[literal]   
.title
some literal content
on two lines.

a normal paragraph.`
				expected := types.Document{
					Elements: []interface{}{
						&types.DelimitedBlock{
							Kind: types.Literal,
							Attributes: types.Attributes{
								types.AttrStyle: types.Literal,
								types.AttrID:    "ID",
								// types.AttrCustomID:         true,
								types.AttrTitle:            "title",
								types.AttrLiteralBlockType: types.LiteralBlockWithAttribute,
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "some literal content",
								},
								&types.StringElement{
									Content: "on two lines.",
								},
							},
						},
						&types.BlankLine{},
						&types.Paragraph{
							Elements: []interface{}{
								&types.StringElement{
									Content: "a normal paragraph.",
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})

		Context("with custom substitutions", func() {

			// testing custom substitutions on a literal block only

			It("should apply the default substitution on block with delimiter", func() {
				source := `:github-url: https://github.com
				
....
a link to https://example.com[] <1>
and <more text> on the +
*next* lines with a link to {github-url}[]

* not a list item
....

<1> a callout
`
				expected := types.Document{
					Attributes: types.Attributes{
						"github-url": "https://github.com",
					},
					Elements: []interface{}{
						&types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						&types.BlankLine{},
						&types.DelimitedBlock{
							Kind: types.Literal,
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to https://example.com[] ",
								},
								types.Callout{
									Ref: 1,
								},
								&types.StringElement{
									Content: "\nand ",
								},
								types.SpecialCharacter{
									Name: "<",
								},
								&types.StringElement{
									Content: "more text",
								},
								types.SpecialCharacter{
									Name: ">",
								},
								&types.StringElement{
									Content: " on the +\n*next* lines with a link to {github-url}[]\n\n* not a list item",
								},
							},
						},
						&types.BlankLine{},
						&types.GenericList{
							Kind: types.CalloutListKind,
							Elements: []types.ListElement{
								&types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												&types.StringElement{
													Content: "a callout",
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
			It("should apply the 'normal' substitution on block with delimiter", func() {
				source := `:github-url: https://github.com
				
[subs="normal"]
....
a link to https://example.com[] <1>
and <more text> on the +
*next* lines with a link to {github-url}[]

* not a list item
....

<1> a callout
`
				expected := types.Document{
					Attributes: types.Attributes{
						"github-url": "https://github.com",
					},
					Elements: []interface{}{
						&types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						&types.BlankLine{},
						&types.DelimitedBlock{
							Kind: types.Literal,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "normal",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to ",
								},
								&types.InlineLink{
									Location: &types.Location{
										Scheme: "https://",
										Path: []interface{}{
											&types.StringElement{
												Content: "example.com",
											},
										},
									},
								},
								&types.StringElement{
									Content: " ",
								},
								types.SpecialCharacter{ // callout is not detected with the `normal` susbtitution
									Name: "<",
								},
								&types.StringElement{
									Content: "1",
								},
								types.SpecialCharacter{
									Name: ">",
								},
								&types.StringElement{
									Content: "\nand ",
								},
								types.SpecialCharacter{
									Name: "<",
								},
								&types.StringElement{
									Content: "more text",
								},
								types.SpecialCharacter{
									Name: ">",
								},
								&types.StringElement{
									Content: " on the",
								},
								types.LineBreak{},
								&types.StringElement{
									Content: "\n",
								},
								&types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										&types.StringElement{
											Content: "next",
										},
									},
								},
								&types.StringElement{
									Content: " lines with a link to ",
								},
								&types.InlineLink{
									Location: &types.Location{
										Scheme: "https://",
										Path: []interface{}{
											&types.StringElement{
												Content: "github.com",
											},
										},
									},
								},
								&types.StringElement{
									Content: "\n\n* not a list item",
								},
							},
						},
						&types.BlankLine{},
						&types.GenericList{
							Kind: types.CalloutListKind,
							Elements: []types.ListElement{
								&types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												&types.StringElement{
													Content: "a callout",
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

			It("should apply the 'quotes,macros' substitution on block with delimiter", func() {
				source := `:github-url: https://github.com
				
[subs="quotes,macros"]
....
a link to https://example.com[] <1>
and <more text> on the +
*next* lines with a link to {github-url}[]

* not a list item
....

<1> a callout
`
				expected := types.Document{
					Attributes: types.Attributes{
						"github-url": "https://github.com",
					},
					Elements: []interface{}{
						&types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						&types.BlankLine{},
						&types.DelimitedBlock{
							Kind: types.Literal,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "quotes,macros",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to ",
								},
								&types.InlineLink{
									Location: &types.Location{
										Scheme: "https://",
										Path: []interface{}{
											&types.StringElement{
												Content: "example.com",
											},
										},
									},
								},
								&types.StringElement{
									Content: " <1>\nand <more text> on the +\n",
								},
								&types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										&types.StringElement{
											Content: "next",
										},
									},
								},
								&types.StringElement{
									Content: " lines with a link to {github-url}[]\n\n* not a list item",
								},
							},
						},
						&types.BlankLine{},
						&types.GenericList{
							Kind: types.CalloutListKind,
							Elements: []types.ListElement{
								&types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												&types.StringElement{
													Content: "a callout",
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

			It("should apply the 'quotes,macros' substitution on block with spaces", func() {
				source := `:github-url: https://github.com
				
[subs="quotes,macros"]
  a link to https://example.com[] <1> 
  and <more text> on the +
  *next* lines with a link to {github-url}[]

<1> a callout
`
				expected := types.Document{
					Attributes: types.Attributes{
						"github-url": "https://github.com",
					},
					Elements: []interface{}{
						&types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						&types.BlankLine{},
						&types.DelimitedBlock{
							Kind: types.Literal,
							Attributes: types.Attributes{
								types.AttrStyle:            types.Literal,
								types.AttrLiteralBlockType: types.LiteralBlockWithSpacesOnFirstLine,
								types.AttrSubstitutions:    "quotes,macros",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "  a link to ",
								},
								&types.InlineLink{
									Location: &types.Location{
										Scheme: "https://",
										Path: []interface{}{
											&types.StringElement{
												Content: "example.com",
											},
										},
									},
								},
								&types.StringElement{
									Content: " <1> ",
								},
								&types.StringElement{
									Content: "  and <more text> on the +",
								},
								&types.StringElement{
									Content: "  ",
								},
								&types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										&types.StringElement{
											Content: "next",
										},
									},
								},
								&types.StringElement{
									Content: " lines with a link to {github-url}[]",
								},
							},
						},
						&types.BlankLine{},
						&types.GenericList{
							Kind: types.CalloutListKind,
							Elements: []types.ListElement{
								&types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												&types.StringElement{
													Content: "a callout",
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

			It("should apply the 'quotes,macros' substitution on block with attribute", func() {
				source := `:github-url: https://github.com
				
[subs="quotes,macros"]
[literal]
a link to https://example.com[] <1>
and <more text> on the +
*next* lines with a link to {github-url}[]

<1> a callout
`
				expected := types.Document{
					Attributes: types.Attributes{
						"github-url": "https://github.com",
					},
					Elements: []interface{}{
						&types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						&types.BlankLine{},
						&types.DelimitedBlock{
							Kind: types.Literal,
							Attributes: types.Attributes{
								types.AttrStyle:            types.Literal,
								types.AttrLiteralBlockType: types.LiteralBlockWithAttribute,
								types.AttrSubstitutions:    "quotes,macros",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to ",
								},
								&types.InlineLink{
									Location: &types.Location{
										Scheme: "https://",
										Path: []interface{}{
											&types.StringElement{
												Content: "example.com",
											},
										},
									},
								},
								&types.StringElement{
									Content: " <1>",
								},
								&types.StringElement{
									Content: "and <more text> on the +",
								},
								&types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										&types.StringElement{
											Content: "next",
										},
									},
								},
								&types.StringElement{
									Content: " lines with a link to {github-url}[]",
								},
							},
						},
						&types.BlankLine{},
						&types.GenericList{
							Kind: types.CalloutListKind,
							Elements: []types.ListElement{
								&types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												&types.StringElement{
													Content: "a callout",
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
})
