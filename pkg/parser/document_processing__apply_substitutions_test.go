package parser_test

import (
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("apply custom substitutions", func() {

	Context("example blocks", func() {

		// in normal blocks, the substiution should be defined and applied on the elements
		// within the blocks
		// TODO: include character replacement (eg: `(C)`)
		source := `:github-url: https://github.com
			
====
[subs="$SUBS"]
a link to https://example.com[] <1>
and <more text> on the +
*next* lines with a link to {github-url}[]

====
`
		Context("explicit substitutions", func() {

			It("should apply the default substitution", func() {
				s := strings.ReplaceAll(source, "[subs=\"$SUBS\"]\n", "")
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
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
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
										&types.SpecialCharacter{
											Name: "<",
										},
										&types.StringElement{
											Content: "1",
										},
										&types.SpecialCharacter{
											Name: ">",
										},
										&types.StringElement{
											Content: "\nand ",
										},
										&types.SpecialCharacter{
											Name: "<",
										},
										&types.StringElement{
											Content: "more text",
										},
										&types.SpecialCharacter{
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
									},
								},
								&types.BlankLine{},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'normal' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "normal")
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
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
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
										&types.SpecialCharacter{ // callout is not detected with the `normal` susbtitution
											Name: "<",
										},
										&types.StringElement{
											Content: "1",
										},
										&types.SpecialCharacter{
											Name: ">",
										},
										&types.StringElement{
											Content: "\nand ",
										},
										&types.SpecialCharacter{
											Name: "<",
										},
										&types.StringElement{
											Content: "more text",
										},
										&types.SpecialCharacter{
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
									},
								},
								&types.BlankLine{},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'quotes' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "quotes")
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
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
									Attributes: types.Attributes{
										types.AttrSubstitutions: "quotes",
									},
									Elements: []interface{}{
										&types.StringElement{
											Content: "a link to https://example.com[] <1>\nand <more text> on the +\n",
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
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'macros' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "macros")
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
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
									Attributes: types.Attributes{
										types.AttrSubstitutions: "macros",
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
											Content: " <1>\nand <more text> on the +\n*next* lines with a link to {github-url}[]",
										},
									},
								},
								&types.BlankLine{},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'attributes' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "attributes")
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
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
									Attributes: types.Attributes{
										types.AttrSubstitutions: "attributes",
									},
									Elements: []interface{}{
										&types.StringElement{
											Content: "a link to https://example.com[] <1>\nand <more text> on the +\n*next* lines with a link to https://github.com[]",
										},
									},
								},
								&types.BlankLine{},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'attributes,macros' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "attributes,macros")
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
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
									Attributes: types.Attributes{
										types.AttrSubstitutions: "attributes,macros",
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
											Content: " <1>\nand <more text> on the +\n*next* lines with a link to ",
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
									},
								},
								&types.BlankLine{},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'specialchars' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "specialchars")
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
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
									Attributes: types.Attributes{
										types.AttrSubstitutions: "specialchars",
									},
									Elements: []interface{}{
										&types.StringElement{
											Content: "a link to https://example.com[] ",
										},
										&types.SpecialCharacter{
											Name: "<",
										},
										&types.StringElement{
											Content: "1",
										},
										&types.SpecialCharacter{
											Name: ">",
										},
										&types.StringElement{
											Content: "\nand ",
										},
										&types.SpecialCharacter{
											Name: "<",
										},
										&types.StringElement{
											Content: "more text",
										},
										&types.SpecialCharacter{
											Name: ">",
										},
										&types.StringElement{
											Content: " on the +\n*next* lines with a link to {github-url}[]",
										},
									},
								},
								&types.BlankLine{},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'replacements' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "replacements")
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
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
									Attributes: types.Attributes{
										types.AttrSubstitutions: "replacements",
									},
									Elements: []interface{}{
										&types.StringElement{
											Content: "a link to https://example.com[] <1>\nand <more text> on the +\n*next* lines with a link to {github-url}[]",
										},
									},
								},
								&types.BlankLine{},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'post_replacements' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "post_replacements")
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
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
									Attributes: types.Attributes{
										types.AttrSubstitutions: "post_replacements",
									},
									Elements: []interface{}{
										&types.StringElement{
											Content: "a link to https://example.com[] <1>\nand <more text> on the",
										},
										types.LineBreak{},
										&types.StringElement{
											Content: "\n*next* lines with a link to {github-url}[]",
										},
									},
								},
								&types.BlankLine{},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'quotes,macros' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "quotes,macros")
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
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
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
											Content: " lines with a link to {github-url}[]",
										},
									},
								},
								&types.BlankLine{},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'macros,quotes' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "macros,quotes")
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
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
									Attributes: types.Attributes{
										types.AttrSubstitutions: "macros,quotes",
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
											Content: " lines with a link to {github-url}[]",
										},
									},
								},
								&types.BlankLine{},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'none' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "none")
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
							Kind: types.Example,
							Elements: []interface{}{
								&types.Paragraph{
									Attributes: types.Attributes{
										types.AttrSubstitutions: "none",
									},
									Elements: []interface{}{
										&types.StringElement{
											Content: "a link to https://example.com[] <1>\nand <more text> on the +\n*next* lines with a link to {github-url}[]",
										},
									},
								},
								&types.BlankLine{},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})
		})

	})

	Context("listing blocks", func() {

		Context("as delimited blocks", func() {
			// testing custom substitutions on listing blocks only, as
			// other verbatim blocks (fenced, literal, source, passthrough)
			// share the same implementation

			source := `:github-url: https://github.com

[subs="$SUBS"]
----
a link to https://example.com[] <1>
and <more text> on the +
*next* lines with a link to {github-url}[]

* not a list item
----
`
			It("should apply the default substitution", func() {
				s := strings.ReplaceAll(source, "[subs=\"$SUBS\"]\n", "") // remove the 'subs' attribute
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
							Kind: types.Listing,
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to https://example.com[] ",
								},
								&types.Callout{
									Ref: 1,
								},
								&types.StringElement{
									Content: "\nand ",
								},
								&types.SpecialCharacter{
									Name: "<",
								},
								&types.StringElement{
									Content: "more text",
								},
								&types.SpecialCharacter{
									Name: ">",
								},
								&types.StringElement{
									Content: " on the +\n*next* lines with a link to {github-url}[]\n\n* not a list item",
								},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'normal' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "normal")
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
							Kind: types.Listing,
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
								&types.SpecialCharacter{ // callout is not detected with the `normal` susbtitution
									Name: "<",
								},
								&types.StringElement{
									Content: "1",
								},
								&types.SpecialCharacter{
									Name: ">",
								},
								&types.StringElement{
									Content: "\nand ",
								},
								&types.SpecialCharacter{
									Name: "<",
								},
								&types.StringElement{
									Content: "more text",
								},
								&types.SpecialCharacter{
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
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'quotes' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "quotes")
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
							Kind: types.Listing,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "quotes",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to https://example.com[] <1>\nand <more text> on the +\n",
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
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'macros' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "macros")
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
							Kind: types.Listing,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "macros",
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
									Content: " <1>\nand <more text> on the +\n*next* lines with a link to {github-url}[]\n\n* not a list item",
								},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'attributes' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "attributes")
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
							Kind: types.Listing,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "attributes",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to https://example.com[] <1>\nand <more text> on the +\n*next* lines with a link to https://github.com[]\n\n* not a list item",
								},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'attributes,macros' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "attributes,macros")
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
							Kind: types.Listing,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "attributes,macros",
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
									Content: " <1>\nand <more text> on the +\n*next* lines with a link to ",
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
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'specialchars' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "specialchars")
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
							Kind: types.Listing,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "specialchars",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to https://example.com[] ",
								},
								&types.SpecialCharacter{
									Name: "<",
								},
								&types.StringElement{
									Content: "1",
								},
								&types.SpecialCharacter{
									Name: ">",
								},
								&types.StringElement{
									Content: "\nand ",
								},
								&types.SpecialCharacter{
									Name: "<",
								},
								&types.StringElement{
									Content: "more text",
								},
								&types.SpecialCharacter{
									Name: ">",
								},
								&types.StringElement{
									Content: " on the +\n*next* lines with a link to {github-url}[]\n\n* not a list item",
								},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'replacements' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "replacements")
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
							Kind: types.Listing,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "replacements",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to https://example.com[] <1>\nand <more text> on the +\n*next* lines with a link to {github-url}[]\n\n* not a list item",
								},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'post_replacements' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "post_replacements")
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
							Kind: types.Listing,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "post_replacements",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to https://example.com[] <1>\nand <more text> on the",
								},
								types.LineBreak{},
								&types.StringElement{
									Content: "\n*next* lines with a link to {github-url}[]\n\n* not a list item",
								},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'quotes,macros' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "quotes,macros")
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
							Kind: types.Listing,
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
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'macros,quotes' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "macros,quotes")
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
							Kind: types.Listing,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "macros,quotes",
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
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'none' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "none")
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
							Kind: types.Listing,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "none",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to https://example.com[] <1>\nand <more text> on the +\n*next* lines with a link to {github-url}[]\n\n* not a list item",
								},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'quotes+' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "quotes+") // same as `quotes,"default"`
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
							Kind: types.Listing,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "quotes+",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to https://example.com[] ",
								},
								&types.Callout{
									Ref: 1,
								},
								&types.StringElement{
									Content: "\nand ",
								},
								&types.SpecialCharacter{
									Name: "<",
								},
								&types.StringElement{
									Content: "more text",
								},
								&types.SpecialCharacter{
									Name: ">",
								},
								&types.StringElement{
									Content: " on the +\n",
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
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the '-callouts' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "-callouts")
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
							Kind: types.Listing,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "-callouts",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to https://example.com[] ",
								},
								&types.SpecialCharacter{
									Name: "<",
								},
								&types.StringElement{
									Content: "1",
								},
								&types.SpecialCharacter{
									Name: ">",
								},
								&types.StringElement{
									Content: "\nand ",
								},
								&types.SpecialCharacter{
									Name: "<",
								},
								&types.StringElement{
									Content: "more text",
								},
								&types.SpecialCharacter{
									Name: ">",
								},
								&types.StringElement{
									Content: " on the +\n*next* lines with a link to {github-url}[]\n\n* not a list item",
								},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the '+quotes' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "+quotes") // default + quotes
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
							Kind: types.Listing,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "+quotes",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to https://example.com[] ",
								},
								&types.Callout{
									Ref: 1,
								},
								&types.StringElement{
									Content: "\nand ",
								},
								&types.SpecialCharacter{
									Name: "<",
								},
								&types.StringElement{
									Content: "more text",
								},
								&types.SpecialCharacter{
									Name: ">",
								},
								&types.StringElement{
									Content: " on the +\n",
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
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the '-quotes' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "-quotes") // default - quotes
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
							Kind: types.Listing,
							Attributes: types.Attributes{
								types.AttrSubstitutions: "-quotes",
							},
							Elements: []interface{}{
								&types.StringElement{
									Content: "a link to https://example.com[] ",
								},
								&types.Callout{
									Ref: 1,
								},
								&types.StringElement{
									Content: "\nand ",
								},
								&types.SpecialCharacter{
									Name: "<",
								},
								&types.StringElement{
									Content: "more text",
								},
								&types.SpecialCharacter{
									Name: ">",
								},
								&types.StringElement{
									Content: " on the +\n*next* lines with a link to {github-url}[]\n\n* not a list item",
								},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should fail when substitution is unknown", func() {
				s := strings.ReplaceAll(source, "$SUBS", "unknown")
				_, err := ParseDocument(s)
				Expect(err).To(HaveOccurred())
			})

			It("should fail when mixing incremental and absolute substitutions", func() {
				s := strings.ReplaceAll(source, "$SUBS", "+attributes,quotes")
				_, err := ParseDocument(s)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("as paragraph blocks", func() {

			source := `:github-url: https://github.com

[listing]
[subs="$SUBS"]
a link to https://example.com[] <1>
and <more text> on the +
*next* lines with a link to {github-url}[]

<1> a callout`

			It("should apply the default substitution", func() {
				s := strings.ReplaceAll(source, "[subs=\"$SUBS\"]\n", "")
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
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle: types.Listing,
							},
							Lines: [][]interface{}{
								{
									&types.StringElement{
										Content: "a link to https://example.com[] ",
									},
									&types.Callout{
										Ref: 1,
									},
								},
								{
									&types.StringElement{
										Content: "\nand ",
									},
									&types.SpecialCharacter{
										Name: "<",
									},
									&types.StringElement{
										Content: "more text",
									},
									&types.SpecialCharacter{
										Name: ">",
									},
									&types.StringElement{
										Content: " on the +",
									},
								},
								{
									&types.StringElement{
										Content: "*next* lines with a link to {github-url}[]",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'quotes' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "quotes")
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
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle:         types.Listing,
								types.AttrSubstitutions: "quotes",
							},
							Lines: [][]interface{}{
								{
									&types.StringElement{
										Content: "a link to https://example.com[] <1>",
									},
								},
								{
									&types.StringElement{
										Content: "and <more text> on the +",
									},
								},
								{
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
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the '+quotes' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "+quotes") // ie, default + quotes
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
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle:         types.Listing,
								types.AttrSubstitutions: "+quotes",
							},
							Lines: [][]interface{}{
								{
									&types.StringElement{
										Content: "a link to https://example.com[] ",
									},
									&types.Callout{
										Ref: 1,
									},
								},
								{
									&types.StringElement{
										Content: "\nand ",
									},
									&types.SpecialCharacter{
										Name: "<",
									},
									&types.StringElement{
										Content: "more text",
									},
									&types.SpecialCharacter{
										Name: ">",
									},
									&types.StringElement{
										Content: " on the +",
									},
								},
								{
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
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})

			It("should apply the 'macros,+quotes,-quotes' substitution", func() {
				s := strings.ReplaceAll(source, "$SUBS", "macros,+quotes,-quotes")
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
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle:         types.Listing,
								types.AttrSubstitutions: "macros,+quotes,-quotes", // ie, "macros" only
							},
							Lines: [][]interface{}{
								{
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
								},
								{
									&types.StringElement{
										Content: "and <more text> on the +",
									},
								},
								{
									&types.StringElement{
										Content: "*next* lines with a link to {github-url}[]",
									},
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
				Expect(ParseDocument(s)).To(MatchDocument(expected))
			})
		})

	})

	Context("passthrough blocks", func() {

		Context("as delimited blocks", func() {

			It("should apply the 'quotes' substitutions", func() {
				source := `[subs=quotes]
.a title
++++
_foo_

*bar*
++++`
				expected := types.Document{
					Attributes: types.Attributes{
						"github-url": "https://github.com",
					},
					Elements: []interface{}{
						types.PassthroughBlock{
							Attributes: types.Attributes{
								types.AttrSubstitutions: "quotes",
								types.AttrTitle:         "a title",
							},
							Lines: [][]interface{}{
								{
									&types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											&types.StringElement{
												Content: "foo",
											},
										},
									},
								},
								{},
								{
									&types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											&types.StringElement{
												Content: "bar",
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
