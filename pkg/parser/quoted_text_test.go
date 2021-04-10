package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("quoted texts", func() {

	Context("in raw documents", func() {

		Context("with single punctuations", func() {

			It("bold text with newline", func() {
				source := "*some bold\ncontent*"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{
											Content: "some bold\ncontent",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("italic text with 3 words in single quote", func() {
				source := "_some italic content_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{
											Content: "some italic content",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("italic text with newline", func() {
				source := "_some italic\ncontent_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{
											Content: "some italic\ncontent",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("monospace text with 3 words", func() {
				source := "`some monospace content`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{
											Content: "some monospace content",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("monospace text with newline", func() {
				source := "`some monospace\ncontent`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{
											Content: "some monospace\ncontent",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("invalid subscript text with 3 words", func() {
				source := "~some subscript content~"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{
									Content: "~some subscript content~",
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("invalid superscript text with 3 words", func() {
				source := "^some superscript content^"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{
									Content: "^some superscript content^",
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("bold text within italic text", func() {
				source := "_some *bold* content_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "some "},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "bold"},
											},
										},
										types.StringElement{Content: " content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("monospace text within bold text within italic quote", func() {
				source := "*some _italic and `monospaced content`_*"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "some "},
										types.QuotedText{
											Kind: types.SingleQuoteItalic,
											Elements: []interface{}{
												types.StringElement{Content: "italic and "},
												types.QuotedText{
													Kind: types.SingleQuoteMonospace,
													Elements: []interface{}{
														types.StringElement{
															Content: "monospaced content",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("italic text within italic text", func() {
				source := "_some _very italic_ content_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "some _very italic"},
									},
								},
								types.StringElement{Content: " content_"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("bold delimiter text within bold text", func() {
				source := "*bold*content*"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "bold*content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("italic delimiter text within italic text", func() {
				source := "_italic_content_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "italic_content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("monospace delimiter text within monospace text", func() {
				source := "`monospace`content`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "monospace`content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("non-bold text then bold text", func() {
				source := "non*bold*content *bold content*"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{
									Content: "non*bold*content ",
								},
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "bold content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})
			It("non-italic text then italic text", func() {
				source := "non_italic_content _italic content_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{
									Content: "non_italic_content ",
								},
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "italic content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("non-monospace text then monospace text", func() {
				source := "non`monospace`content `monospace content`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{
									Content: "non`monospace`content ",
								},
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "monospace content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("subscript text attached", func() {
				source := "O~2~ is a molecule"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "O"},
								types.QuotedText{
									Kind: types.SingleQuoteSubscript,
									Elements: []interface{}{
										types.StringElement{Content: "2"},
									},
								},
								types.StringElement{Content: " is a molecule"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("superscript text attached", func() {
				source := "M^me^ White"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "M"},
								types.QuotedText{
									Kind: types.SingleQuoteSuperscript,
									Elements: []interface{}{
										types.StringElement{Content: "me"},
									},
								},
								types.StringElement{Content: " White"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("invalid subscript text with 3 words", func() {
				source := "~some subscript content~"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "~some subscript content~"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})
		})

		Context("with double punctuations", func() {

			It("bold text of 1 word in double quote", func() {
				source := "**hello**"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "hello"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("bold text with newline", func() {
				source := "**some bold\ncontent**"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "some bold\ncontent"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("italic text with 3 words in double quote", func() {
				source := "__some italic content__"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "some italic content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("italic text with newline", func() {
				source := "__some italic\ncontent__"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "some italic\ncontent"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("monospace text with 3 words in double quote", func() {
				source := "`` some monospace content ``"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: " some monospace content "},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("monospace text with newline", func() {
				source := "``some monospace\ncontent``"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "some monospace\ncontent"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("superscript text within italic text", func() {
				source := "__some ^superscript^ content__"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "some "},
										types.QuotedText{
											Kind: types.SingleQuoteSuperscript,
											Elements: []interface{}{
												types.StringElement{Content: "superscript"},
											},
										},
										types.StringElement{Content: " content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("superscript text within italic text within bold quote", func() {
				source := "**some _italic and ^superscriptcontent^_**"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "some "},
										types.QuotedText{
											Kind: types.SingleQuoteItalic,
											Elements: []interface{}{
												types.StringElement{Content: "italic and "},
												types.QuotedText{
													Kind: types.SingleQuoteSuperscript,
													Elements: []interface{}{
														types.StringElement{Content: "superscriptcontent"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})
		})

		Context("inline", func() {

			It("inline content with bold text", func() {
				source := "a paragraph with *some bold content*"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "a paragraph with "},
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "some bold content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("inline content with invalid bold text - use case 1", func() {
				source := "a paragraph with *some bold content"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "a paragraph with *some bold content"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("inline content with invalid bold text - use case 2", func() {
				source := "a paragraph with *some bold content *"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "a paragraph with *some bold content *"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("inline content with invalid bold text - use case 3", func() {
				source := "a paragraph with * some bold content*"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "a paragraph with * some bold content*"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("invalid italic text within bold text", func() {
				source := "some *bold and _italic content _ together*."
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "some "},
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "bold and _italic content _ together"},
									},
								},
								types.StringElement{Content: "."},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("italic text within invalid bold text", func() {
				source := "some *bold and _italic content_ together *."
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "some *bold and "},
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "italic content"},
									},
								},
								types.StringElement{Content: " together *."},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("inline content with invalid subscript text - use case 1", func() {
				source := "a paragraph with ~some subscript content"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "a paragraph with ~some subscript content"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("inline content with invalid subscript text - use case 2", func() {
				source := "a paragraph with ~some subscript content ~"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "a paragraph with ~some subscript content ~"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("inline content with invalid subscript text - use case 3", func() {
				source := "a paragraph with ~ some subscript content~"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "a paragraph with ~ some subscript content~"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("inline content with invalid superscript text - use case 1", func() {
				source := "a paragraph with ^some superscript content"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "a paragraph with ^some superscript content"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("inline content with invalid superscript text - use case 2", func() {
				source := "a paragraph with ^some superscript content ^"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "a paragraph with ^some superscript content ^"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("inline content with invalid superscript text - use case 3", func() {
				source := "a paragraph with ^ some superscript content^"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "a paragraph with ^ some superscript content^"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("marked text with newline", func() {
				source := "#some marked\ncontent#"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMarked,
									Elements: []interface{}{
										types.StringElement{Content: "some marked\ncontent"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("double marked text with newline", func() {
				source := "##some marked\ncontent##"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteMarked,
									Elements: []interface{}{
										types.StringElement{Content: "some marked\ncontent"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

		})

		Context("with attributes", func() {

			It("simple dot.role italics", func() {
				source := "[.myrole]_italics_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "italics"},
									},
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{"myrole"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("simple dot.role italics unconstrained", func() {
				source := "it[.uncle]__al__ic"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{
									Content: "it",
								},
								types.QuotedText{
									Kind: types.DoubleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "al"},
									},
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{"uncle"},
									},
								},
								types.StringElement{
									Content: "ic",
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("simple dot.role bold", func() {
				source := "[.myrole]*bold*"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "bold"},
									},
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{"myrole"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("simple dot.role bold unconstrained", func() {
				source := "it[.uncle]**al**ic"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{
									Content: "it",
								},
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "al"},
									},
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{"uncle"},
									},
								},
								types.StringElement{
									Content: "ic",
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("simple dot.role mono", func() {
				source := "[.myrole]`true`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "true"},
									},
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{"myrole"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("simple dot.role mono unconstrained", func() {
				source := "int[.uncle]``eg``rate"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{
									Content: "int",
								},
								types.QuotedText{
									Kind: types.DoubleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "eg"},
									},
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{"uncle"},
									},
								},
								types.StringElement{
									Content: "rate",
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("role with extra attribute", func() {
				source := "[myrole,and=nothing]_italics_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{"myrole"},
										"and":           "nothing",
									},
									Elements: []interface{}{
										types.StringElement{Content: "italics"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("dot.role with extra attribute", func() {
				source := "[.myrole,and=nothing]_italics_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{"myrole"},
										"and":           "nothing",
									},
									Elements: []interface{}{
										types.StringElement{Content: "italics"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("short-hand ID only", func() {
				source := "[#here]*bold*"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "bold"},
									},
									Attributes: types.Attributes{
										types.AttrID: "here",
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("short-hand role only", func() {
				source := "[bob]**bold**"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "bold"},
									},
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{"bob"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("short-hand dot.role only", func() {
				source := "[.bob]**bold**"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "bold"},
									},
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{"bob"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("short-hand role with special characters", func() {
				source := `["a <role>"]**bold**`
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{
											[]interface{}{
												types.StringElement{
													Content: "a ",
												},
												types.SpecialCharacter{
													Name: "<",
												},
												types.StringElement{
													Content: "role",
												},
												types.SpecialCharacter{
													Name: ">",
												},
											},
										},
									},
									Elements: []interface{}{
										types.StringElement{Content: "bold"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("short-hand dot.role with special characters", func() {
				source := `[."a <role>"]**bold**`
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{
											[]interface{}{
												types.StringElement{
													Content: "a ",
												},
												types.SpecialCharacter{
													Name: "<",
												},
												types.StringElement{
													Content: "role",
												},
												types.SpecialCharacter{
													Name: ">",
												},
											},
										},
									},
									Elements: []interface{}{
										types.StringElement{Content: "bold"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("marked short-hand dot.role only", func() {
				source := "[.bob]##the builder##"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteMarked,
									Elements: []interface{}{
										types.StringElement{Content: "the builder"},
									},
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{"bob"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("short-hand multiple roles and id", func() {
				source := "[.role1#anchor.role2.role3]**bold**[#here.second.class]_text_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{"role1", "role2", "role3"},
										types.AttrID:    "anchor",
									},
									Elements: []interface{}{
										types.StringElement{Content: "bold"},
									},
								},
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{"second", "class"},
										types.AttrID:    "here",
									},
									Elements: []interface{}{
										types.StringElement{Content: "text"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("empty role", func() {
				source := "[]**bold**"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "bold"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("quoted dot.role", func() {
				source := "[.'here, again']**bold**"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "bold"},
									},
									// NB: This will confuse the renderer.
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{"here, again"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("quoted dot.role with special chars", func() {
				source := "[.\"something <wicked>\"]**bold**"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "bold"},
									},
									// NB: This will confuse the renderer.
									Attributes: types.Attributes{
										types.AttrRoles: []interface{}{
											[]interface{}{
												types.StringElement{
													Content: "something ",
												},
												types.SpecialCharacter{
													Name: "<",
												},
												types.StringElement{
													Content: "wicked",
												},
												types.SpecialCharacter{
													Name: ">",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			// This demonstrates that we cannot inject malicious data in these attributes.
			// The content is escaped by the renderer, not the parser.
			It("bad syntax", func() {
				source := "[.<something \"wicked>]**bold**"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{
									Content: "[.",
								},
								types.SpecialCharacter{
									Name: "<",
								},
								types.StringElement{
									Content: "something \"wicked",
								},
								types.SpecialCharacter{
									Name: ">",
								},
								types.StringElement{
									Content: "]",
								},
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "bold"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

		})

		Context("with nested quoted text", func() {

			It("italic text within bold text", func() {
				source := "some *bold and _italic content_ together*."
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "some "},
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "bold and "},
										types.QuotedText{
											Kind: types.SingleQuoteItalic,
											Elements: []interface{}{
												types.StringElement{Content: "italic content"},
											},
										},
										types.StringElement{Content: " together"},
									},
								},
								types.StringElement{Content: "."},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("single-quote bold within single-quote bold text", func() {
				// here we don't allow for bold text within bold text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "*some *nested bold* content*"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "some *nested bold"},
									},
								},
								types.StringElement{Content: " content*"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("double-quote bold within double-quote bold text", func() {
				// here we don't allow for bold text within bold text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "**some **nested bold** content**"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "some "},
									},
								},
								types.StringElement{Content: "nested bold"},
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: " content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("single-quote bold within double-quote bold text", func() {
				// here we don't allow for bold text within bold text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "**some *nested bold* content**"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "some "},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "nested bold"},
											},
										},
										types.StringElement{Content: " content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("double-quote bold within single-quote bold text", func() {
				// here we don't allow for bold text within bold text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "*some **nested bold** content*"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "some "},
										types.QuotedText{
											Kind: types.DoubleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "nested bold"},
											},
										},
										types.StringElement{Content: " content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("single-quote italic within single-quote italic text", func() {
				// here we don't allow for italic text within italic text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "_some _nested italic_ content_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "some _nested italic"},
									},
								},
								types.StringElement{Content: " content_"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("double-quote italic within double-quote italic text", func() {
				// here we don't allow for italic text within italic text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "__some __nested italic__ content__"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "some "},
									},
								},
								types.StringElement{Content: "nested italic"},
								types.QuotedText{
									Kind: types.DoubleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: " content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("single-quote italic within double-quote italic text", func() {
				// here we allow for italic text within italic text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "_some __nested italic__ content_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "some "},
										types.QuotedText{
											Kind: types.DoubleQuoteItalic,
											Elements: []interface{}{
												types.StringElement{Content: "nested italic"},
											},
										},
										types.StringElement{Content: " content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("double-quote italic within single-quote italic text", func() {
				// here we allow for italic text within italic text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "_some __nested italic__ content_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "some "},
										types.QuotedText{
											Kind: types.DoubleQuoteItalic,
											Elements: []interface{}{
												types.StringElement{Content: "nested italic"},
											},
										},
										types.StringElement{Content: " content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("single-quote monospace within single-quote monospace text", func() {
				// here we don't allow for monospace text within monospace text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "`some `nested monospace` content`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "some `nested monospace"},
									},
								},
								types.StringElement{Content: " content`"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("double-quote monospace within double-quote monospace text", func() {
				// here we don't allow for monospace text within monospace text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "``some ``nested monospace`` content``"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.DoubleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "some "},
									},
								},
								types.StringElement{Content: "nested monospace"},
								types.QuotedText{
									Kind: types.DoubleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: " content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("single-quote monospace within double-quote monospace text", func() {
				// here we allow for monospace text within monospace text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "`some ``nested monospace`` content`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "some "},
										types.QuotedText{
											Kind: types.DoubleQuoteMonospace,
											Elements: []interface{}{
												types.StringElement{Content: "nested monospace"},
											},
										},
										types.StringElement{Content: " content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("double-quote monospace within single-quote monospace text", func() {
				// here we allow for monospace text within monospace text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "`some ``nested monospace`` content`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "some "},
										types.QuotedText{
											Kind: types.DoubleQuoteMonospace,
											Elements: []interface{}{
												types.StringElement{Content: "nested monospace"},
											},
										},
										types.StringElement{Content: " content"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("quoted text within marked text", func() {
				source := "some #marked and _italic_ and *bold* and `monospaced` content together#."
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "some "},
								types.QuotedText{
									Kind: types.SingleQuoteMarked,
									Elements: []interface{}{
										types.StringElement{Content: "marked and "},
										types.QuotedText{
											Kind: types.SingleQuoteItalic,
											Elements: []interface{}{
												types.StringElement{Content: "italic"},
											},
										},
										types.StringElement{Content: " and "},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "bold"},
											},
										},
										types.StringElement{Content: " and "},
										types.QuotedText{
											Kind: types.SingleQuoteMonospace,
											Elements: []interface{}{
												types.StringElement{Content: "monospaced"},
											},
										},
										types.StringElement{Content: " content together"},
									},
								},
								types.StringElement{Content: "."},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("unbalanced bold in monospace - case 1", func() {
				source := "`*a`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "*a"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("unbalanced bold in monospace - case 2", func() {
				source := "`a*b`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "a*b"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("italic in monospace", func() {
				source := "`_a_`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.QuotedText{
											Kind: types.SingleQuoteItalic,
											Elements: []interface{}{
												types.StringElement{Content: "a"},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("unbalanced italic in monospace", func() {
				source := "`a_b`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "a_b"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("unparsed bold in monospace", func() {
				source := "`a*b*`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "a*b*"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("parsed subscript in monospace", func() {
				source := "`a~b~`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "a"},
										types.QuotedText{
											Kind: types.SingleQuoteSubscript,
											Elements: []interface{}{
												types.StringElement{Content: "b"},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("multiline in single quoted monospace - case 1", func() {
				source := "`a\nb`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "a\nb"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("multiline in double quoted monospace - case 1", func() {
				source := "`a\nb`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "a\nb"},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("multiline in single quoted  monospace - case 2", func() {
				source := "`a\n*b*`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "a\n"},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "b"},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("multiline in double quoted  monospace - case 2", func() {
				source := "`a\n*b*`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "a\n"},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "b"},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("link in bold", func() {
				source := "*a link:/[b]*"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "a "},
										types.InlineLink{
											Attributes: types.Attributes{
												types.AttrInlineLinkText: "b",
											},
											Location: types.Location{
												Path: []interface{}{
													types.StringElement{
														Content: "/",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("image in bold", func() {
				source := "*a image:foo.png[]*"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "a "},
										types.InlineImage{
											Location: types.Location{
												Path: []interface{}{
													types.StringElement{
														Content: "foo.png",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("singleplus passthrough in bold", func() {
				source := "*a +image:foo.png[]+*"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "a "},
										types.InlinePassthrough{
											Kind: types.SinglePlusPassthrough,
											Elements: []interface{}{
												types.StringElement{Content: "image:foo.png[]"},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("tripleplus passthrough in bold", func() {
				source := "*a +++image:foo.png[]+++*"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{Content: "a "},
										types.InlinePassthrough{
											Kind: types.TriplePlusPassthrough,
											Elements: []interface{}{
												types.StringElement{Content: "image:foo.png[]"},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("link in italic", func() {
				source := "_a link:/[b]_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "a "},
										types.InlineLink{
											Attributes: types.Attributes{
												types.AttrInlineLinkText: "b",
											},
											Location: types.Location{
												Path: []interface{}{
													types.StringElement{
														Content: "/",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("image in italic", func() {
				source := "_a image:foo.png[]_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "a "},
										types.InlineImage{
											Location: types.Location{
												Path: []interface{}{
													types.StringElement{
														Content: "foo.png",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("singleplus passthrough in italic", func() {
				source := "_a +image:foo.png[]+_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "a "},
										types.InlinePassthrough{
											Kind: types.SinglePlusPassthrough,
											Elements: []interface{}{
												types.StringElement{Content: "image:foo.png[]"},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("tripleplus passthrough in italic", func() {
				source := "_a +++image:foo.png[]+++_"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{Content: "a "},
										types.InlinePassthrough{
											Kind: types.TriplePlusPassthrough,
											Elements: []interface{}{
												types.StringElement{Content: "image:foo.png[]"},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("link in monospace", func() {
				source := "`a link:/[b]`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "a "},
										types.InlineLink{
											Attributes: types.Attributes{
												types.AttrInlineLinkText: "b",
											},
											Location: types.Location{
												Path: []interface{}{
													types.StringElement{
														Content: "/",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("image in monospace", func() {
				source := "`a image:foo.png[]`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "a "},
										types.InlineImage{
											Location: types.Location{
												Path: []interface{}{
													types.StringElement{
														Content: "foo.png",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("singleplus passthrough in monospace", func() {
				source := "`a +image:foo.png[]+`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "a "},
										types.InlinePassthrough{
											Kind: types.SinglePlusPassthrough,
											Elements: []interface{}{
												types.StringElement{Content: "image:foo.png[]"},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("tripleplus passthrough in monospace", func() {
				source := "`a +++image:foo.png[]+++`"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{Content: "a "},
										types.InlinePassthrough{
											Kind: types.TriplePlusPassthrough,
											Elements: []interface{}{
												types.StringElement{Content: "image:foo.png[]"},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

		})

		Context("unbalanced quoted text", func() {

			Context("unbalanced bold text", func() {

				It("unbalanced bold text - extra on left", func() {
					source := "**some bold content*"
					expected := types.DocumentFragments{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "*some bold content"},
										},
									},
								},
							},
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("unbalanced bold text - extra on right", func() {
					source := "*some bold content**"
					expected := types.DocumentFragments{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "some bold content"},
										},
									},
									types.StringElement{Content: "*"},
								},
							},
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})
			})

			Context("unbalanced italic text", func() {

				It("unbalanced italic text - extra on left", func() {
					source := "__some italic content_"
					expected := types.DocumentFragments{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "_some italic content"},
										},
									},
								},
							},
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("unbalanced italic text - extra on right", func() {
					source := "_some italic content__"
					expected := types.DocumentFragments{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "some italic content"},
										},
									},
									types.StringElement{Content: "_"},
								},
							},
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})
			})

			Context("unbalanced monospace text", func() {

				It("unbalanced monospace text - extra on left", func() {
					source := "``some monospace content`"
					expected := types.DocumentFragments{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "`some monospace content"},
										},
									},
								},
							},
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("unbalanced monospace text - extra on right", func() {
					source := "`some monospace content``"
					expected := types.DocumentFragments{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "some monospace content"},
										},
									},
									types.StringElement{Content: "`"},
								},
							},
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})
			})

			It("inline content with unbalanced bold text", func() {
				source := "a paragraph with *some bold content"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "a paragraph with *some bold content"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

		})

		Context("prevented substitution", func() {

			Context("prevented bold text substitution", func() {

				Context("without nested quoted text", func() {

					It("escaped bold text with single backslash", func() {
						source := `\*bold content*`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "*bold content*"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped bold text with multiple backslashes", func() {
						source := `\\*bold content*`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `\*bold content*`},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped bold text with double quote", func() {
						source := `\\**bold content**`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `**bold content**`},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped bold text with double quote and more backslashes", func() {
						source := `\\\**bold content**`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `\**bold content**`},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped bold text with unbalanced double quote", func() {
						source := `\**bold content*`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `**bold content*`},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped bold text with unbalanced double quote and more backslashes", func() {
						source := `\\\**bold content*`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `\\**bold content*`},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})
				})

				Context("with nested quoted text", func() {

					It("escaped bold text with nested italic text", func() {
						source := `\*_italic content_*`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "*"},
										types.QuotedText{
											Kind: types.SingleQuoteItalic,
											Elements: []interface{}{
												types.StringElement{Content: "italic content"},
											},
										},
										types.StringElement{Content: "*"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped bold text with unbalanced double quote and nested italic test", func() {
						source := `\**_italic content_*`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "**"},
										types.QuotedText{
											Kind: types.SingleQuoteItalic,
											Elements: []interface{}{
												types.StringElement{Content: "italic content"},
											},
										},
										types.StringElement{Content: "*"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped bold text with nested italic", func() {
						source := `\*bold _and italic_ content*`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "*bold "},
										types.QuotedText{
											Kind: types.SingleQuoteItalic,
											Elements: []interface{}{
												types.StringElement{Content: "and italic"},
											},
										},
										types.StringElement{Content: " content*"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})
				})

			})

			Context("prevented italic text substitution", func() {

				Context("without nested quoted text", func() {

					It("escaped italic text with single quote", func() {
						source := `\_italic content_`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "_italic content_"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped italic text with single quote and more backslashes", func() {
						source := `\\_italic content_`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `\_italic content_`},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped italic text with double quote with 2 backslashes", func() {
						source := `\\__italic content__`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `__italic content__`},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped italic text with double quote with 3 backslashes", func() {
						source := `\\\__italic content__`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `\__italic content__`},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped italic text with unbalanced double quote", func() {
						source := `\__italic content_`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `__italic content_`},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped italic text with unbalanced double quote and more backslashes", func() {
						source := `\\\__italic content_`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `\\__italic content_`}, // only 1 backslash remove
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})
				})

				Context("with nested quoted text", func() {

					It("escaped italic text with nested monospace text", func() {
						source := `\` + "_`monospace content`_" // gives: \_`monospace content`_
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "_"},
										types.QuotedText{
											Kind: types.SingleQuoteMonospace,
											Elements: []interface{}{
												types.StringElement{Content: "monospace content"},
											},
										},
										types.StringElement{Content: "_"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped italic text with unbalanced double quote and nested bold test", func() {
						source := `\__*bold content*_`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "__"},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "bold content"},
											},
										},
										types.StringElement{Content: "_"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped italic text with nested bold text", func() {
						source := `\_italic *and bold* content_`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "_italic "},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "and bold"},
											},
										},
										types.StringElement{Content: " content_"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})
				})
			})

			Context("prevented monospace text substitution", func() {

				Context("without nested quoted text", func() {

					It("escaped monospace text with single quote", func() {
						source := `\` + "`monospace content`"
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "`monospace content`"}, // backslash removed
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped monospace text with single quote and more backslashes", func() {
						source := `\\` + "`monospace content`"
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `\` + "`monospace content`"}, // only 1 backslash removed
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped monospace text with double quote", func() {
						source := `\\` + "`monospace content``"
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `\` + "`monospace content``"}, // 2 back slashes "consumed"
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped monospace text with double quote and more backslashes", func() {
						source := `\\\` + "``monospace content``" // 3 backslashes
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `\` + "``monospace content``"}, // 2 back slashes "consumed"
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped monospace text with unbalanced double quote", func() {
						source := `\` + "``monospace content`"
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "``monospace content`"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped monospace text with unbalanced double quote and more backslashes", func() {
						source := `\\\` + "``monospace content`" // 3 backslashes
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `\\` + "``monospace content`"}, // 2 backslashes removed
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})
				})

				Context("with nested quoted text", func() {

					It("escaped monospace text with nested bold text", func() {
						source := `\` + "`*bold content*`" // gives: \`*bold content*`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "`"},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "bold content"},
											},
										},
										types.StringElement{Content: "`"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped monospace text with unbalanced double backquote and nested bold test", func() {
						source := `\` + "``*bold content*`"
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "``"},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "bold content"},
											},
										},
										types.StringElement{Content: "`"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped monospace text with nested bold text", func() {
						source := `\` + "`monospace *and bold* content`"
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "`monospace "},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "and bold"},
											},
										},
										types.StringElement{Content: " content`"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})
				})
			})

			Context("prevented subscript text substitution", func() {

				Context("without nested quoted text", func() {

					It("escaped subscript text with single quote", func() {
						source := `\~subscriptcontent~`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "~subscriptcontent~"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped subscript text with single quote and more backslashes", func() {
						source := `\\~subscriptcontent~`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `\~subscriptcontent~`}, // only 1 backslash removed
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

				})

				Context("with nested quoted text", func() {

					It("escaped subscript text with nested bold text", func() {
						source := `\~*boldcontent*~`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "~"},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "boldcontent"},
											},
										},
										types.StringElement{Content: "~"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped subscript text with nested bold text", func() {
						source := `\~subscript *and bold* content~`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `\~subscript `},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "and bold"},
											},
										},
										types.StringElement{Content: " content~"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})
				})
			})

			Context("prevented superscript text substitution", func() {

				Context("without nested quoted text", func() {

					It("escaped superscript text with single quote", func() {
						source := `\^superscriptcontent^`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "^superscriptcontent^"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped superscript text with single quote and more backslashes", func() {
						source := `\\^superscriptcontent^`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `\^superscriptcontent^`}, // only 1 backslash removed
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

				})

				Context("with nested quoted text", func() {

					It("escaped superscript text with nested bold text - case 1", func() {
						source := `\^*bold content*^` // valid escaped superscript since it has no space within
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `^`},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "bold content"},
											},
										},
										types.StringElement{Content: "^"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped superscript text with unbalanced double backquote and nested bold test", func() {
						source := `\^*bold content*^`
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "^"},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "bold content"},
											},
										},
										types.StringElement{Content: "^"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})

					It("escaped superscript text with nested bold text - case 2", func() {
						source := `\^superscript *and bold* content^` // invalid superscript text since it has spaces within
						expected := types.DocumentFragments{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: `\^superscript `},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "and bold"},
											},
										},
										types.StringElement{Content: " content^"},
									},
								},
							},
						}
						Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
					})
				})
			})
		})

		Context("with quoted string", func() {

			It("apostrophes in single bold", func() {
				source := "this *mother's mothers' mothers`'*\n"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "this "},
								types.QuotedText{
									Kind: types.SingleQuoteBold,
									Elements: []interface{}{
										types.StringElement{
											Content: "mother\u2019s mothers' mothers\u2019",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("apostrophes in double bold", func() {
				source := "this **mother's mothers' mothers`'**\n"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "this "},
								types.QuotedText{
									Kind: types.DoubleQuoteBold,
									Elements: []interface{}{
										types.StringElement{
											Content: "mother\u2019s mothers' mothers\u2019",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("apostrophes in single italic", func() {
				source := "this _mother's mothers' mothers`'_\n"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "this "},
								types.QuotedText{
									Kind: types.SingleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{
											Content: "mother\u2019s mothers' mothers\u2019",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("apostrophes in double italic", func() {
				source := "this __mother's mothers' mothers`'__\n"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "this "},
								types.QuotedText{
									Kind: types.DoubleQuoteItalic,
									Elements: []interface{}{
										types.StringElement{
											Content: "mother\u2019s mothers' mothers\u2019",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("apostrophes in single mono", func() {
				source := "this `mother's mothers`' day`\n" // no typographic quotes here
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "this "},
								types.QuotedText{
									Kind: types.SingleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{
											Content: "mother\u2019s mothers\u2019 day",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("apostrophes in double mono", func() {
				source := "this ``mother's mothers`' day``\n"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "this "},
								types.QuotedText{
									Kind: types.DoubleQuoteMonospace,
									Elements: []interface{}{
										types.StringElement{
											Content: "mother\u2019s mothers\u2019 day",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("apostrophes in single marked", func() {
				source := "this #mother's mothers' mothers`'#\n"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "this "},
								types.QuotedText{
									Kind: types.SingleQuoteMarked,
									Elements: []interface{}{
										types.StringElement{
											Content: "mother\u2019s mothers' mothers\u2019",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("apostrophes in double marked", func() {
				source := "this ##mother's mothers' mothers`'##\n"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "this "},
								types.QuotedText{
									Kind: types.DoubleQuoteMarked,
									Elements: []interface{}{
										types.StringElement{
											Content: "mother\u2019s mothers' mothers\u2019",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})
		})
	})

	Context("in final documents", func() {

		Context("with single punctuation", func() {

			It("bold text with 1 word", func() {
				source := "*hello*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "hello"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("bold text with 2 words", func() {
				source := "*bold    content*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "bold    content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("bold text with 3 words", func() {
				source := "*some bold content*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "some bold content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("italic text with 3 words in single quote", func() {
				source := "_some italic content_"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "some italic content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("monospace text with 3 words", func() {
				source := "`some monospace content`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "some monospace content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("invalid subscript text with 3 words", func() {
				source := "~some subscript content~"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "~some subscript content~"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("invalid superscript text with 3 words", func() {
				source := "^some superscript content^"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "^some superscript content^"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("bold text within italic text", func() {
				source := "_some *bold* content_"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "some "},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "bold"},
												},
											},
											types.StringElement{Content: " content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("monospace text within bold text within italic quote", func() {
				source := "*some _italic and `monospaced content`_*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "some "},
											types.QuotedText{
												Kind: types.SingleQuoteItalic,
												Elements: []interface{}{
													types.StringElement{Content: "italic and "},
													types.QuotedText{
														Kind: types.SingleQuoteMonospace,
														Elements: []interface{}{
															types.StringElement{Content: "monospaced content"},
														},
													},
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

			It("italic text within italic text", func() {
				source := "_some _very italic_ content_"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "some _very italic"},
										},
									},
									types.StringElement{Content: " content_"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("subscript text attached", func() {
				source := "O~2~ is a molecule"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "O"},
									types.QuotedText{
										Kind: types.SingleQuoteSubscript,
										Elements: []interface{}{
											types.StringElement{Content: "2"},
										},
									},
									types.StringElement{Content: " is a molecule"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("superscript text attached", func() {
				source := "M^me^ White"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "M"},
									types.QuotedText{
										Kind: types.SingleQuoteSuperscript,
										Elements: []interface{}{
											types.StringElement{Content: "me"},
										},
									},
									types.StringElement{Content: " White"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("invalid subscript text with 3 words", func() {
				source := "~some subscript content~"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "~some subscript content~"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("bold text across paragraph", func() {
				source := "*some bold\n\ncontent*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "*some bold"},
								},
							},
						},
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "content*"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("italic text across paragraph", func() {
				source := "_some italic\n\ncontent_"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "_some italic"},
								},
							},
						},
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "content_"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("monospace text across paragraph", func() {
				source := "`some monospace\n\ncontent`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "`some monospace"},
								},
							},
						},
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "content`"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("marked text across paragraph", func() {
				source := "#some marked\n\ncontent#"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "#some marked",
									},
								},
							},
						},
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "content#",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

		})

		Context("with double punctuation", func() {

			It("bold text of 1 word in double quote", func() {
				source := "**hello**"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.DoubleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "hello"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("italic text with 3 words in double quote", func() {
				source := "__some italic content__"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.DoubleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "some italic content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("monospace text with 3 words in double quote", func() {
				source := "`` some monospace content ``"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.DoubleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: " some monospace content "},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("superscript text within italic text", func() {
				source := "__some ^superscript^ content__"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.DoubleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "some "},
											types.QuotedText{
												Kind: types.SingleQuoteSuperscript,
												Elements: []interface{}{
													types.StringElement{Content: "superscript"},
												},
											},
											types.StringElement{Content: " content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("superscript text within italic text within bold quote", func() {
				source := "**some _italic and ^superscriptcontent^_**"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.DoubleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "some "},
											types.QuotedText{
												Kind: types.SingleQuoteItalic,
												Elements: []interface{}{
													types.StringElement{Content: "italic and "},
													types.QuotedText{
														Kind: types.SingleQuoteSuperscript,
														Elements: []interface{}{
															types.StringElement{Content: "superscriptcontent"},
														},
													},
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

			It("bold text across paragraph", func() {
				source := "**some bold\n\ncontent**"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "**some bold",
									},
								},
							},
						},
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "content**"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("italic text across paragraph", func() {
				source := "__some italic\n\ncontent__"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "__some italic",
									},
								},
							},
						},
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "content__"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("monospace text across paragraph", func() {
				source := "``some monospace\n\ncontent``"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "``some monospace",
									},
								},
							},
						},
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "content``"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("double marked text across paragraph", func() {
				source := "##some marked\n\ncontent##"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "##some marked",
									},
								},
							},
						},
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "content##"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

		})

		Context("quoted text inline", func() {

			It("inline content with bold text", func() {
				source := "a paragraph with *some bold content*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a paragraph with "},
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "some bold content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("inline content with invalid bold text - use case 1", func() {
				source := "a paragraph with *some bold content"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a paragraph with *some bold content"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("inline content with invalid bold text - use case 2", func() {
				source := "a paragraph with *some bold content *"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a paragraph with *some bold content *"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("inline content with invalid bold text - use case 3", func() {
				source := "a paragraph with * some bold content*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a paragraph with * some bold content*"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("invalid italic text within bold text", func() {
				source := "some *bold and _italic content _ together*."
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "some "},
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "bold and _italic content _ together"},
										},
									},
									types.StringElement{Content: "."},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("italic text within invalid bold text", func() {
				source := "some *bold and _italic content_ together *."
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "some *bold and "},
									types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "italic content"},
										},
									},
									types.StringElement{Content: " together *."},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("inline content with invalid subscript text - use case 1", func() {
				source := "a paragraph with ~some subscript content"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a paragraph with ~some subscript content"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("inline content with invalid subscript text - use case 2", func() {
				source := "a paragraph with ~some subscript content ~"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a paragraph with ~some subscript content ~"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("inline content with invalid subscript text - use case 3", func() {
				source := "a paragraph with ~ some subscript content~"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a paragraph with ~ some subscript content~"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("inline content with invalid superscript text - use case 1", func() {
				source := "a paragraph with ^some superscript content"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a paragraph with ^some superscript content"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("inline content with invalid superscript text - use case 2", func() {
				source := "a paragraph with ^some superscript content ^"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a paragraph with ^some superscript content ^"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("inline content with invalid superscript text - use case 3", func() {
				source := "a paragraph with ^ some superscript content^"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a paragraph with ^ some superscript content^"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})

		Context("with nested quoted text", func() {

			It("italic text within bold text", func() {
				source := "some *bold and _italic content_ together*."
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "some "},
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "bold and "},
											types.QuotedText{
												Kind: types.SingleQuoteItalic,
												Elements: []interface{}{
													types.StringElement{Content: "italic content"},
												},
											},
											types.StringElement{Content: " together"},
										},
									},
									types.StringElement{Content: "."},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("single-quote bold within single-quote bold text", func() {
				// here we don't allow for bold text within bold text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "*some *nested bold* content*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "some *nested bold"},
										},
									},
									types.StringElement{Content: " content*"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("double-quote bold within double-quote bold text", func() {
				// here we don't allow for bold text within bold text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "**some **nested bold** content**"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.DoubleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "some "},
										},
									},
									types.StringElement{Content: "nested bold"},
									types.QuotedText{
										Kind: types.DoubleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: " content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("single-quote bold within double-quote bold text", func() {
				// here we don't allow for bold text within bold text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "**some *nested bold* content**"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.DoubleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "some "},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "nested bold"},
												},
											},
											types.StringElement{Content: " content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("double-quote bold within single-quote bold text", func() {
				// here we don't allow for bold text within bold text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "*some **nested bold** content*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "some "},
											types.QuotedText{
												Kind: types.DoubleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "nested bold"},
												},
											},
											types.StringElement{Content: " content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("single-quote italic within single-quote italic text", func() {
				// here we don't allow for italic text within italic text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "_some _nested italic_ content_"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "some _nested italic"},
										},
									},
									types.StringElement{Content: " content_"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("double-quote italic within double-quote italic text", func() {
				// here we don't allow for italic text within italic text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "__some __nested italic__ content__"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.DoubleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "some "},
										},
									},
									types.StringElement{Content: "nested italic"},
									types.QuotedText{
										Kind: types.DoubleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: " content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("single-quote italic within double-quote italic text", func() {
				// here we allow for italic text within italic text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "__some _nested italic_ content__"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.DoubleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "some "},
											types.QuotedText{
												Kind: types.SingleQuoteItalic,
												Elements: []interface{}{
													types.StringElement{Content: "nested italic"},
												},
											},
											types.StringElement{Content: " content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("double-quote italic within single-quote italic text", func() {
				// here we allow for italic text within italic text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "_some __nested italic__ content_"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "some "},
											types.QuotedText{
												Kind: types.DoubleQuoteItalic,
												Elements: []interface{}{
													types.StringElement{Content: "nested italic"},
												},
											},
											types.StringElement{Content: " content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("single-quote monospace within single-quote monospace text", func() {
				// here we don't allow for monospace text within monospace text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "`some `nested monospace` content`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "some `nested monospace"},
										},
									},
									types.StringElement{Content: " content`"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("double-quote monospace within double-quote monospace text", func() {
				// here we don't allow for monospace text within monospace text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "``some ``nested monospace`` content``"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.DoubleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "some "},
										},
									},
									types.StringElement{Content: "nested monospace"},
									types.QuotedText{
										Kind: types.DoubleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: " content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("single-quote monospace within double-quote monospace text", func() {
				// here we allow for monospace text within monospace text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "`some ``nested monospace`` content`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "some "},
											types.QuotedText{
												Kind: types.DoubleQuoteMonospace,
												Elements: []interface{}{
													types.StringElement{Content: "nested monospace"},
												},
											},
											types.StringElement{Content: " content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("double-quote monospace within single-quote monospace text", func() {
				// here we allow for monospace text within monospace text, to comply with the existing implementations (asciidoc and asciidoctor)
				source := "`some ``nested monospace`` content`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "some "},
											types.QuotedText{
												Kind: types.DoubleQuoteMonospace,
												Elements: []interface{}{
													types.StringElement{Content: "nested monospace"},
												},
											},
											types.StringElement{Content: " content"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("unbalanced bold in monospace - case 1", func() {
				source := "`*a`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "*a"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("unbalanced bold in monospace - case 2", func() {
				source := "`a*b`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "a*b"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("italic in monospace", func() {
				source := "`_a_`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.QuotedText{
												Kind: types.SingleQuoteItalic,
												Elements: []interface{}{
													types.StringElement{Content: "a"},
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

			It("unbalanced italic in monospace", func() {
				source := "`a_b`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "a_b"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("unparsed bold in monospace", func() {
				source := "`a*b*`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "a*b*"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("parsed subscript in monospace", func() {
				source := "`a~b~`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "a"},
											types.QuotedText{
												Kind: types.SingleQuoteSubscript,
												Elements: []interface{}{
													types.StringElement{Content: "b"},
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

			It("multiline in single quoted monospace - case 1", func() {
				source := "`a\nb`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "a\nb"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("multiline in double quoted monospace - case 1", func() {
				source := "`a\nb`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "a\nb"},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("multiline in single quoted  monospace - case 2", func() {
				source := "`a\n*b*`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "a\n"},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "b"},
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

			It("multiline in double quoted  monospace - case 2", func() {
				source := "`a\n*b*`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "a\n"},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "b"},
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

			It("link in bold", func() {
				source := "*a link:/[b]*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlineLink{
												Attributes: types.Attributes{
													types.AttrInlineLinkText: "b",
												},
												Location: types.Location{
													Path: []interface{}{
														types.StringElement{
															Content: "/",
														},
													},
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

			It("image in bold", func() {
				source := "*a image:foo.png[]*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlineImage{
												Location: types.Location{
													Path: []interface{}{
														types.StringElement{
															Content: "foo.png",
														},
													},
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

			It("singleplus passthrough in bold", func() {
				source := "*a +image:foo.png[]+*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlinePassthrough{
												Kind: types.SinglePlusPassthrough,
												Elements: []interface{}{
													types.StringElement{Content: "image:foo.png[]"},
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

			It("tripleplus passthrough in bold", func() {
				source := "*a +++image:foo.png[]+++*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlinePassthrough{
												Kind: types.TriplePlusPassthrough,
												Elements: []interface{}{
													types.StringElement{Content: "image:foo.png[]"},
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

			It("link in italic", func() {
				source := "_a link:/[b]_"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlineLink{
												Attributes: types.Attributes{
													types.AttrInlineLinkText: "b",
												},
												Location: types.Location{
													Path: []interface{}{
														types.StringElement{
															Content: "/",
														},
													},
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

			It("image in italic", func() {
				source := "_a image:foo.png[]_"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlineImage{
												Location: types.Location{
													Path: []interface{}{
														types.StringElement{
															Content: "foo.png",
														},
													},
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

			It("singleplus passthrough in italic", func() {
				source := "_a +image:foo.png[]+_"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlinePassthrough{
												Kind: types.SinglePlusPassthrough,
												Elements: []interface{}{
													types.StringElement{Content: "image:foo.png[]"},
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

			It("tripleplus passthrough in italic", func() {
				source := "_a +++image:foo.png[]+++_"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlinePassthrough{
												Kind: types.TriplePlusPassthrough,
												Elements: []interface{}{
													types.StringElement{Content: "image:foo.png[]"},
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

			It("link in monospace", func() {
				source := "`a link:/[b]`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlineLink{
												Attributes: types.Attributes{
													types.AttrInlineLinkText: "b",
												},
												Location: types.Location{
													Path: []interface{}{
														types.StringElement{
															Content: "/",
														},
													},
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

			It("image in monospace", func() {
				source := "`a image:foo.png[]`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlineImage{
												Location: types.Location{
													Path: []interface{}{
														types.StringElement{
															Content: "foo.png",
														},
													},
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

			It("singleplus passthrough in monospace", func() {
				source := "`a +image:foo.png[]+`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlinePassthrough{
												Kind: types.SinglePlusPassthrough,
												Elements: []interface{}{
													types.StringElement{Content: "image:foo.png[]"},
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

			It("tripleplus passthrough in monospace", func() {
				source := "`a +++image:foo.png[]+++`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlinePassthrough{
												Kind: types.TriplePlusPassthrough,
												Elements: []interface{}{
													types.StringElement{Content: "image:foo.png[]"},
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

		Context("unbalanced quoted text", func() {

			Context("unbalanced bold text", func() {

				It("unbalanced bold text - extra on left", func() {
					source := "**some bold content*"
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "*some bold content"},
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("unbalanced bold text - extra on right", func() {
					source := "*some bold content**"
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{

										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{Content: "some bold content"},
											},
										},
										types.StringElement{Content: "*"},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})
			})

			Context("unbalanced italic text", func() {

				It("unbalanced italic text - extra on left", func() {
					source := "__some italic content_"
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{

										types.QuotedText{
											Kind: types.SingleQuoteItalic,
											Elements: []interface{}{
												types.StringElement{Content: "_some italic content"},
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("unbalanced italic text - extra on right", func() {
					source := "_some italic content__"
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.QuotedText{
											Kind: types.SingleQuoteItalic,
											Elements: []interface{}{
												types.StringElement{Content: "some italic content"},
											},
										},
										types.StringElement{Content: "_"},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})
			})

			Context("unbalanced monospace text", func() {

				It("unbalanced monospace text - extra on left", func() {
					source := "``some monospace content`"
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.QuotedText{
											Kind: types.SingleQuoteMonospace,
											Elements: []interface{}{
												types.StringElement{Content: "`some monospace content"},
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("unbalanced monospace text - extra on right", func() {
					source := "`some monospace content``"
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.QuotedText{
											Kind: types.SingleQuoteMonospace,
											Elements: []interface{}{
												types.StringElement{Content: "some monospace content"},
											},
										},
										types.StringElement{Content: "`"},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})
			})

			It("inline content with unbalanced bold text", func() {
				source := "a paragraph with *some bold content"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a paragraph with *some bold content"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

		})

		Context("prevented substitution", func() {

			Context("prevented bold text substitution", func() {

				Context("without nested quoted text", func() {

					It("escaped bold text with single backslash", func() {
						source := `\*bold content*`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "*bold content*"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped bold text with multiple backslashes", func() {
						source := `\\*bold content*`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `\*bold content*`},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped bold text with double quote", func() {
						source := `\\**bold content**`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `**bold content**`},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped bold text with double quote and more backslashes", func() {
						source := `\\\**bold content**`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `\**bold content**`},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped bold text with unbalanced double quote", func() {
						source := `\**bold content*`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `**bold content*`},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped bold text with unbalanced double quote and more backslashes", func() {
						source := `\\\**bold content*`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `\\**bold content*`},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})
				})

				Context("with nested quoted text", func() {

					It("escaped bold text with nested italic text", func() {
						source := `\*_italic content_*`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "*"},
											types.QuotedText{
												Kind: types.SingleQuoteItalic,
												Elements: []interface{}{
													types.StringElement{Content: "italic content"},
												},
											},
											types.StringElement{Content: "*"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped bold text with unbalanced double quote and nested italic test", func() {
						source := `\**_italic content_*`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "**"},
											types.QuotedText{
												Kind: types.SingleQuoteItalic,
												Elements: []interface{}{
													types.StringElement{Content: "italic content"},
												},
											},
											types.StringElement{Content: "*"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped bold text with nested italic", func() {
						source := `\*bold _and italic_ content*`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "*bold "},
											types.QuotedText{
												Kind: types.SingleQuoteItalic,
												Elements: []interface{}{
													types.StringElement{Content: "and italic"},
												},
											},
											types.StringElement{Content: " content*"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})
				})

			})

			Context("prevented italic text substitution", func() {

				Context("without nested quoted text", func() {

					It("escaped italic text with single quote", func() {
						source := `\_italic content_`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "_italic content_"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped italic text with single quote and more backslashes", func() {
						source := `\\_italic content_`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `\_italic content_`},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped italic text with double quote with 2 backslashes", func() {
						source := `\\__italic content__`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `__italic content__`},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped italic text with double quote with 3 backslashes", func() {
						source := `\\\__italic content__`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `\__italic content__`},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped italic text with unbalanced double quote", func() {
						source := `\__italic content_`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `__italic content_`},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped italic text with unbalanced double quote and more backslashes", func() {
						source := `\\\__italic content_`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `\\__italic content_`}, // only 1 backslash remove
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})
				})

				Context("with nested quoted text", func() {

					It("escaped italic text with nested monospace text", func() {
						source := `\` + "_`monospace content`_" // gives: \_`monospace content`_
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "_"},
											types.QuotedText{
												Kind: types.SingleQuoteMonospace,
												Elements: []interface{}{
													types.StringElement{Content: "monospace content"},
												},
											},
											types.StringElement{Content: "_"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped italic text with unbalanced double quote and nested bold test", func() {
						source := `\__*bold content*_`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "__"},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "bold content"},
												},
											},
											types.StringElement{Content: "_"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped italic text with nested bold text", func() {
						source := `\_italic *and bold* content_`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "_italic "},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "and bold"},
												},
											},
											types.StringElement{Content: " content_"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})
				})
			})

			Context("prevented monospace text substitution", func() {

				Context("without nested quoted text", func() {

					It("escaped monospace text with single quote", func() {
						source := `\` + "`monospace content`"
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "`monospace content`"}, // backslash removed
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped monospace text with single quote and more backslashes", func() {
						source := `\\` + "`monospace content`"
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `\` + "`monospace content`"}, // only 1 backslash removed
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped monospace text with double quote", func() {
						source := `\\` + "`monospace content``"
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `\` + "`monospace content``"}, // 2 back slashes "consumed"
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped monospace text with double quote and more backslashes", func() {
						source := `\\\` + "``monospace content``" // 3 backslashes
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `\` + "``monospace content``"}, // 2 back slashes "consumed"
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped monospace text with unbalanced double quote", func() {
						source := `\` + "``monospace content`"
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "``monospace content`"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped monospace text with unbalanced double quote and more backslashes", func() {
						source := `\\\` + "``monospace content`" // 3 backslashes
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `\\` + "``monospace content`"}, // 2 backslashes removed
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})
				})

				Context("with nested quoted text", func() {

					It("escaped monospace text with nested bold text", func() {
						source := `\` + "`*bold content*`" // gives: \`*bold content*`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "`"},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "bold content"},
												},
											},
											types.StringElement{Content: "`"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped monospace text with unbalanced double backquote and nested bold test", func() {
						source := `\` + "``*bold content*`"
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "``"},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "bold content"},
												},
											},
											types.StringElement{Content: "`"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped monospace text with nested bold text", func() {
						source := `\` + "`monospace *and bold* content`"
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "`monospace "},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "and bold"},
												},
											},
											types.StringElement{Content: " content`"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})
				})
			})

			Context("prevented subscript text substitution", func() {

				Context("without nested quoted text", func() {

					It("escaped subscript text with single quote", func() {
						source := `\~subscriptcontent~`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "~subscriptcontent~"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped subscript text with single quote and more backslashes", func() {
						source := `\\~subscriptcontent~`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `\~subscriptcontent~`}, // only 1 backslash removed
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

				})

				Context("with nested quoted text", func() {

					It("escaped subscript text with nested bold text", func() {
						source := `\~*boldcontent*~`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "~"},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "boldcontent"},
												},
											},
											types.StringElement{Content: "~"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped subscript text with nested bold text", func() {
						source := `\~subscript *and bold* content~`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `\~subscript `},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "and bold"},
												},
											},
											types.StringElement{Content: " content~"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})
				})
			})

			Context("prevented superscript text substitution", func() {

				Context("without nested quoted text", func() {

					It("escaped superscript text with single quote", func() {
						source := `\^superscriptcontent^`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "^superscriptcontent^"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped superscript text with single quote and more backslashes", func() {
						source := `\\^superscriptcontent^`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `\^superscriptcontent^`}, // only 1 backslash removed
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

				})

				Context("with nested quoted text", func() {

					It("escaped superscript text with nested bold text - case 1", func() {
						source := `\^*bold content*^` // valid escaped superscript since it has no space within
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `^`},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "bold content"},
												},
											},
											types.StringElement{Content: "^"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped superscript text with unbalanced double backquote and nested bold test", func() {
						source := `\^*bold content*^`
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: "^"},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "bold content"},
												},
											},
											types.StringElement{Content: "^"},
										},
									},
								},
							},
						}
						Expect(ParseDocument(source)).To(MatchDocument(expected))
					})

					It("escaped superscript text with nested bold text - case 2", func() {
						source := `\^superscript *and bold* content^` // invalid superscript text since it has spaces within
						expected := types.Document{
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{Content: `\^superscript `},
											types.QuotedText{
												Kind: types.SingleQuoteBold,
												Elements: []interface{}{
													types.StringElement{Content: "and bold"},
												},
											},
											types.StringElement{Content: " content^"},
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

		Context("nested images", func() {

			It("image in bold", func() {
				source := "*a image:foo.png[]*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlineImage{
												Location: types.Location{
													Path: []interface{}{
														types.StringElement{
															Content: "foo.png",
														},
													},
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

			It("image in italic", func() {
				source := "_a image:foo.png[]_"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlineImage{
												Location: types.Location{
													Path: []interface{}{
														types.StringElement{
															Content: "foo.png",
														},
													},
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

			It("image in monospace", func() {
				source := "`a image:foo.png[]`"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{Content: "a "},
											types.InlineImage{
												Location: types.Location{
													Path: []interface{}{
														types.StringElement{
															Content: "foo.png",
														},
													},
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
