package parser_test

import (
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("verse blocks", func() {

	Context("in raw documents", func() {

		Context("as delimited blocks", func() {

			It("single line verse with author and title", func() {
				source := `[verse, john doe, verse title]
____
some *verse* content
____`
				expected := []types.DocumentFragment{
					{
						Elements: []interface{}{

							types.VerseBlock{
								Attributes: types.Attributes{
									types.AttrStyle:       types.Verse,
									types.AttrQuoteAuthor: "john doe",
									types.AttrQuoteTitle:  "verse title",
								},
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "some ",
										},
										&types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												&types.StringElement{
													Content: "verse",
												},
											},
										},
										&types.StringElement{
											Content: " content",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("multi-line verse with unrendered list and author only", func() {
				source := `[verse, john doe,   ]
____
- some 
- verse 
- content 
____
`
				expected := []types.DocumentFragment{
					{
						Elements: []interface{}{

							types.VerseBlock{
								Attributes: types.Attributes{
									types.AttrStyle:       types.Verse,
									types.AttrQuoteAuthor: "john doe",
								},
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "- some ",
										},
									},
									{
										&types.StringElement{
											Content: "- verse ",
										},
									},
									{
										&types.StringElement{
											Content: "- content ",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("multi-line verse with title only", func() {
				source := `[verse, ,verse title]
____
some verse content 
____
`
				expected := []types.DocumentFragment{
					{
						Elements: []interface{}{

							types.VerseBlock{
								Attributes: types.Attributes{
									types.AttrStyle:      types.Verse,
									types.AttrQuoteTitle: "verse title",
								},
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "some verse content ",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("multi-line verse with unrendered lists and block without author and title", func() {
				source := `[verse]
____
* some
----
* verse 
----
* content
____`
				expected := []types.DocumentFragment{
					{
						Elements: []interface{}{

							types.VerseBlock{
								Attributes: types.Attributes{
									types.AttrStyle: types.Verse,
								},
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "* some",
										},
									},
									{
										&types.StringElement{
											Content: "----",
										},
									},
									{
										&types.StringElement{
											Content: "* verse ",
										},
									},
									{
										&types.StringElement{
											Content: "----",
										},
									},
									{
										&types.StringElement{
											Content: "* content",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("multi-line verse with unrendered list without author and title", func() {
				source := `[verse]
____
* foo


	* bar
____`
				expected := []types.DocumentFragment{
					{
						Elements: []interface{}{

							types.VerseBlock{
								Attributes: types.Attributes{
									types.AttrStyle: types.Verse,
								},
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "* foo",
										},
									},
									{},
									{},
									{
										&types.StringElement{
											Content: "\t* bar",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("empty verse without author and title", func() {
				source := `[verse]
____
____`
				expected := []types.DocumentFragment{
					{
						Elements: []interface{}{

							types.VerseBlock{
								Attributes: types.Attributes{
									types.AttrStyle: types.Verse,
								},
								Lines: [][]interface{}{
									{},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			It("unclosed verse without author and title", func() {
				source := `[verse]
____
foo
`
				expected := []types.DocumentFragment{
					{
						Elements: []interface{}{

							types.VerseBlock{
								Attributes: types.Attributes{
									types.AttrStyle: types.Verse,
								},
								Lines: [][]interface{}{
									{
										&types.StringElement{
											Content: "foo",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragmentGroups(expected))
			})

			Context("with custom substitutions", func() {

				source := `:github-url: https://github.com
				
[subs="$SUBS"]
[verse, john doe, verse title]
____
a link to https://example.com[] <1>
and <more text> on the +
*next* lines with a link to {github-url}[]

* not a list item
____

<1> a callout
`

				It("should apply the default substitution", func() {
					s := strings.ReplaceAll(source, "[subs=\"$SUBS\"]\n", "")
					expected := []types.DocumentFragment{
						{
							Elements: []interface{}{

								&types.AttributeDeclaration{
									Name:  "github-url",
									Value: "https://github.com",
								},
								&types.BlankLine{},
								types.VerseBlock{
									Attributes: types.Attributes{
										types.AttrStyle:       types.Verse,
										types.AttrQuoteAuthor: "john doe",
										types.AttrQuoteTitle:  "verse title",
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
										},
										{
											&types.StringElement{
												Content: "and ",
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
										{},
										{
											&types.StringElement{
												Content: "* not a list item",
											},
										},
									},
								},
								&types.BlankLine{},
								types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
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
					Expect(AssembleDocumentFragments(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'normal' substitution", func() {
					s := strings.ReplaceAll(source, "$SUBS", "normal")
					expected := []types.DocumentFragment{
						{
							Elements: []interface{}{
								&types.AttributeDeclaration{
									Name:  "github-url",
									Value: "https://github.com",
								},
								&types.BlankLine{},
								types.VerseBlock{
									Attributes: types.Attributes{
										types.AttrStyle:         types.Verse,
										types.AttrQuoteAuthor:   "john doe",
										types.AttrQuoteTitle:    "verse title",
										types.AttrSubstitutions: "normal",
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
										},
										{
											&types.StringElement{
												Content: "and ",
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
										{},
										{
											&types.StringElement{
												Content: "* not a list item",
											},
										},
									},
								},
								&types.BlankLine{},
								types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
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
					Expect(AssembleDocumentFragments(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'quotes' substitution", func() {
					s := strings.ReplaceAll(source, "$SUBS", "quotes")
					expected := []types.DocumentFragment{
						{
							Elements: []interface{}{
								&types.AttributeDeclaration{
									Name:  "github-url",
									Value: "https://github.com",
								},
								&types.BlankLine{},
								types.VerseBlock{
									Attributes: types.Attributes{
										types.AttrStyle:         types.Verse,
										types.AttrQuoteAuthor:   "john doe",
										types.AttrQuoteTitle:    "verse title",
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
										{},
										{
											&types.StringElement{
												Content: "* not a list item",
											},
										},
									},
								},
								&types.BlankLine{},
								types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
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
					Expect(AssembleDocumentFragments(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'macros' substitution", func() {
					s := strings.ReplaceAll(source, "$SUBS", "macros")
					expected := []types.DocumentFragment{
						{
							Elements: []interface{}{
								&types.AttributeDeclaration{
									Name:  "github-url",
									Value: "https://github.com",
								},
								&types.BlankLine{},
								types.VerseBlock{
									Attributes: types.Attributes{
										types.AttrStyle:         types.Verse,
										types.AttrQuoteAuthor:   "john doe",
										types.AttrQuoteTitle:    "verse title",
										types.AttrSubstitutions: "macros",
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
										{},
										{
											&types.StringElement{
												Content: "* not a list item",
											},
										},
									},
								},
								&types.BlankLine{},
								types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
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
					Expect(AssembleDocumentFragments(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'attributes' substitution", func() {
					s := strings.ReplaceAll(source, "$SUBS", "attributes")
					expected := []types.DocumentFragment{
						{
							Elements: []interface{}{
								&types.AttributeDeclaration{
									Name:  "github-url",
									Value: "https://github.com",
								},
								&types.BlankLine{},
								types.VerseBlock{
									Attributes: types.Attributes{
										types.AttrStyle:         types.Verse,
										types.AttrQuoteAuthor:   "john doe",
										types.AttrQuoteTitle:    "verse title",
										types.AttrSubstitutions: "attributes",
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
											&types.StringElement{
												Content: "*next* lines with a link to https://github.com[]",
											},
										},
										{},
										{
											&types.StringElement{
												Content: "* not a list item",
											},
										},
									},
								},
								&types.BlankLine{},
								types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
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
					Expect(AssembleDocumentFragments(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'attributes,macros' substitution", func() {
					s := strings.ReplaceAll(source, "$SUBS", "attributes,macros")
					expected := []types.DocumentFragment{
						{
							Elements: []interface{}{
								&types.AttributeDeclaration{
									Name:  "github-url",
									Value: "https://github.com",
								},
								&types.BlankLine{},
								types.VerseBlock{
									Attributes: types.Attributes{
										types.AttrStyle:         types.Verse,
										types.AttrQuoteAuthor:   "john doe",
										types.AttrQuoteTitle:    "verse title",
										types.AttrSubstitutions: "attributes,macros",
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
												Content: "*next* lines with a link to ",
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
										{},
										{
											&types.StringElement{
												Content: "* not a list item",
											},
										},
									},
								},
								&types.BlankLine{},
								types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
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
					Expect(AssembleDocumentFragments(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'specialchars' substitution", func() {
					s := strings.ReplaceAll(source, "$SUBS", "specialchars")
					expected := []types.DocumentFragment{
						{
							Elements: []interface{}{
								&types.AttributeDeclaration{
									Name:  "github-url",
									Value: "https://github.com",
								},
								&types.BlankLine{},
								types.VerseBlock{
									Attributes: types.Attributes{
										types.AttrStyle:         types.Verse,
										types.AttrQuoteAuthor:   "john doe",
										types.AttrQuoteTitle:    "verse title",
										types.AttrSubstitutions: "specialchars",
									},
									Lines: [][]interface{}{
										{
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
										},
										{
											&types.StringElement{
												Content: "and ",
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
										{},
										{
											&types.StringElement{
												Content: "* not a list item",
											},
										},
									},
								},
								&types.BlankLine{},
								types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
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
					Expect(AssembleDocumentFragments(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'replacements' substitution", func() {
					s := strings.ReplaceAll(source, "$SUBS", "replacements")
					expected := []types.DocumentFragment{
						{
							Elements: []interface{}{

								&types.AttributeDeclaration{
									Name:  "github-url",
									Value: "https://github.com",
								},
								&types.BlankLine{},
								types.VerseBlock{
									Attributes: types.Attributes{
										types.AttrStyle:         types.Verse,
										types.AttrQuoteAuthor:   "john doe",
										types.AttrQuoteTitle:    "verse title",
										types.AttrSubstitutions: "replacements",
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
											&types.StringElement{
												Content: "*next* lines with a link to {github-url}[]",
											},
										},
										{},
										{
											&types.StringElement{
												Content: "* not a list item",
											},
										},
									},
								},
								&types.BlankLine{},
								types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
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
					Expect(AssembleDocumentFragments(s)).To(MatchDocumentFragmentGroups(expected))
				})

				It("should apply the 'post_replacements' substitution", func() {
					s := strings.ReplaceAll(source, "$SUBS", "post_replacements")
					expected := []types.DocumentFragment{
						{
							Elements: []interface{}{
								&types.AttributeDeclaration{
									Name:  "github-url",
									Value: "https://github.com",
								},
								&types.BlankLine{},
								types.VerseBlock{
									Attributes: types.Attributes{
										types.AttrStyle:         types.Verse,
										types.AttrQuoteAuthor:   "john doe",
										types.AttrQuoteTitle:    "verse title",
										types.AttrSubstitutions: "post_replacements",
									},
									Lines: [][]interface{}{
										{
											&types.StringElement{
												Content: "a link to https://example.com[] <1>",
											},
										},
										{
											&types.StringElement{
												Content: "and <more text> on the",
											},
											types.LineBreak{},
										},
										{
											&types.StringElement{
												Content: "*next* lines with a link to {github-url}[]",
											},
										},
										{},
										{
											&types.StringElement{
												Content: "* not a list item",
											},
										},
									},
								},
								&types.BlankLine{},
								types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
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
					Expect(AssembleDocumentFragments(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'quotes,macros' substitution", func() {
					s := strings.ReplaceAll(source, "$SUBS", "quotes,macros")
					expected := []types.DocumentFragment{
						{
							Elements: []interface{}{
								&types.AttributeDeclaration{
									Name:  "github-url",
									Value: "https://github.com",
								},
								&types.BlankLine{},
								types.VerseBlock{
									Attributes: types.Attributes{
										types.AttrStyle:         types.Verse,
										types.AttrQuoteAuthor:   "john doe",
										types.AttrQuoteTitle:    "verse title",
										types.AttrSubstitutions: "quotes,macros",
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
										{},
										{
											&types.StringElement{
												Content: "* not a list item",
											},
										},
									},
								},
								&types.BlankLine{},
								types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
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
					Expect(AssembleDocumentFragments(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'macros,quotes' substitution", func() {
					s := strings.ReplaceAll(source, "$SUBS", "macros,quotes")
					expected := []types.DocumentFragment{
						{
							Elements: []interface{}{
								&types.AttributeDeclaration{
									Name:  "github-url",
									Value: "https://github.com",
								},
								&types.BlankLine{},
								types.VerseBlock{
									Attributes: types.Attributes{
										types.AttrStyle:         types.Verse,
										types.AttrQuoteAuthor:   "john doe",
										types.AttrQuoteTitle:    "verse title",
										types.AttrSubstitutions: "macros,quotes",
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
										{},
										{
											&types.StringElement{
												Content: "* not a list item",
											},
										},
									},
								},
								&types.BlankLine{},
								types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
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
					Expect(AssembleDocumentFragments(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'none' substitution", func() {
					s := strings.ReplaceAll(source, "$SUBS", "none")
					expected := []types.DocumentFragment{
						{
							Elements: []interface{}{
								&types.AttributeDeclaration{
									Name:  "github-url",
									Value: "https://github.com",
								},
								&types.BlankLine{},
								types.VerseBlock{
									Attributes: types.Attributes{
										types.AttrStyle:         types.Verse,
										types.AttrQuoteAuthor:   "john doe",
										types.AttrQuoteTitle:    "verse title",
										types.AttrSubstitutions: "none",
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
											&types.StringElement{
												Content: "*next* lines with a link to {github-url}[]",
											},
										},
										{},
										{
											&types.StringElement{
												Content: "* not a list item",
											},
										},
									},
								},
								&types.BlankLine{},
								types.CalloutListElement{
									Ref: 1,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
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
					Expect(AssembleDocumentFragments(s)).To(MatchDocumentFragments(expected))
				})
			})
		})
	})

	Context("in final documents", func() {

		Context("as delimited blocks", func() {

			It("single line verse with author and title", func() {
				source := `[verse, john doe, verse title]
____
some *verse* content
____`
				expected := types.Document{
					Elements: []interface{}{
						types.VerseBlock{
							Attributes: types.Attributes{
								types.AttrStyle:       types.Verse,
								types.AttrQuoteAuthor: "john doe",
								types.AttrQuoteTitle:  "verse title",
							},
							Lines: [][]interface{}{
								{
									&types.StringElement{
										Content: "some ",
									},
									&types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											&types.StringElement{
												Content: "verse",
											},
										},
									},
									&types.StringElement{
										Content: " content",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("multi-line verse with unrendered list and author only", func() {
				source := `[verse, john doe,   ]
____
- some 
- verse 
- content 
____
`
				expected := types.Document{
					Elements: []interface{}{
						types.VerseBlock{
							Attributes: types.Attributes{
								types.AttrStyle:       types.Verse,
								types.AttrQuoteAuthor: "john doe",
							},
							Lines: [][]interface{}{
								{
									&types.StringElement{
										Content: "- some ",
									},
								},
								{
									&types.StringElement{
										Content: "- verse ",
									},
								},
								{
									&types.StringElement{
										Content: "- content ",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("multi-line verse with title only", func() {
				source := `[verse, ,verse title]
____
some verse content 
____
`
				expected := types.Document{
					Elements: []interface{}{
						types.VerseBlock{
							Attributes: types.Attributes{
								types.AttrStyle:      types.Verse,
								types.AttrQuoteTitle: "verse title",
							},
							Lines: [][]interface{}{
								{
									&types.StringElement{
										Content: "some verse content ",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("multi-line verse with unrendered lists and block without author and title", func() {
				source := `[verse]
____
* some
----
* verse 
----
* content
____`
				expected := types.Document{
					Elements: []interface{}{
						types.VerseBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Verse,
							},
							Lines: [][]interface{}{
								{
									&types.StringElement{
										Content: "* some",
									},
								},
								{
									&types.StringElement{
										Content: "----",
									},
								},
								{
									&types.StringElement{
										Content: "* verse ",
									},
								},
								{
									&types.StringElement{
										Content: "----",
									},
								},
								{
									&types.StringElement{
										Content: "* content",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("multi-line verse with unrendered list without author and title", func() {
				source := `[verse]
____
* foo


	* bar
____`
				expected := types.Document{
					Elements: []interface{}{
						types.VerseBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Verse,
							},
							Lines: [][]interface{}{
								{
									&types.StringElement{
										Content: "* foo",
									},
								},
								{},
								{},
								{
									&types.StringElement{
										Content: "\t* bar",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("empty verse without author and title", func() {
				source := `[verse]
____
____`
				expected := types.Document{
					Elements: []interface{}{
						types.VerseBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Verse,
							},
							Lines: [][]interface{}{
								{},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("unclosed verse without author and title", func() {
				source := `[verse]
____
foo
`
				expected := types.Document{
					Elements: []interface{}{
						types.VerseBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Verse,
							},
							Lines: [][]interface{}{
								{
									&types.StringElement{
										Content: "foo",
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
