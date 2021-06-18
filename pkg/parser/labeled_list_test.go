package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("labeled lists", func() {

	Context("in final documents", func() {

		It("with a term and description on 2 lines", func() {
			source := `Item1::
Item 1 description
on 2 lines`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item1",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{
												Content: "Item 1 description\non 2 lines",
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

		It("with a single term and no description", func() {
			source := `Item1::`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item1",
									},
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("with a quoted text in term and in description", func() {
			source := "`foo()`::\n" +
				`This function is _untyped_.`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									&types.QuotedText{
										Kind: types.SingleQuoteMonospace,
										Elements: []interface{}{
											types.StringElement{
												Content: "foo()",
											},
										},
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{
												Content: "This function is ",
											},
											&types.QuotedText{
												Kind: types.SingleQuoteItalic,
												Elements: []interface{}{
													types.StringElement{
														Content: "untyped",
													},
												},
											},
											types.StringElement{
												Content: ".",
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
		It("with a index term", func() {
			source := "((`foo`))::\n" +
				`This function is _untyped_.`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.IndexTerm{
										Term: []interface{}{
											&types.QuotedText{
												Kind: types.SingleQuoteMonospace,
												Elements: []interface{}{
													types.StringElement{
														Content: "foo",
													},
												},
											},
										},
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{
												Content: "This function is ",
											},
											&types.QuotedText{
												Kind: types.SingleQuoteItalic,
												Elements: []interface{}{
													types.StringElement{
														Content: "untyped",
													},
												},
											},
											types.StringElement{
												Content: ".",
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

		It("with a concealed index term in term", func() {
			source := "(((foo,bar)))::\n" +
				`This function is _untyped_.`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.ConcealedIndexTerm{
										Term1: "foo",
										Term2: "bar",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{
												Content: "This function is ",
											},
											&types.QuotedText{
												Kind: types.SingleQuoteItalic,
												Elements: []interface{}{
													types.StringElement{
														Content: "untyped",
													},
												},
											},
											types.StringElement{
												Content: ".",
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

		It("with a horizontal layout attribute", func() {
			source := `[horizontal]
Item1:: foo`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Attributes: types.Attributes{
							"style": "horizontal",
						},
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item1",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{Content: "foo"},
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

		It("with a single term and a blank line", func() {
			source := `Item1::
			`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item1",
									},
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("with multiple sibling items", func() {
			source := `Item 1::
Item 1 description
Item 2:: 
Item 2 description
Item 3:: 
Item 3 description`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item 1",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{Content: "Item 1 description"},
										},
									},
								},
							},
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item 2",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{Content: "Item 2 description"},
										},
									},
								},
							},
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item 3",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{Content: "Item 3 description"},
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

		It("with multiple nested items", func() {
			source := `Item 1::
Item 1 description
Item 2:::
Item 2 description
Item 3::::
Item 3 description`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item 1",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{Content: "Item 1 description"},
										},
									},
									&types.GenericList{
										Kind: types.LabeledListKind,
										Elements: []interface{}{
											&types.LabeledListElement{
												Style: types.TripleColons,
												Term: []interface{}{
													types.StringElement{
														Content: "Item 2",
													},
												},
												Elements: []interface{}{
													&types.Paragraph{
														Elements: []interface{}{
															types.StringElement{Content: "Item 2 description"},
														},
													},
													&types.GenericList{
														Kind: types.LabeledListKind,
														Elements: []interface{}{
															&types.LabeledListElement{
																Style: types.QuadrupleColons,
																Term: []interface{}{
																	types.StringElement{
																		Content: "Item 3",
																	},
																},
																Elements: []interface{}{
																	&types.Paragraph{
																		Elements: []interface{}{
																			types.StringElement{Content: "Item 3 description"},
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
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("with nested unordered list - case 1", func() {
			source := `Empty item:: 
* foo
* bar
Item with description:: something simple`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Empty item",
									},
								},
								Elements: []interface{}{
									&types.GenericList{
										Kind: types.UnorderedListKind,
										Elements: []interface{}{
											&types.UnorderedListElement{
												BulletStyle: types.OneAsterisk,
												CheckStyle:  types.NoCheck,
												Elements: []interface{}{
													&types.Paragraph{
														Elements: []interface{}{
															types.StringElement{Content: "foo"},
														},
													},
												},
											},
											&types.UnorderedListElement{
												BulletStyle: types.OneAsterisk,
												CheckStyle:  types.NoCheck,
												Elements: []interface{}{
													&types.Paragraph{
														Elements: []interface{}{
															types.StringElement{Content: "bar"},
														},
													},
												},
											},
										},
									},
								},
							},
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item with description",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{Content: "something simple"},
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

		It("with a single item and paragraph", func() {
			source := `Item 1::
foo
bar

a normal paragraph.`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item 1",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{Content: "foo\nbar"},
										},
									},
								},
							},
						},
					},

					&types.Paragraph{
						Elements: []interface{}{
							types.StringElement{Content: "a normal paragraph."},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("with item continuation", func() {
			source := `Item 1::
+
----
a fenced block
----
Item 2:: something simple
+
----
another fenced block
----`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item 1",
									},
								},
								Elements: []interface{}{
									&types.ListElementContinuation{},
									&types.DelimitedBlock{
										Kind: types.Listing,
										Elements: []interface{}{
											types.StringElement{
												Content: "a fenced block",
											},
										},
									},
								},
							},
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item 2",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{Content: "something simple"},
										},
									},
									&types.ListElementContinuation{},
									&types.DelimitedBlock{
										Kind: types.Listing,
										Elements: []interface{}{
											types.StringElement{
												Content: "another fenced block",
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

		It("without item continuation", func() {
			source := `Item 1::
----
a fenced block
----
Item 2:: something simple
----
another fenced block
----`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item 1",
									},
								},
							},
						},
					},
					&types.DelimitedBlock{
						Kind: types.Listing,
						Elements: []interface{}{
							types.StringElement{
								Content: "a fenced block",
							},
						},
					},
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item 2",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{Content: "something simple"},
										},
									},
								},
							},
						},
					},
					&types.DelimitedBlock{
						Kind: types.Listing,
						Elements: []interface{}{
							types.StringElement{
								Content: "another fenced block",
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("with nested unordered list - case 2", func() {
			source := `Labeled item::
- unordered item`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Labeled item",
									},
								},
								Elements: []interface{}{
									&types.GenericList{
										Kind: types.UnorderedListKind,
										Elements: []interface{}{
											&types.UnorderedListElement{
												BulletStyle: types.Dash,
												CheckStyle:  types.NoCheck,
												Elements: []interface{}{
													&types.Paragraph{
														Elements: []interface{}{
															types.StringElement{Content: "unordered item"},
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
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("with title", func() {
			source := `.Labeled, single-line
first term:: definition of the first term
second term:: definition of the second term`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Attributes: types.Attributes{
							types.AttrTitle: "Labeled, single-line",
						},
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "first term",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{
												Content: "definition of the first term",
											},
										},
									},
								},
							},
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "second term",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{
												Content: "definition of the second term",
											},
										},
									},
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

		It("max level of labeled items - case 1", func() {
			source := `.Labeled, max nesting
level 1:: description 1
level 2::: description 2
level 3:::: description 3
level 1:: description 1`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Attributes: types.Attributes{
							types.AttrTitle: "Labeled, max nesting",
						},
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "level 1",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{
												Content: "description 1",
											},
										},
									},
									&types.GenericList{
										Kind: types.LabeledListKind,
										Elements: []interface{}{
											&types.LabeledListElement{
												Style: types.TripleColons,
												Term: []interface{}{
													types.StringElement{
														Content: "level 2",
													},
												},
												Elements: []interface{}{
													&types.Paragraph{
														Elements: []interface{}{
															types.StringElement{
																Content: "description 2",
															},
														},
													},
													&types.GenericList{
														Kind: types.LabeledListKind,
														Elements: []interface{}{
															&types.LabeledListElement{
																Style: types.QuadrupleColons,
																Term: []interface{}{
																	types.StringElement{
																		Content: "level 3",
																	},
																},
																Elements: []interface{}{
																	&types.Paragraph{
																		Elements: []interface{}{
																			types.StringElement{
																				Content: "description 3",
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
									},
								},
							},
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "level 1",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{
												Content: "description 1",
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

		It("max level of labeled items - case 2", func() {
			source := `.Labeled, max nesting
level 1:: description 1
level 2::: description 2
level 3:::: description 3
level 2::: description 2`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Attributes: types.Attributes{
							types.AttrTitle: "Labeled, max nesting",
						},
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "level 1",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{
												Content: "description 1",
											},
										},
									},
									&types.GenericList{
										Kind: types.LabeledListKind,
										Elements: []interface{}{
											&types.LabeledListElement{
												Style: types.TripleColons,
												Term: []interface{}{
													types.StringElement{
														Content: "level 2",
													},
												},
												Elements: []interface{}{
													&types.Paragraph{
														Elements: []interface{}{
															types.StringElement{
																Content: "description 2",
															},
														},
													},
													&types.GenericList{
														Kind: types.LabeledListKind,
														Elements: []interface{}{
															&types.LabeledListElement{
																Style: types.QuadrupleColons,
																Term: []interface{}{
																	types.StringElement{
																		Content: "level 3",
																	},
																},
																Elements: []interface{}{
																	&types.Paragraph{
																		Elements: []interface{}{
																			types.StringElement{
																				Content: "description 3",
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
											&types.LabeledListElement{
												Style: types.TripleColons,
												Term: []interface{}{
													types.StringElement{
														Content: "level 2",
													},
												},
												Elements: []interface{}{
													&types.Paragraph{
														Elements: []interface{}{
															types.StringElement{
																Content: "description 2",
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
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("item with predefined attribute", func() {
			source := `level 1:: {amp}`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "level 1",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.PredefinedAttribute{Name: "amp"},
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

		It("item with a colon the term", func() {
			source := `what: ever:: text`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "what: ever",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{Content: "text"},
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

		It("with multiple item continuations", func() {
			source := `Item 1::
content 1
+
NOTE: note

Item 2::
content 2
+
addition
+
IMPORTANT: important
+
TIP: tip`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item 1",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{Content: "content 1"},
										},
									},
									&types.ListElementContinuation{},
									&types.Paragraph{
										Attributes: types.Attributes{
											types.AttrStyle: types.Note,
										},
										Elements: []interface{}{
											types.StringElement{Content: "note"},
										},
									},
								},
							},
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "Item 2",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{Content: "content 2"},
										},
									},
									&types.ListElementContinuation{},
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{Content: "addition"},
										},
									},
									&types.ListElementContinuation{},
									&types.Paragraph{
										Attributes: types.Attributes{
											types.AttrStyle: types.Important,
										},
										Elements: []interface{}{
											types.StringElement{Content: "important"},
										},
									},
									&types.ListElementContinuation{},
									&types.Paragraph{
										Attributes: types.Attributes{
											types.AttrStyle: types.Tip,
										},
										Elements: []interface{}{
											types.StringElement{Content: "tip"},
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

		It("with list item continuations", func() {
			source := `item::
This is the first line of the first paragraph.
This is the second line of the first paragraph.
+
This is the first line of the continuation paragraph.
This is the second line of the continuation paragraph.
+
This is the next continuation paragraph.
+
TIP: We can embed admonitions too!
`
			expected := types.Document{
				Elements: []interface{}{
					&types.GenericList{
						Kind: types.LabeledListKind,
						Elements: []interface{}{
							&types.LabeledListElement{
								Style: types.DoubleColons,
								Term: []interface{}{
									types.StringElement{
										Content: "item",
									},
								},
								Elements: []interface{}{
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{
												Content: "This is the first line of the first paragraph.\nThis is the second line of the first paragraph.",
											},
										},
									},
									&types.ListElementContinuation{},
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{
												Content: "This is the first line of the continuation paragraph.\nThis is the second line of the continuation paragraph.",
											},
										},
									},
									&types.ListElementContinuation{},
									&types.Paragraph{
										Elements: []interface{}{
											types.StringElement{
												Content: "This is the next continuation paragraph.",
											},
										},
									},
									&types.ListElementContinuation{},
									&types.Paragraph{
										Attributes: types.Attributes{
											types.AttrStyle: types.Tip,
										},
										Elements: []interface{}{
											types.StringElement{
												Content: "We can embed admonitions too!",
											},
										},
									},
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
	})

})
