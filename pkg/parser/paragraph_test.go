package parser_test

import (
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("paragraphs", func() {

	Context("raw documents", func() {

		Context("regular paragraphs", func() {

			It("with basic content", func() {
				source := `foo
bar
baz`
				expected := types.DocumentFragments{
					types.InlineElements{
						types.StringElement{Content: "foo"},
					},
					types.InlineElements{
						types.StringElement{Content: "bar"},
					},
					types.InlineElements{
						types.StringElement{Content: "baz"},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with explicit line break", func() {
				source := `foo +
bar
baz`
				expected := types.DocumentFragments{
					types.InlineElements{
						types.StringElement{Content: "foo"},
						types.LineBreak{},
					},
					types.InlineElements{
						types.StringElement{Content: "bar"},
					},
					types.InlineElements{
						types.StringElement{Content: "baz"},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with hardbreaks attribute", func() {
				source := `[%hardbreaks]
foo
bar
baz`
				expected := types.DocumentFragments{
					types.Attributes{
						types.AttrOptions: []interface{}{"hardbreaks"},
					},
					types.InlineElements{
						types.StringElement{Content: "foo"},
					},
					types.InlineElements{
						types.StringElement{Content: "bar"},
					},
					types.InlineElements{
						types.StringElement{Content: "baz"},
					},
				}
				result, err := ParseRawSource(source)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(MatchDocumentFragments(expected))
			})

			It("with title attribute", func() {
				source := `[title=My Title]
foo
baz`
				expected := types.DocumentFragments{
					types.Attributes{
						types.AttrTitle: "My Title",
					},
					types.InlineElements{
						types.StringElement{Content: "foo"},
					},
					types.InlineElements{
						types.StringElement{Content: "baz"},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with custom title attribute - explicit and unquoted", func() {
				source := `:title: cookies
				
[title=my {title}]
foo
baz`
				expected := types.DocumentFragments{
					types.AttributeDeclaration{
						Name:  "title",
						Value: "cookies",
					},
					types.BlankLine{},
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
					types.InlineElements{
						types.StringElement{Content: "foo"},
					},
					types.InlineElements{
						types.StringElement{Content: "baz"},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with custom title attribute - explicit and single quoted", func() {
				source := `:title: cookies
				
[title='my {title}']
foo
baz`
				expected := types.DocumentFragments{
					types.AttributeDeclaration{
						Name:  "title",
						Value: "cookies",
					},
					types.BlankLine{},
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
					types.InlineElements{
						types.StringElement{Content: "foo"},
					},
					types.InlineElements{
						types.StringElement{Content: "baz"},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with custom title attribute - explicit and double quoted", func() {
				source := `:title: cookies
				
[title="my {title}"]
foo
baz`
				expected := types.DocumentFragments{
					types.AttributeDeclaration{
						Name:  "title",
						Value: "cookies",
					},
					types.BlankLine{},
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
					types.InlineElements{
						types.StringElement{Content: "foo"},
					},
					types.InlineElements{
						types.StringElement{Content: "baz"},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with custom title attribute - implicit", func() {
				source := `:title: cookies
				
.my {title}
foo
baz`
				expected := types.DocumentFragments{
					types.AttributeDeclaration{
						Name:  "title",
						Value: "cookies",
					},
					types.BlankLine{},
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
					types.InlineElements{
						types.StringElement{Content: "foo"},
					},
					types.InlineElements{
						types.StringElement{Content: "baz"},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with multiple attributes without blanklines in-between", func() {
				source := `[%hardbreaks.role1.role2]
[#anchor]
foo
baz`
				expected := types.DocumentFragments{
					types.Attributes{
						types.AttrID:      "anchor",
						types.AttrRoles:   []interface{}{"role1", "role2"},
						types.AttrOptions: []interface{}{"hardbreaks"},
					},
					types.InlineElements{
						types.StringElement{Content: "foo"},
					},
					types.InlineElements{
						types.StringElement{Content: "baz"},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with multiple attributes and blanklines in-between", func() {
				source := `[%hardbreaks.role1.role2]

[#anchor]

foo
baz`
				expected := types.DocumentFragments{
					types.Attributes{
						types.AttrRoles:   []interface{}{"role1", "role2"},
						types.AttrOptions: []interface{}{"hardbreaks"},
					},
					types.BlankLine{},
					types.Attributes{
						types.AttrID: "anchor",
					},
					types.BlankLine{},
					types.InlineElements{
						types.StringElement{Content: "foo"},
					},
					types.InlineElements{
						types.StringElement{Content: "baz"},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with paragraph roles and attribute - case 1", func() {
				source := `[.role1%hardbreaks.role2]
foo
bar
baz`
				expected := types.DocumentFragments{
					types.Attributes{
						types.AttrOptions: []interface{}{"hardbreaks"},
						types.AttrRoles:   []interface{}{"role1", "role2"},
					},
					types.InlineElements{
						types.StringElement{Content: "foo"},
					},
					types.InlineElements{
						types.StringElement{Content: "bar"},
					},
					types.InlineElements{
						types.StringElement{Content: "baz"},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("with paragraph roles - case 2", func() {
				source := `[.role1%hardbreaks]
[.role2]
foo
bar
baz`
				expected := types.DocumentFragments{
					types.Attributes{
						types.AttrOptions: []interface{}{"hardbreaks"},
						types.AttrRoles:   []interface{}{"role1", "role2"},
					},
					types.InlineElements{
						types.StringElement{Content: "foo"},
					},
					types.InlineElements{
						types.StringElement{Content: "bar"},
					},
					types.InlineElements{
						types.StringElement{Content: "baz"},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("not treat plusplus as line break", func() {
				source := `C++
foo`
				expected := types.DocumentFragments{
					types.InlineElements{
						types.StringElement{Content: "C++"},
					},
					types.InlineElements{
						types.StringElement{Content: "foo"},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			Context("with counters", func() {

				It("default", func() {
					source := `foo{counter:foo} bar{counter2:foo} baz{counter:foo} bob{counter:bob}`
					expected := types.DocumentFragments{
						types.InlineElements{
							types.StringElement{
								Content: "foo",
							},
							types.CounterSubstitution{
								Name: "foo",
							},
							types.StringElement{
								Content: " bar",
							},
							types.CounterSubstitution{
								Hidden: true,
								Name:   "foo",
							},
							types.StringElement{
								Content: " baz",
							},
							types.CounterSubstitution{
								Name: "foo",
							},
							types.StringElement{
								Content: " bob",
							},
							types.CounterSubstitution{
								Name: "bob",
							},
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("with numeric start", func() {
					source := `foo{counter:foo:2} bar{counter2:foo} baz{counter:foo} bob{counter:bob:10}`
					expected := types.DocumentFragments{
						types.InlineElements{
							types.StringElement{
								Content: "foo",
							},
							types.CounterSubstitution{
								Name:  "foo",
								Value: 2,
							},
							types.StringElement{
								Content: " bar",
							},
							types.CounterSubstitution{
								Hidden: true,
								Name:   "foo",
							},
							types.StringElement{
								Content: " baz",
							},
							types.CounterSubstitution{
								Name: "foo",
							},
							types.StringElement{
								Content: " bob",
							},
							types.CounterSubstitution{
								Name:  "bob",
								Value: 10,
							},
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})

				It("with alphanumeric start", func() {
					source := `foo{counter:foo:b} bar{counter2:foo} baz{counter:foo} bob{counter:bob:z}`
					expected := types.DocumentFragments{
						types.InlineElements{
							types.StringElement{
								Content: "foo",
							},
							types.CounterSubstitution{
								Name:  "foo",
								Value: int32(98), // asciicode for `b`
							},
							types.StringElement{
								Content: " bar",
							},
							types.CounterSubstitution{
								Hidden: true,
								Name:   "foo",
							},
							types.StringElement{
								Content: " baz",
							},
							types.CounterSubstitution{
								Name: "foo",
							},
							types.StringElement{
								Content: " bob",
							},
							types.CounterSubstitution{
								Name:  "bob",
								Value: int32(122), // asciicode for `z`
							},
						},
					}
					Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
				})
			})

			Context("with custom substitutions", func() {

				// using the same input for all substitution tests
				source := `:github-url: https://github.com

[subs="$SUBS"]
a link to https://github.com[] <using the *inline link macro*>
another one using attribute substitution: {github-url}[]...
// a single-line comment.`

				It("should apply the 'default' substitution on multiple lines", func() {
					// quoted text is parsed but inline link macro is not
					s := strings.ReplaceAll(source, "[subs=\"$SUBS\"]\n", "")
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						types.BlankLine{},
						// [subs="$SUBS"] was removed in this test case
						types.InlineElements{
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
								Content: " ",
							},
							types.SpecialCharacter{
								Name: "<",
							},
							types.StringElement{
								Content: "using the ",
							},
							types.QuotedText{
								Kind: types.SingleQuoteBold,
								Elements: []interface{}{
									types.StringElement{
										Content: "inline link macro",
									},
								},
							},
							types.SpecialCharacter{
								Name: ">",
							},
						},
						types.InlineElements{
							types.StringElement{
								Content: "another one using attribute substitution: ",
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
								Content: "\u2026\u200b", // symbol for ellipsis, applied by the 'replacements' substitution
							},
						},
						types.InlineElements{
							types.SingleLineComment{
								Content: " a single-line comment.",
							},
						},
					}
					Expect(ParseRawSource(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'normal' substitution on multiple lines", func() {
					// quoted text is parsed but inline link macro is not
					s := strings.ReplaceAll(source, "$SUBS", "normal")
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						types.BlankLine{},
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrSubstitutions: "normal",
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
										Content: " ",
									},
									types.SpecialCharacter{
										Name: "<",
									},
									types.StringElement{
										Content: "using the ",
									},
									types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											types.StringElement{
												Content: "inline link macro",
											},
										},
									},
									types.SpecialCharacter{
										Name: ">",
									},
								},
								{
									types.StringElement{
										Content: "another one using attribute substitution: ",
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
										Content: "\u2026\u200b", // symbol for ellipsis, applied by the 'replacements' substitution
									},
								},
								{
									types.SingleLineComment{
										Content: " a single-line comment.",
									},
								},
							},
						},
					}
					Expect(ParseRawSource(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'quotes' substitution on multiple lines", func() {
					// quoted text is parsed but inline link macro is not
					s := strings.ReplaceAll(source, "$SUBS", "quotes")
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						types.BlankLine{},
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
								{
									types.SingleLineComment{
										Content: " a single-line comment.",
									},
								},
							},
						},
					}
					Expect(ParseRawSource(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'macros' substitution on multiple lines", func() {
					// quoted text is not parsed but inline link macro is
					s := strings.ReplaceAll(source, "$SUBS", "macros")
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						types.BlankLine{},
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
								{
									types.SingleLineComment{
										Content: " a single-line comment.",
									},
								},
							},
						},
					}
					Expect(ParseRawSource(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'attributes' substitution on multiple lines", func() {
					// quoted text is not parsed but inline link macro is
					s := strings.ReplaceAll(source, "$SUBS", "attributes")
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						types.BlankLine{},
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
								{
									types.SingleLineComment{
										Content: " a single-line comment.",
									},
								},
							},
						},
					}
					Expect(ParseRawSource(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'attributes,macros' substitution on multiple lines", func() {
					// quoted text is not parsed but inline link macro is
					s := strings.ReplaceAll(source, "$SUBS", "attributes,macros")
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						types.BlankLine{},
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
								{
									types.SingleLineComment{
										Content: " a single-line comment.",
									},
								},
							},
						},
					}
					Expect(ParseRawSource(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'specialchars' substitution on multiple lines", func() {
					// quoted text is not parsed but inline link macro is
					s := strings.ReplaceAll(source, "$SUBS", "specialchars")
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						types.BlankLine{},
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
								{
									types.SingleLineComment{
										Content: " a single-line comment.",
									},
								},
							},
						},
					}
					Expect(ParseRawSource(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'replacements' substitution on multiple lines", func() {
					// quoted text is not parsed but inline link macro is
					s := strings.ReplaceAll(source, "$SUBS", "replacements")
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						types.BlankLine{},
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
								{
									types.SingleLineComment{
										Content: " a single-line comment.",
									},
								},
							},
						},
					}
					Expect(ParseRawSource(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'quotes,macros' substitutions", func() {
					// quoted text and inline link macro are both parsed
					s := strings.ReplaceAll(source, "$SUBS", "quotes,macros")
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						types.BlankLine{},
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
								{
									types.SingleLineComment{
										Content: " a single-line comment.",
									},
								},
							},
						},
					}
					Expect(ParseRawSource(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'macros,quotes' substitutions", func() {
					// quoted text and inline link macro are both parsed
					// (same as above, but with subs in reversed order)
					s := strings.ReplaceAll(source, "$SUBS", "macros,quotes")
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						types.BlankLine{},
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
								{
									types.SingleLineComment{
										Content: " a single-line comment.",
									},
								},
							},
						},
					}
					Expect(ParseRawSource(s)).To(MatchDocumentFragments(expected))
				})

				It("should apply the 'none' substitution", func() {
					s := strings.ReplaceAll(source, "$SUBS", "none")
					expected := types.DocumentFragments{
						types.AttributeDeclaration{
							Name:  "github-url",
							Value: "https://github.com",
						},
						types.BlankLine{},
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
								{
									types.SingleLineComment{
										Content: " a single-line comment.",
									},
								},
							},
						},
					}
					Expect(ParseRawSource(s)).To(MatchDocumentFragments(expected))
				})
			})
		})

		Context("admonition paragraphs", func() {

			It("note admonition paragraph", func() {
				source := `NOTE: this is a note.`
				expected := types.DocumentFragments{
					types.Paragraph{
						Attributes: types.Attributes{
							types.AttrStyle: types.Note,
						},
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "this is a note."},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("warning admonition paragraph", func() {
				source := `WARNING: this is a multiline
warning!`
				expected := types.DocumentFragments{
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
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("admonition note paragraph with id and title", func() {
				source := `[[foo]]
.bar
NOTE: this is a note.`
				expected := types.DocumentFragments{
					types.Paragraph{
						Attributes: types.Attributes{
							types.AttrStyle: types.Note,
							types.AttrID:    "foo",
							types.AttrTitle: "bar",
						},
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "this is a note."},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("caution admonition paragraph with single line", func() {
				source := `[CAUTION]
this is a caution!`
				expected := types.DocumentFragments{
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
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("multiline caution admonition paragraph with title and id", func() {
				source := `[[foo]]
[CAUTION] 
.bar
this is a 
*caution*!`
				expected := types.DocumentFragments{
					types.Paragraph{
						Attributes: types.Attributes{
							types.AttrStyle: types.Caution,
							types.AttrID:    "foo",
							types.AttrTitle: "bar",
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
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("multiple admonition paragraphs", func() {
				source := `[NOTE]
No space after the [NOTE]!

[CAUTION]
And no space after [CAUTION] either.`
				expected := types.DocumentFragments{
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
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})
		})

		Context("verse paragraphs", func() {

			It("paragraph as a verse with author and title", func() {
				source := `[verse, john doe, verse title]
I am a verse paragraph.`
				expected := types.DocumentFragments{
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
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("paragraph as a verse with author, title and other attributes", func() {
				source := `[[universal]]
[verse, john doe, verse title]
.universe
I am a verse paragraph.`
				expected := types.DocumentFragments{
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
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("paragraph as a verse with empty title", func() {
				source := `[verse, john doe, ]
I am a verse paragraph.`
				expected := types.DocumentFragments{
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
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("paragraph as a verse without title", func() {
				source := `[verse, john doe ]
I am a verse paragraph.`
				expected := types.DocumentFragments{
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
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("paragraph as a verse with empty author", func() {
				source := `[verse,  ]
I am a verse paragraph.`
				expected := types.DocumentFragments{
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
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("paragraph as a verse without author", func() {
				source := `[verse]
I am a verse paragraph.`
				expected := types.DocumentFragments{
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
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("image block as a verse", func() {
				// assume that the author meant to use an image, so the `verse` attribute will be ignored during rendering
				source := `[verse, john doe, verse title]
image::foo.png[]`
				expected := types.DocumentFragments{
					types.Paragraph{
						Attributes: types.Attributes{
							types.AttrStyle:       types.Verse,
							types.AttrQuoteAuthor: "john doe",
							types.AttrQuoteTitle:  "verse title",
						},
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "image::foo.png[]"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})
		})

		Context("quote paragraphs", func() {

			It("paragraph as a quote with author and title", func() {
				source := `[quote, john doe, quote title]
I am a quote paragraph.`
				expected := types.DocumentFragments{
					types.Paragraph{
						Attributes: types.Attributes{
							types.AttrStyle:       types.Quote,
							types.AttrQuoteAuthor: "john doe",
							types.AttrQuoteTitle:  "quote title",
						},
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "I am a quote paragraph."},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("paragraph as a quote with author, title and other attributes", func() {
				source := `[[universal]]
[quote, john doe, quote title]
.universe
I am a quote paragraph.`
				expected := types.DocumentFragments{
					types.Paragraph{
						Attributes: types.Attributes{
							types.AttrStyle:       types.Quote,
							types.AttrQuoteAuthor: "john doe",
							types.AttrQuoteTitle:  "quote title",
							types.AttrID:          "universal",
							// types.AttrCustomID:    true,
							types.AttrTitle: "universe",
						},
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "I am a quote paragraph."},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("paragraph as a quote with empty title", func() {
				source := `[quote, john doe, ]
I am a quote paragraph.`
				expected := types.DocumentFragments{
					types.Paragraph{
						Attributes: types.Attributes{
							types.AttrStyle:       types.Quote,
							types.AttrQuoteAuthor: "john doe",
						},
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "I am a quote paragraph."},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("paragraph as a quote without title", func() {
				source := `[quote, john doe ]
I am a quote paragraph.`
				expected := types.DocumentFragments{
					types.Paragraph{
						Attributes: types.Attributes{
							types.AttrStyle:       types.Quote,
							types.AttrQuoteAuthor: "john doe",
						},
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "I am a quote paragraph."},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("paragraph as a quote with empty author", func() {
				source := `[quote,  ]
I am a quote paragraph.`
				expected := types.DocumentFragments{
					types.Paragraph{
						Attributes: types.Attributes{
							types.AttrStyle: types.Quote,
						},
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "I am a quote paragraph."},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("paragraph as a quote without author", func() {
				source := `[quote]
I am a quote paragraph.`
				expected := types.DocumentFragments{
					types.Paragraph{
						Attributes: types.Attributes{
							types.AttrStyle: types.Quote,
						},
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "I am a quote paragraph."},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("image block is NOT a quote", func() {
				Skip("needs clarification...")
				source := `[quote, john doe, quote title]
image::foo.png[]`
				expected := types.DocumentFragments{
					types.ImageBlock{
						Location: types.Location{
							Scheme: "",
							Path:   []interface{}{types.StringElement{Content: "foo.png"}},
						},
						Attributes: types.Attributes{
							types.AttrImageAlt: "quote",
							types.AttrWidth:    "john doe",
							types.AttrHeight:   "quote title",
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})
		})

		Context("thematic breaks", func() {

			It("thematic break form1 by itself", func() {
				source := "***"
				expected := types.DocumentFragments{
					types.ThematicBreak{},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("thematic break form2 by itself", func() {
				source := "* * *"
				expected := types.DocumentFragments{
					types.ThematicBreak{},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("thematic break form3 by itself", func() {
				source := "---"
				expected := types.DocumentFragments{
					types.ThematicBreak{},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("thematic break form4 by itself", func() {
				source := "- - -"
				expected := types.DocumentFragments{
					types.ThematicBreak{},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("thematic break form5 by itself", func() {
				source := "___"
				expected := types.DocumentFragments{
					types.ThematicBreak{},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("thematic break form4 by itself", func() {
				source := "_ _ _"
				expected := types.DocumentFragments{
					types.ThematicBreak{},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			It("thematic break with leading text", func() {
				source := "text ***"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "text ***"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})

			// NB: three asterisks gets confused with bullets if with trailing text
			It("thematic break with trailing text", func() {
				source := "* * * text"
				expected := types.DocumentFragments{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "* * * text"},
							},
						},
					},
				}
				Expect(ParseRawSource(source)).To(MatchDocumentFragments(expected))
			})
		})
	})

	Context("final documents", func() {

		Context("default paragraph", func() {

			It("with title attribute", func() {
				source := `[title=My Title]
foo
baz`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrTitle: "My Title",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "foo"},
								},
								{
									types.StringElement{Content: "baz"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with custom title attribute - explicit and unquoted", func() {
				source := `:title: cookies
				
[title=my {title}]
foo
baz`
				expected := types.Document{
					Elements: []interface{}{
						types.AttributeDeclaration{
							Name:  "title",
							Value: "cookies",
						},
						types.BlankLine{},
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrTitle: "my cookies",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "foo"},
								},
								{
									types.StringElement{Content: "baz"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with custom title attribute - explicit and single quoted", func() {
				source := `:title: cookies
				
[title='my {title}']
foo
baz`
				expected := types.Document{
					Elements: []interface{}{
						types.AttributeDeclaration{
							Name:  "title",
							Value: "cookies",
						},
						types.BlankLine{},
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrTitle: "my cookies",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "foo"},
								},
								{
									types.StringElement{Content: "baz"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with custom title attribute - explicit and double quoted", func() {
				source := `:title: cookies
				
[title="my {title}"]
foo
baz`
				expected := types.Document{
					Elements: []interface{}{
						types.AttributeDeclaration{
							Name:  "title",
							Value: "cookies",
						},
						types.BlankLine{},
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrTitle: "my cookies",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "foo"},
								},
								{
									types.StringElement{Content: "baz"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with custom title attribute - implicit", func() {
				source := `:title: cookies
				
.my {title}
foo
baz`
				expected := types.Document{
					Elements: []interface{}{
						types.AttributeDeclaration{
							Name:  "title",
							Value: "cookies",
						},
						types.BlankLine{},
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrTitle: "my cookies",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "foo"},
								},
								{
									types.StringElement{Content: "baz"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with multiple attributes without blanklines in-between", func() {
				source := `[%hardbreaks.role1.role2]
[#anchor]
foo
baz`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrID:      "anchor",
								types.AttrRoles:   []interface{}{"role1", "role2"},
								types.AttrOptions: []interface{}{"hardbreaks"},
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "foo"},
								},
								{
									types.StringElement{Content: "baz"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with multiple attributes and blanklines in-between", func() {
				source := `[%hardbreaks.role1.role2]

[#anchor]

foo
baz`
				expected := types.Document{
					Elements: []interface{}{
						types.BlankLine{}, // attribute not retained in blankline
						types.BlankLine{}, // attribute not retained in blankline
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "foo"},
								},
								{
									types.StringElement{Content: "baz"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with paragraph roles and attribute - case 1", func() {
				source := `[.role1%hardbreaks.role2]
foo
bar
baz`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrOptions: []interface{}{"hardbreaks"},
								types.AttrRoles:   []interface{}{"role1", "role2"},
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "foo"},
								},
								{
									types.StringElement{Content: "bar"},
								},
								{
									types.StringElement{Content: "baz"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with paragraph roles - case 2", func() {
				source := `[.role1%hardbreaks]
[.role2]
foo
bar
baz`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Attributes: types.Attributes{
								types.AttrOptions: []interface{}{"hardbreaks"},
								types.AttrRoles:   []interface{}{"role1", "role2"},
							},
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "foo"},
								},
								{
									types.StringElement{Content: "bar"},
								},
								{
									types.StringElement{Content: "baz"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("not treat plusplus as line break", func() {
				source := `C++
foo`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "C++"},
								},
								{
									types.StringElement{Content: "foo"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			Context("with counters", func() {

				It("default", func() {
					source := `foo{counter:foo} bar{counter2:foo} baz{counter:foo} bob{counter:bob}`
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "foo1 bar baz3 bob1"},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("with numeric start", func() {
					source := `foo{counter:foo:2} bar{counter2:foo} baz{counter:foo} bob{counter:bob:10}`
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "foo2 bar baz4 bob10"},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("with alphanumeric start", func() {
					source := `foo{counter:foo:b} bar{counter2:foo} baz{counter:foo} bob{counter:bob:z}`
					expected := types.Document{
						Attributes: types.Attributes{
							types.AttrIDPrefix: "bar_",
						},
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "foob bar bazd bobz"},
									},
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
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "a paragraph"},
								},
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
							Lines: [][]interface{}{
								{
									types.PredefinedAttribute{
										Name: "blank",
									},
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
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "hello "},
									types.PredefinedAttribute{Name: "plus"},
									types.StringElement{Content: " world"},
								},
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

		Context("quote paragraphs", func() {

			It("inline image within a quote", func() {
				source := `[quote, john doe, quote title]
a foo image:foo.png[]`
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
										Content: "a foo ",
									},
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
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})

	})
})
