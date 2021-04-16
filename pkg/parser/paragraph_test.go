package parser_test

import (
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("paragraphs", func() {

	Context("in raw documents", func() {

		Context("regular paragraphs", func() {

			It("with basic content", func() {
				source := `cookie
chocolate
pasta`
				expected := types.DocumentFragments{
					types.DocumentFragment{
						LineOffset: 1,
						Content: []interface{}{
							types.RawLine("cookie"),
							types.RawLine("chocolate"),
							types.RawLine("pasta"),
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
			})

			It("with hardbreaks attribute", func() {
				source := `[%hardbreaks]
cookie
chocolate
pasta`
				expected := types.DocumentFragments{
					types.DocumentFragment{
						LineOffset: 1,
						Content: []interface{}{
							types.Attributes{
								types.AttrOptions: []interface{}{"hardbreaks"},
							},
							types.RawLine("cookie"),
							types.RawLine("chocolate"),
							types.RawLine("pasta"),
						},
					},
				}
				result, err := ParseDocumentFragments(source)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(MatchDocumentFragments(expected))
			})

			It("with title attribute", func() {
				source := `[title=my title]
cookie
pasta`
				expected := types.DocumentFragments{
					types.DocumentFragment{
						LineOffset: 1,
						Content: []interface{}{
							types.Attributes{
								types.AttrTitle: "my title",
							},
							types.RawLine("cookie"),
							types.RawLine("pasta"),
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
			})

			It("with custom title attribute - explicit and unquoted", func() {
				source := `:title: cookies
				
[title=my {title}]
cookie
pasta`
				expected := types.DocumentFragments{
					types.DocumentFragment{
						LineOffset: 1,
						Content: []interface{}{
							types.AttributeDeclaration{
								Name:  "title",
								Value: "cookies",
							},
						},
					},
					types.DocumentFragment{
						LineOffset: 3,
						Content: []interface{}{
							types.Attributes{
								types.AttrTitle: []interface{}{
									types.StringElement{
										Content: "my ",
									},
									types.AttributeSubstitution{
										Name: "title",
									},
								},
							},
							types.RawLine("cookie"),
							types.RawLine("pasta"),
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
			})

			It("with multiple attributes and blanklines in-between", func() {
				source := `[%hardbreaks.role1.role2]

[#anchor]

cookie
pasta`
				expected := types.DocumentFragments{
					types.DocumentFragment{
						LineOffset: 1,
						Content: []interface{}{
							types.Attributes{
								types.AttrRoles:   []interface{}{"role1", "role2"},
								types.AttrOptions: []interface{}{"hardbreaks"},
							},
						},
					},
					types.DocumentFragment{
						LineOffset: 3,
						Content: []interface{}{
							types.Attributes{
								types.AttrID: "anchor",
							},
						},
					},
					types.DocumentFragment{
						LineOffset: 5,
						Content: []interface{}{
							types.RawLine("cookie"),
							types.RawLine("pasta"),
						},
					},
				}
				Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
			})

			Context("with custom substitutions", func() {

				// using the same input for all substitution tests
				source := `:github-url: https://github.com

[subs="$SUBS"]
a link to https://github.com[] <using the *inline link macro*>
another one using attribute substitution: {github-url}[]...
// a single-line comment.`

				It("should read multiple lines", func() {
					s := strings.ReplaceAll(source, "$SUBS", "normal")
					expected := types.DocumentFragments{
						types.DocumentFragment{
							LineOffset: 1,
							Content: []interface{}{
								types.AttributeDeclaration{
									Name:  "github-url",
									Value: "https://github.com",
								},
							},
						},
						types.DocumentFragment{
							LineOffset: 3,
							Content: []interface{}{
								types.Attributes{
									types.AttrSubstitutions: "normal",
								},
								types.RawLine("a link to https://github.com[] <using the *inline link macro*>"),
								types.RawLine("another one using attribute substitution: {github-url}[]..."),
								types.SingleLineComment{
									Content: " a single-line comment.",
								},
							},
						},
					}
					result, err := ParseDocumentFragments(s)
					Expect(err).NotTo(HaveOccurred())
					Expect(result).To(MatchDocumentFragments(expected))
				})

			})
		})
	})

	Context("in final documents", func() {

		Context("regular paragraphs", func() {

			It("with title attribute", func() {
				source := `[title=my title]
cookie
pasta`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrTitle: "my title",
							},
							Elements: []interface{}{
								types.StringElement{Content: "cookie\npasta"},
							},
						},
					},
				}
				result, err := ParseDocument(source)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(MatchDocument(expected))
			})

			It("with custom title attribute - explicit and unquoted", func() {
				source := `:title: cookies
				
[title=my {title}]
cookie
pasta`
				expected := types.Document{
					Elements: []interface{}{
						types.AttributeDeclaration{
							Name:  "title",
							Value: "cookies",
						},
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrTitle: "my cookies",
							},
							Elements: []interface{}{
								types.StringElement{Content: "cookie\npasta"},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with custom title attribute - explicit and single quoted", func() {
				source := `:title: cookies
				
[title='my {title}']
cookie
pasta`
				expected := types.Document{
					Elements: []interface{}{
						types.AttributeDeclaration{
							Name:  "title",
							Value: "cookies",
						},
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrTitle: "my cookies",
							},
							Elements: []interface{}{
								types.StringElement{Content: "cookie\npasta"},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with custom title attribute - explicit and double quoted", func() {
				source := `:title: cookies
				
[title="my {title}"]
cookie
pasta`
				expected := types.Document{
					Elements: []interface{}{
						types.AttributeDeclaration{
							Name:  "title",
							Value: "cookies",
						},
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrTitle: "my cookies",
							},
							Elements: []interface{}{
								types.StringElement{Content: "cookie\npasta"},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with custom title attribute - implicit", func() {
				source := `:title: cookies
				
.my {title}
cookie
pasta`
				expected := types.Document{
					Elements: []interface{}{
						types.AttributeDeclaration{
							Name:  "title",
							Value: "cookies",
						},
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrTitle: "my cookies",
							},
							Elements: []interface{}{
								types.StringElement{Content: "cookie\npasta"},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with multiple attributes without blanklines in-between", func() {
				source := `[%hardbreaks.role1.role2]
[#anchor]
cookie
pasta`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrID:      "anchor",
								types.AttrRoles:   []interface{}{"role1", "role2"},
								types.AttrOptions: []interface{}{"hardbreaks"},
							},
							Elements: []interface{}{
								types.StringElement{Content: "cookie\npasta"},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with multiple attributes and blanklines in-between", func() {
				source := `[%hardbreaks.role1.role2]

[#anchor]

cookie
pasta`
				expected := types.Document{
					Elements: []interface{}{
						types.BlankLine{}, // attribute not retained in blankline
						types.BlankLine{}, // attribute not retained in blankline
						types.Paragraph{
							Elements: []interface{}{
								types.StringElement{Content: "cookie\npasta"},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with paragraph roles and attribute - case 1", func() {
				source := `[.role1%hardbreaks.role2]
cookie
pasta`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrOptions: []interface{}{"hardbreaks"},
								types.AttrRoles:   []interface{}{"role1", "role2"},
							},
							Elements: []interface{}{
								types.StringElement{Content: "cookie\npasta"},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with paragraph roles - case 2", func() {
				source := `[.role1%hardbreaks]
[.role2]
cookie
pasta`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrOptions: []interface{}{"hardbreaks"},
								types.AttrRoles:   []interface{}{"role1", "role2"},
							},
							Elements: []interface{}{
								types.StringElement{Content: "cookie\npasta"},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("not treat plusplus as line break", func() {
				source := `C++
cookie`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Elements: []interface{}{
								types.StringElement{Content: "C++\ncookie"},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			Context("with counters", func() {

				It("default", func() {
					source := `cookie{counter:cookie} chocolate{counter2:cookie} pasta{counter:cookie} bob{counter:bob}`
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Elements: []interface{}{
									types.StringElement{Content: "foo1 chocolate baz3 bob1"},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("with numeric start", func() {
					source := `cookie{counter:cookie:2} chocolate{counter2:cookie} pasta{counter:cookie} bob{counter:bob:10}`
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Elements: []interface{}{
									types.StringElement{Content: "foo2 chocolate baz4 bob10"},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("with alphanumeric start", func() {
					source := `cookie{counter:cookie:b} chocolate{counter2:cookie} pasta{counter:cookie} bob{counter:bob:z}`
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrIDPrefix: "bar_",
						},
						Elements: []interface{}{
							types.Paragraph{
								Elements: []interface{}{
									types.StringElement{Content: "foob chocolate bazd bobz"},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})
			})

			It("paragraph with custom id prefix and title", func() {
				source := `:idprefix: bar_
			
.a title
a paragraph`
				expected := types.Document{
					Attributes: types.Attributes{
						types.AttrIDPrefix: "bar_",
					},
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrTitle: "a title", // there is no default ID. Only custom IDs
							},
							Elements: []interface{}{
								types.StringElement{Content: "a paragraph"},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("empty paragraph", func() {
				source := `{blank}`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Elements: []interface{}{
								types.PredefinedAttribute{
									Name: "blank",
								},
							},
						},
					},
				}
				result, err := ParseDocument(source)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(MatchDocument(expected))
			})

			It("paragraph with predefined attribute", func() {
				source := "hello {plus} world"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Elements: []interface{}{
								types.StringElement{Content: "hello "},
								types.PredefinedAttribute{Name: "plus"},
								types.StringElement{Content: " world"},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			Context("with custom substitutions", func() {

				// using the same input for all substitution tests
				source := `:github-url: https://github.com
					
[subs="$SUBS"]
a link to https://github.com[] <using the *inline link macro*>
another one using attribute substitution: {github-url}[]...
// a single-line comment`

				It("should apply the 'none' substitution", func() {
					s := strings.ReplaceAll(source, "$SUBS", "none")
					expected := types.Document{
						Attributes: types.Attributes{
							"github-url": "https://github.com",
						},
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrSubstitutions: "none",
								},
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to https://github.com[] <using the *inline link macro*>"},
									},
									{
										types.StringElement{Content: "another one using attribute substitution: {github-url}[]..."},
									},
								},
							},
						},
					}
					Expect(ParseDocument(s)).To(MatchDocument(expected))
				})

				It("should apply the 'quotes' substitution on multiple lines", func() {
					// quoted text is parsed but inline link macro is not
					s := strings.ReplaceAll(source, "$SUBS", "quotes")
					expected := types.Document{
						Attributes: types.Attributes{
							"github-url": "https://github.com",
						},
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrSubstitutions: "quotes",
								},
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "a link to https://github.com[] <using the ",
										},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{
													Content: "inline link macro",
												},
											},
										},
										types.StringElement{
											Content: ">",
										},
									},
									{
										types.StringElement{
											Content: "another one using attribute substitution: {github-url}[]...",
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(s)).To(MatchDocument(expected))
				})

				It("should apply the 'macros' substitution on multiple lines", func() {
					// quoted text is not parsed but inline link macro is
					s := strings.ReplaceAll(source, "$SUBS", "macros")
					expected := types.Document{
						Attributes: types.Attributes{
							"github-url": "https://github.com",
						},
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrSubstitutions: "macros",
								},
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "a link to ",
										},
										types.InlineLink{
											Location: types.Location{
												Scheme: "https://",
												Path: []interface{}{
													types.StringElement{
														Content: "github.com",
													},
												},
											},
										},
										types.StringElement{
											Content: " <using the *inline link macro*>",
										},
									},
									{
										types.StringElement{
											Content: "another one using attribute substitution: {github-url}[]...",
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(s)).To(MatchDocument(expected))
				})

				It("should apply the 'attributes' substitution on multiple lines", func() {
					// quoted text is not parsed but inline link macro is
					s := strings.ReplaceAll(source, "$SUBS", "attributes")
					expected := types.Document{
						Attributes: types.Attributes{
							"github-url": "https://github.com",
						},
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrSubstitutions: "attributes",
								},
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to https://github.com[] <using the *inline link macro*>"},
									},
									{
										types.StringElement{Content: "another one using attribute substitution: https://github.com[]..."},
									},
								},
							},
						},
					}
					Expect(ParseDocument(s)).To(MatchDocument(expected))
				})

				It("should apply the 'attributes,macros' substitution on multiple lines", func() {
					// quoted text is not parsed but inline link macro is
					s := strings.ReplaceAll(source, "$SUBS", "attributes,macros")
					expected := types.Document{
						Attributes: types.Attributes{
							"github-url": "https://github.com",
						},
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrSubstitutions: "attributes,macros",
								},
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to "},
										types.InlineLink{
											Location: types.Location{
												Scheme: "https://",
												Path: []interface{}{
													types.StringElement{
														Content: "github.com",
													},
												},
											},
										},
										types.StringElement{Content: " <using the *inline link macro*>"},
									},
									{
										types.StringElement{Content: "another one using attribute substitution: "},
										types.InlineLink{
											Location: types.Location{
												Scheme: "https://",
												Path: []interface{}{
													types.StringElement{
														Content: "github.com",
													},
												},
											},
										},
										types.StringElement{Content: "..."},
									},
								},
							},
						},
					}
					Expect(ParseDocument(s)).To(MatchDocument(expected))
				})

				It("should apply the 'specialchars' substitution on multiple lines", func() {
					// quoted text is not parsed but inline link macro is
					s := strings.ReplaceAll(source, "$SUBS", "specialchars")
					expected := types.Document{
						Attributes: types.Attributes{
							"github-url": "https://github.com",
						},
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrSubstitutions: "specialchars",
								},
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to https://github.com[] "},
										types.SpecialCharacter{Name: "<"},
										types.StringElement{Content: "using the *inline link macro*"},
										types.SpecialCharacter{Name: ">"},
									},
									{
										types.StringElement{Content: "another one using attribute substitution: {github-url}[]..."},
									},
								},
							},
						},
					}
					Expect(ParseDocument(s)).To(MatchDocument(expected))
				})

				It("should apply the replacements substitution on multiple lines", func() {
					// quoted text is not parsed but inline link macro is
					s := strings.ReplaceAll(source, "$SUBS", "replacements")
					expected := types.Document{
						Attributes: types.Attributes{
							"github-url": "https://github.com",
						},
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrSubstitutions: "replacements",
								},
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to https://github.com[] <using the *inline link macro*>"},
									},
									{
										types.StringElement{Content: "another one using attribute substitution: {github-url}[]\u2026\u200b"},
									},
								},
							},
						},
					}
					Expect(ParseDocument(s)).To(MatchDocument(expected))
				})

				It("should apply the 'quotes' and 'macros' substitutions", func() {
					// quoted text and inline link macro are both parsed
					s := strings.ReplaceAll(source, "$SUBS", "quotes,macros")
					expected := types.Document{
						Attributes: types.Attributes{
							"github-url": "https://github.com",
						},
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrSubstitutions: "quotes,macros",
								},
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "a link to ",
										},
										types.InlineLink{
											Location: types.Location{
												Scheme: "https://",
												Path: []interface{}{
													types.StringElement{
														Content: "github.com",
													},
												},
											},
										},
										types.StringElement{
											Content: " <using the ",
										},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{
													Content: "inline link macro",
												},
											},
										},
										types.StringElement{
											Content: ">",
										},
									},
									{
										types.StringElement{Content: "another one using attribute substitution: {github-url}[]..."},
									},
								},
							},
						},
					}
					Expect(ParseDocument(s)).To(MatchDocument(expected))
				})

				It("should apply the 'macros' and 'quotes' substitutions", func() {
					// quoted text and inline link macro are both parsed
					// (same as above, but with subs in reversed order)
					s := strings.ReplaceAll(source, "$SUBS", "macros,quotes")
					expected := types.Document{
						Attributes: types.Attributes{
							"github-url": "https://github.com",
						},
						Elements: []interface{}{
							types.Paragraph{
								Attributes: types.Attributes{
									types.AttrSubstitutions: "macros,quotes",
								},
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "a link to ",
										},
										types.InlineLink{
											Location: types.Location{
												Scheme: "https://",
												Path: []interface{}{
													types.StringElement{
														Content: "github.com",
													},
												},
											},
										},
										types.StringElement{
											Content: " <using the ",
										},
										types.QuotedText{
											Kind: types.SingleQuoteBold,
											Elements: []interface{}{
												types.StringElement{
													Content: "inline link macro",
												},
											},
										},
										types.StringElement{
											Content: ">",
										},
									},
									{
										types.StringElement{Content: "another one using attribute substitution: {github-url}[]..."},
									},
								},
							},
						},
					}
					Expect(ParseDocument(s)).To(MatchDocument(expected))
				})
			})
		})

		Context("admonition paragraphs", func() {

			It("note admonition paragraph", func() {
				source := `NOTE: this is a note.`
				expected := types.Document{
					Elements: []interface{}{
						types.Attributes{
							types.AttrStyle: types.Note,
						},
						types.InlineElements{
							types.StringElement{Content: "this is a note."},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("warning admonition paragraph", func() {
				source := `WARNING: this is a multiline
warning!`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle: types.Warning,
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "this is a multiline"},
								},
								{
									types.StringElement{Content: "warning!"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("admonition note paragraph with id and title", func() {
				source := `[[cookie]]
.chocolate
NOTE: this is a note.`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle: types.Note,
								types.AttrID:    "cookie",
								types.AttrTitle: "chocolate",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "this is a note."},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("caution admonition paragraph with single line", func() {
				source := `[CAUTION]
this is a caution!`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle: types.Caution,
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "this is a caution!"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("multiline caution admonition paragraph with title and id", func() {
				source := `[[cookie]]
[CAUTION] 
.chocolate
this is a 
*caution*!`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle: types.Caution,
								types.AttrID:    "cookie",
								types.AttrTitle: "chocolate",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "this is a "},
								},
								{
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{
												Content: "caution",
											},
										},
									},
									types.StringElement{
										Content: "!",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("multiple admonition paragraphs", func() {
				source := `[NOTE]
No space after the [NOTE]!

[CAUTION]
And no space after [CAUTION] either.`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle: types.Note,
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "No space after the [NOTE]!"},
								},
							},
						},
						types.BlankLine{},
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle: types.Caution,
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "And no space after [CAUTION] either."},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})

		Context("quote paragraphs", func() {

			It("inline image within a quote", func() {
				source := `[quote, john doe, quote title]
a cookie image:cookie.png[]`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle:       types.Quote,
								types.AttrQuoteAuthor: "john doe",
								types.AttrQuoteTitle:  "quote title",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "a cookie ",
									},
									types.InlineImage{
										Location: types.Location{
											Path: []interface{}{
												types.StringElement{
													Content: "cookie.png",
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

		Context("verse paragraphs", func() {

			It("paragraph as a verse with author and title", func() {
				source := `[verse, john doe, verse title]
I am a verse paragraph.`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle:       types.Verse,
								types.AttrQuoteAuthor: "john doe",
								types.AttrQuoteTitle:  "verse title",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "I am a verse paragraph."},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("paragraph as a verse with author, title and other attributes", func() {
				source := `[[universal]]
[verse, john doe, verse title]
.universe
I am a verse paragraph.`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle:       types.Verse,
								types.AttrQuoteAuthor: "john doe",
								types.AttrQuoteTitle:  "verse title",
								types.AttrID:          "universal",
								// types.AttrCustomID:    true,
								types.AttrTitle: "universe",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "I am a verse paragraph."},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("paragraph as a verse with empty title", func() {
				source := `[verse, john doe, ]
I am a verse paragraph.`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle:       types.Verse,
								types.AttrQuoteAuthor: "john doe",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "I am a verse paragraph."},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("paragraph as a verse without title", func() {
				source := `[verse, john doe ]
I am a verse paragraph.`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle:       types.Verse,
								types.AttrQuoteAuthor: "john doe",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "I am a verse paragraph."},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("paragraph as a verse with empty author", func() {
				source := `[verse,  ]
I am a verse paragraph.`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle: types.Verse,
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "I am a verse paragraph."},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("paragraph as a verse without author", func() {
				source := `[verse]
I am a verse paragraph.`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle: types.Verse,
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "I am a verse paragraph."},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("image block as a verse", func() {
				// assume that the author meant to use an image, so the `verse` attribute will be ignored during rendering
				source := `[verse, john doe, verse title]
image::cookie.png[]`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrStyle:       types.Verse,
								types.AttrQuoteAuthor: "john doe",
								types.AttrQuoteTitle:  "verse title",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "image::cookie.png[]"},
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
