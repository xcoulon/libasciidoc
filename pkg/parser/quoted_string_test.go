package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("quoted strings", func() {

	Context("in final documents", func() {

		It("simple single quoted string", func() {
			source := "'`curly was single`'"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.SingleQuote,
								Elements: []interface{}{
									&types.StringElement{Content: "curly was single"},
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
		It("interior spaces with single quoted string", func() {
			source := "'` curly was single `'"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.StringElement{Content: "'` curly was single \u2019"},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("interior ending space with single quoted string", func() {
			source := "'`curly was single `'"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.StringElement{Content: "'`curly was single \u2019"},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("interior leading space with single quoted string", func() {
			source := "'` curly was single`'"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.StringElement{Content: "'` curly was single\u2019"},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("bold in single quoted string", func() {
			source := "'`curly *was* single`'"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.SingleQuote,
								Elements: []interface{}{
									&types.StringElement{Content: "curly "},
									&types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											&types.StringElement{Content: "was"},
										},
									},
									&types.StringElement{Content: " single"},
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("italics in single quoted string", func() {
			source := "'`curly _was_ single`'"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.SingleQuote,
								Elements: []interface{}{
									&types.StringElement{Content: "curly "},
									&types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											&types.StringElement{Content: "was"},
										},
									},
									&types.StringElement{Content: " single"},
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
		It("span in single quoted string", func() {
			source := "'`curly [.strikeout]#was#_is_ single`'"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.SingleQuote,
								Elements: []interface{}{
									&types.StringElement{
										Content: "curly ",
									},
									&types.QuotedText{
										Kind: types.SingleQuoteMarked,
										Attributes: types.Attributes{
											types.AttrRoles: []interface{}{"strikeout"},
										},
										Elements: []interface{}{
											&types.StringElement{
												Content: "was",
											},
										},
									},
									&types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											&types.StringElement{
												Content: "is",
											},
										},
									},

									&types.StringElement{
										Content: " single",
									},
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
		It("curly in monospace string", func() {
			source := "'`curly `is` single`'"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.SingleQuote,
								Elements: []interface{}{
									&types.StringElement{Content: "curly "},
									&types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											&types.StringElement{Content: "is"},
										},
									},
									&types.StringElement{Content: " single"},
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
		It("curly as monospace string", func() {
			source := "'``curly``'"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.SingleQuote,
								Elements: []interface{}{
									&types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											&types.StringElement{Content: "curly"},
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

		It("curly with nested double curly", func() {
			source := "'`single\"`double`\"`'"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.SingleQuote,
								Elements: []interface{}{
									&types.StringElement{Content: "single"},
									&types.QuotedString{
										Kind: types.DoubleQuote,
										Elements: []interface{}{
											&types.StringElement{Content: "double"},
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

		It("curly in monospace string", func() {
			source := "`'`curly`'`"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedText{
								Kind: types.SingleQuoteMonospace,
								Elements: []interface{}{
									&types.QuotedString{
										Kind: types.SingleQuote,
										Elements: []interface{}{
											&types.StringElement{Content: "curly"},
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
		It("curly in italics", func() {
			source := "_'`curly`'_"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedText{
								Kind: types.SingleQuoteItalic,
								Elements: []interface{}{
									&types.QuotedString{
										Kind: types.SingleQuote,
										Elements: []interface{}{
											&types.StringElement{Content: "curly"},
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
		It("curly in bold", func() {
			source := "*'`curly`'*"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedText{
								Kind: types.SingleQuoteBold,
								Elements: []interface{}{
									&types.QuotedString{
										Kind: types.SingleQuote,
										Elements: []interface{}{
											&types.StringElement{Content: "curly"},
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

		It("curly in link", func() {
			source := "https://www.example.com/a['`example`']"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.InlineLink{
								Location: &types.Location{
									Scheme: "https://",
									Path: []interface{}{
										&types.StringElement{Content: "www.example.com/a"},
									},
								},
								Attributes: types.Attributes{
									types.AttrInlineLinkText: []interface{}{
										&types.QuotedString{
											Kind: types.SingleQuote,
											Elements: []interface{}{
												&types.StringElement{
													Content: "example",
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
		It("curly in quoted link", func() {
			source := "https://www.example.com/a[\"an '`example`'\"]"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.InlineLink{
								Location: &types.Location{
									Scheme: "https://",
									Path: []interface{}{
										&types.StringElement{Content: "www.example.com/a"},
									},
								},
								Attributes: types.Attributes{
									types.AttrInlineLinkText: []interface{}{
										&types.StringElement{
											Content: "an ",
										},
										&types.QuotedString{
											Kind: types.SingleQuote,
											Elements: []interface{}{
												&types.StringElement{
													Content: "example",
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

		It("image in curly", func() {
			source := "'`a image:foo.png[]`'"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.SingleQuote,
								Elements: []interface{}{
									&types.StringElement{Content: "a "},
									types.InlineImage{
										Location: &types.Location{
											Path: []interface{}{
												&types.StringElement{
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
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("icon in curly", func() {
			source := "'`a icon:note[]`'"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.SingleQuote,
								Elements: []interface{}{
									&types.StringElement{Content: "a "},
									types.Icon{
										Class: "note",
									},
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("simple double quoted string", func() {
			source := "\"`curly was single`\""
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.DoubleQuote,
								Elements: []interface{}{
									&types.StringElement{Content: "curly was single"},
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("interior spaces with double quoted string", func() {
			source := "\"` curly was single `\""
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.StringElement{Content: "\"` curly was single `\""},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
		It("interior ending space with double quoted string", func() {
			source := "\"`curly was single `\""
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.StringElement{Content: "\"`curly was single `\""},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
		It("interior leading space with double quoted string", func() {
			source := "\"` curly was single`\""
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.StringElement{Content: "\"` curly was single`\""},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("bold in double quoted string", func() {
			source := "\"`curly *was* single`\""
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.DoubleQuote,
								Elements: []interface{}{
									&types.StringElement{Content: "curly "},
									&types.QuotedText{
										Kind: types.SingleQuoteBold,
										Elements: []interface{}{
											&types.StringElement{Content: "was"},
										},
									},
									&types.StringElement{Content: " single"},
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
		It("italics in double quoted string", func() {
			source := "\"`curly _was_ single`\""
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.DoubleQuote,
								Elements: []interface{}{
									&types.StringElement{Content: "curly "},
									&types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											&types.StringElement{Content: "was"},
										},
									},
									&types.StringElement{Content: " single"},
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("span in double quoted string", func() {
			source := "\"`curly [.strikeout]#was#_is_ single`\""
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.DoubleQuote,
								Elements: []interface{}{
									&types.StringElement{
										Content: "curly ",
									},
									&types.QuotedText{
										Kind: types.SingleQuoteMarked,
										Attributes: types.Attributes{
											types.AttrRoles: []interface{}{"strikeout"},
										},
										Elements: []interface{}{
											&types.StringElement{
												Content: "was",
											},
										},
									},
									&types.QuotedText{
										Kind: types.SingleQuoteItalic,
										Elements: []interface{}{
											&types.StringElement{
												Content: "is",
											},
										},
									},
									&types.StringElement{
										Content: " single",
									},
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("double curly in monospace string", func() {
			source := "\"`curly `is` single`\""
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.DoubleQuote,
								Elements: []interface{}{
									&types.StringElement{Content: "curly "},
									&types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											&types.StringElement{Content: "is"},
										},
									},
									&types.StringElement{Content: " single"},
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
		It("double curly as monospace string", func() {
			source := "\"``curly``\""
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.DoubleQuote,
								Elements: []interface{}{
									&types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											&types.StringElement{Content: "curly"},
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
		It("double curly with nested single curly", func() {
			source := "\"`double'`single`'`\""
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.DoubleQuote,
								Elements: []interface{}{
									&types.StringElement{Content: "double"},
									&types.QuotedString{
										Kind: types.SingleQuote,
										Elements: []interface{}{
											&types.StringElement{Content: "single"},
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
		It("double curly in monospace string", func() {
			source := "`\"`curly`\"`"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedText{
								Kind: types.SingleQuoteMonospace,
								Elements: []interface{}{
									&types.QuotedString{
										Kind: types.DoubleQuote,
										Elements: []interface{}{
											&types.StringElement{Content: "curly"},
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
		It("double curly in italics", func() {
			source := "_\"`curly`\"_"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedText{
								Kind: types.SingleQuoteItalic,
								Elements: []interface{}{
									&types.QuotedString{
										Kind: types.DoubleQuote,
										Elements: []interface{}{
											&types.StringElement{Content: "curly"},
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
		It("double curly in bold", func() {
			source := "*\"`curly`\"*"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedText{
								Kind: types.SingleQuoteBold,
								Elements: []interface{}{
									&types.QuotedString{
										Kind: types.DoubleQuote,
										Elements: []interface{}{
											&types.StringElement{Content: "curly"},
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

		// In a link, the quotes are ambiguous, and we default to assuming they are for enclosing
		// the link text.  Nest them explicitly if this is needed.
		It("double curly in link (becomes mono)", func() {
			source := "https://www.example.com/a[\"`example`\"]"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.InlineLink{
								Location: &types.Location{
									Scheme: "https://",
									Path: []interface{}{
										&types.StringElement{Content: "www.example.com/a"},
									},
								},
								Attributes: types.Attributes{
									types.AttrInlineLinkText: []interface{}{
										&types.QuotedString{
											Kind: types.DoubleQuote,
											Elements: []interface{}{
												&types.StringElement{
													Content: "example",
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

		// This is the unambiguous form.
		It("curly in quoted link", func() {
			source := "https://www.example.com/a[\"\"`example`\"\"]"
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.InlineLink{
								Location: &types.Location{
									Scheme: "https://",
									Path: []interface{}{
										&types.StringElement{Content: "www.example.com/a"},
									},
								},
								Attributes: types.Attributes{
									types.AttrInlineLinkText: []interface{}{
										&types.QuotedString{
											Kind: types.DoubleQuote,
											Elements: []interface{}{
												&types.StringElement{
													Content: "example",
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
		It("image in double curly", func() {
			source := "\"`a image:foo.png[]`\""
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.DoubleQuote,
								Elements: []interface{}{
									&types.StringElement{Content: "a "},
									types.InlineImage{
										Location: &types.Location{
											Path: []interface{}{
												&types.StringElement{
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
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
		It("icon in double curly", func() {
			source := "\"`a icon:note[]`\""
			expected := types.Document{
				Elements: []interface{}{
					&types.Paragraph{
						Elements: []interface{}{
							&types.QuotedString{
								Kind: types.DoubleQuote,
								Elements: []interface{}{
									&types.StringElement{Content: "a "},
									types.Icon{
										Class: "note",
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

	It("curly in title", func() {
		source := "== a '`curly`' episode"
		expected := types.Document{
			Elements: []interface{}{
				types.Section{
					Attributes: types.Attributes{
						types.AttrID: "_a_episode",
					},
					Title: []interface{}{
						&types.StringElement{Content: "a "},
						&types.QuotedString{
							Kind: types.SingleQuote,
							Elements: []interface{}{
								&types.StringElement{Content: "curly"},
							},
						},
						&types.StringElement{Content: " episode"},
					},
					Elements: []interface{}{},
				},
			},
		}
		Expect(ParseDocument(source)).To(MatchDocument(expected))
	})

	It("curly in list element", func() {
		source := "* a '`curly`' episode"
		expected := types.Document{
			Elements: []interface{}{
				&types.GenericList{
					Kind: types.OrderedListKind,
					Elements: []types.ListElement{
						&types.UnorderedListElement{
							CheckStyle:  types.NoCheck,
							BulletStyle: types.OneAsterisk,
							Elements: []interface{}{
								&types.Paragraph{
									Elements: []interface{}{
										&types.StringElement{Content: "a "},
										&types.QuotedString{
											Kind: types.SingleQuote,
											Elements: []interface{}{
												&types.StringElement{Content: "curly"},
											},
										},
										&types.StringElement{Content: " episode"},
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

	It("curly in labeled list", func() {
		source := "'`term`':: something '`quoted`'"
		expected := types.Document{
			Elements: []interface{}{
				&types.GenericList{
					Kind: types.LabeledListKind,
					Elements: []types.ListElement{
						&types.LabeledListElement{
							Term: []interface{}{
								&types.StringElement{Content: "'`term`'"}, // parsed later
							},
							Elements: []interface{}{
								&types.Paragraph{
									Elements: []interface{}{
										&types.StringElement{Content: "something "},
										&types.QuotedString{
											Kind: types.SingleQuote,
											Elements: []interface{}{
												&types.StringElement{Content: "quoted"},
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

	It("double curly in title", func() {
		source := "== a \"`curly`\" episode"
		expected := types.Document{
			Elements: []interface{}{
				types.Section{
					Attributes: types.Attributes{
						types.AttrID: "_a_episode",
					},
					Title: []interface{}{
						&types.StringElement{Content: "a "},
						&types.QuotedString{
							Kind: types.DoubleQuote,
							Elements: []interface{}{
								&types.StringElement{Content: "curly"},
							},
						},
						&types.StringElement{Content: " episode"},
					},
					Elements: []interface{}{},
				},
			},
		}
		Expect(ParseDocument(source)).To(MatchDocument(expected))
	})

	It("double curly in labeled list", func() {
		source := "\"`term`\":: something \"`quoted`\""
		expected := types.Document{
			Elements: []interface{}{
				types.LabeledListElement{
					Term: []interface{}{
						&types.StringElement{Content: "\"`term`\""}, // parsed later
					},
					Elements: []interface{}{
						&types.Paragraph{
							Elements: []interface{}{
								&types.StringElement{Content: "something "},
								&types.QuotedString{
									Kind: types.DoubleQuote,
									Elements: []interface{}{
										&types.StringElement{Content: "quoted"},
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

	It("double in list element", func() {
		source := "* a \"`curly`\" episode"
		expected := types.Document{
			Elements: []interface{}{
				&types.GenericList{
					Kind: types.OrderedListKind,
					Elements: []types.ListElement{
						&types.UnorderedListElement{
							CheckStyle:  types.NoCheck,
							BulletStyle: types.OneAsterisk,
							Elements: []interface{}{
								&types.Paragraph{
									Elements: []interface{}{
										&types.StringElement{Content: "a "},
										&types.QuotedString{
											Kind: types.DoubleQuote,
											Elements: []interface{}{
												&types.StringElement{Content: "curly"},
											},
										},
										&types.StringElement{Content: " episode"},
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
