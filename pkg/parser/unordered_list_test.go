package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("unordered lists", func() {

	Context("in final documents", func() {

		Context("valid content", func() {

			It("with a basic single item", func() {
				source := `* a list item`
				expected := types.Document{
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
												types.StringElement{Content: "a list item"},
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

			It("with ID, title, role and a single item", func() {
				source := `.mytitle
[#listID]
[.myrole]
* a list item`
				expected := types.Document{
					Elements: []interface{}{
						&types.GenericList{
							Kind: types.UnorderedListKind,
							Attributes: types.Attributes{
								types.AttrID:    "listID",
								types.AttrTitle: "mytitle",
								types.AttrRoles: []interface{}{"myrole"},
							},
							Elements: []interface{}{
								&types.UnorderedListElement{
									BulletStyle: types.OneAsterisk,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												types.StringElement{Content: "a list item"},
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

			It("with a title and a single item", func() {
				source := `.a title
	* a list item`
				expected := types.Document{
					Elements: []interface{}{
						&types.GenericList{
							Kind: types.UnorderedListKind,
							Attributes: types.Attributes{
								types.AttrTitle: "a title",
							},
							Elements: []interface{}{
								&types.UnorderedListElement{
									BulletStyle: types.OneAsterisk,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												types.StringElement{Content: "a list item"},
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

			It("with 2 items with stars", func() {
				source := `* a first item
					* a second item with *bold content*`
				expected := types.Document{
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
												types.StringElement{Content: "a first item"},
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
												types.StringElement{Content: "a second item with "},
												&types.QuotedText{
													Kind: types.SingleQuoteBold,
													Elements: []interface{}{
														types.StringElement{Content: "bold content"},
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

			It("with multiple levels", func() {
				source := `.Unordered list title
		* list item 1
		** nested list item A
		*** nested nested list item A.1
		*** nested nested list item A.2
		** nested list item B
		*** nested nested list item B.1
		*** nested nested list item B.2
		* list item 2`
				expected := types.Document{
					Elements: []interface{}{
						&types.GenericList{
							Kind: types.UnorderedListKind,
							Attributes: types.Attributes{
								types.AttrTitle: "Unordered list title",
							},
							Elements: []interface{}{
								&types.UnorderedListElement{
									BulletStyle: types.OneAsterisk,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												types.StringElement{Content: "list item 1"},
											},
										},
										&types.GenericList{
											Kind: types.UnorderedListKind,
											Elements: []interface{}{
												&types.UnorderedListElement{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "nested list item A"},
															},
														},
														&types.GenericList{
															Kind: types.UnorderedListKind,
															Elements: []interface{}{
																&types.UnorderedListElement{
																	BulletStyle: types.ThreeAsterisks,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		&types.Paragraph{
																			Elements: []interface{}{
																				types.StringElement{Content: "nested nested list item A.1"},
																			},
																		},
																	},
																},
																&types.UnorderedListElement{
																	BulletStyle: types.ThreeAsterisks,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		&types.Paragraph{
																			Elements: []interface{}{
																				types.StringElement{Content: "nested nested list item A.2"},
																			},
																		},
																	},
																},
															},
														},
													},
												},
												&types.UnorderedListElement{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "nested list item B"},
															},
														},
														&types.GenericList{
															Kind: types.UnorderedListKind,
															Elements: []interface{}{
																&types.UnorderedListElement{
																	BulletStyle: types.ThreeAsterisks,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		&types.Paragraph{
																			Elements: []interface{}{
																				types.StringElement{Content: "nested nested list item B.1"},
																			},
																		},
																	},
																},
																&types.UnorderedListElement{
																	BulletStyle: types.ThreeAsterisks,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		&types.Paragraph{
																			Elements: []interface{}{
																				types.StringElement{Content: "nested nested list item B.2"},
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
								&types.UnorderedListElement{
									BulletStyle: types.OneAsterisk,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												types.StringElement{Content: "list item 2"},
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

			It("with 2 items with carets", func() {
				source := "- a first item\n" +
					"- a second item with *bold content*"
				expected := types.Document{
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
												types.StringElement{Content: "a first item"},
											},
										},
									},
								},
								&types.UnorderedListElement{
									BulletStyle: types.Dash,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												types.StringElement{Content: "a second item with "},
												&types.QuotedText{
													Kind: types.SingleQuoteBold,
													Elements: []interface{}{
														types.StringElement{Content: "bold content"},
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

			It("with items with mixed styles", func() {
				source := `- a parent item
					* a child item
					- another parent item
					* another child item
					** with a sub child item`
				expected := types.Document{
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
												types.StringElement{Content: "a parent item"},
											},
										},
										&types.GenericList{
											Kind: types.UnorderedListKind,
											Elements: []interface{}{
												&types.UnorderedListElement{
													BulletStyle: types.OneAsterisk,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "a child item"},
															},
														},
													},
												},
											},
										},
									},
								},
								&types.UnorderedListElement{
									BulletStyle: types.Dash,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												types.StringElement{Content: "another parent item"},
											},
										},
										&types.GenericList{
											Kind: types.UnorderedListKind,
											Elements: []interface{}{
												&types.UnorderedListElement{
													BulletStyle: types.OneAsterisk,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "another child item"},
															},
														},
														&types.GenericList{
															Kind: types.UnorderedListKind,
															Elements: []interface{}{
																&types.UnorderedListElement{
																	BulletStyle: types.TwoAsterisks,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		&types.Paragraph{
																			Elements: []interface{}{
																				types.StringElement{Content: "with a sub child item"},
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

			It("with 2 items with empty line in-between", func() {
				// fist line after list item is swallowed
				source := "* a first item\n" +
					"\n" +
					"* a second item with *bold content*"
				expected := types.Document{
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
												types.StringElement{Content: "a first item"},
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
												types.StringElement{Content: "a second item with "},
												&types.QuotedText{
													Kind: types.SingleQuoteBold,
													Elements: []interface{}{
														types.StringElement{Content: "bold content"},
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
			It("with 2 items on multiple lines", func() {
				source := `* item 1
  on 2 lines.
* item 2
on 2 lines, too.`
				expected := types.Document{
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
												types.StringElement{Content: "item 1\n  on 2 lines."}, // heading spaces are not trimmed
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
												types.StringElement{Content: "item 2\non 2 lines, too."},
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
			It("unordered lists with 2 empty lines in-between", func() {
				// the first blank lines after the first list is swallowed (for the list item)
				source := "* an item in the first list\n" +
					"\n" +
					"\n" +
					"* an item in the second list"
				expected := types.Document{
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
												types.StringElement{Content: "an item in the first list"},
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
												types.StringElement{Content: "an item in the second list"},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected)) // parse the whole document to get 2 lists
			})

			It("with items on 3 levels", func() {
				source := `* item 1
	** item 1.1
	** item 1.2
	*** item 1.2.1
	** item 1.3
	** item 1.4
	* item 2
	** item 2.1`
				expected := types.Document{
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
												types.StringElement{Content: "item 1"},
											},
										},
										&types.GenericList{
											Kind: types.UnorderedListKind,
											Elements: []interface{}{
												&types.UnorderedListElement{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "item 1.1"},
															},
														},
													},
												},
												&types.UnorderedListElement{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "item 1.2"},
															},
														},
														&types.GenericList{
															Kind: types.UnorderedListKind,
															Elements: []interface{}{
																&types.UnorderedListElement{
																	BulletStyle: types.ThreeAsterisks,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		&types.Paragraph{
																			Elements: []interface{}{
																				types.StringElement{Content: "item 1.2.1"},
																			},
																		},
																	},
																},
															},
														},
													},
												},
												&types.UnorderedListElement{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "item 1.3"},
															},
														},
													},
												},
												&types.UnorderedListElement{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "item 1.4"},
															},
														},
													},
												},
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
												types.StringElement{Content: "item 2"},
											},
										},
										&types.GenericList{
											Kind: types.UnorderedListKind,
											Elements: []interface{}{
												&types.UnorderedListElement{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "item 2.1"},
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

			It("max level of unordered items - case 1", func() {
				source := `.Unordered, max nesting
* level 1
** level 2
*** level 3
**** level 4
***** level 5
* level 1`
				expected := types.Document{
					Elements: []interface{}{
						&types.GenericList{
							Kind: types.UnorderedListKind,
							Attributes: types.Attributes{
								types.AttrTitle: "Unordered, max nesting",
							},
							Elements: []interface{}{
								&types.UnorderedListElement{
									BulletStyle: types.OneAsterisk,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												types.StringElement{Content: "level 1"},
											},
										},
										&types.GenericList{
											Kind: types.UnorderedListKind,
											Elements: []interface{}{
												&types.UnorderedListElement{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "level 2"},
															},
														},
														&types.GenericList{
															Kind: types.UnorderedListKind,
															Elements: []interface{}{
																&types.UnorderedListElement{
																	BulletStyle: types.ThreeAsterisks,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		&types.Paragraph{
																			Elements: []interface{}{
																				types.StringElement{Content: "level 3"},
																			},
																		},
																		&types.GenericList{
																			Kind: types.UnorderedListKind,
																			Elements: []interface{}{
																				&types.UnorderedListElement{
																					BulletStyle: types.FourAsterisks,
																					CheckStyle:  types.NoCheck,
																					Elements: []interface{}{
																						&types.Paragraph{
																							Elements: []interface{}{
																								types.StringElement{Content: "level 4"},
																							},
																						},
																						&types.GenericList{
																							Kind: types.UnorderedListKind,
																							Elements: []interface{}{
																								&types.UnorderedListElement{
																									BulletStyle: types.FiveAsterisks,
																									CheckStyle:  types.NoCheck,
																									Elements: []interface{}{
																										&types.Paragraph{
																											Elements: []interface{}{
																												types.StringElement{Content: "level 5"},
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
												},
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
												types.StringElement{Content: "level 1"},
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

			It("max level of unordered items - case 2", func() {
				source := `.Unordered, max nesting
* level 1
** level 2
*** level 3
**** level 4
***** level 5
** level 2`
				expected := types.Document{
					Elements: []interface{}{
						&types.GenericList{
							Kind: types.UnorderedListKind,
							Attributes: types.Attributes{
								types.AttrTitle: "Unordered, max nesting",
							},
							Elements: []interface{}{
								&types.UnorderedListElement{
									BulletStyle: types.OneAsterisk,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												types.StringElement{Content: "level 1"},
											},
										},
										&types.GenericList{
											Kind: types.UnorderedListKind,
											Elements: []interface{}{
												&types.UnorderedListElement{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "level 2"},
															},
														},
														&types.GenericList{
															Kind: types.UnorderedListKind,
															Elements: []interface{}{
																&types.UnorderedListElement{
																	BulletStyle: types.ThreeAsterisks,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		&types.Paragraph{
																			Elements: []interface{}{
																				types.StringElement{Content: "level 3"},
																			},
																		},
																		&types.GenericList{
																			Kind: types.UnorderedListKind,
																			Elements: []interface{}{
																				&types.UnorderedListElement{
																					BulletStyle: types.FourAsterisks,
																					CheckStyle:  types.NoCheck,
																					Elements: []interface{}{
																						&types.Paragraph{
																							Elements: []interface{}{
																								types.StringElement{Content: "level 4"},
																							},
																						},
																						&types.GenericList{
																							Kind: types.UnorderedListKind,
																							Elements: []interface{}{
																								&types.UnorderedListElement{
																									BulletStyle: types.FiveAsterisks,
																									CheckStyle:  types.NoCheck,
																									Elements: []interface{}{
																										&types.Paragraph{
																											Elements: []interface{}{
																												types.StringElement{Content: "level 5"},
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
												},
												&types.UnorderedListElement{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "level 2"},
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

			It("unordered list item with predefined attribute", func() {
				source := `* {amp}`
				expected := types.Document{
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
		})

		Context("invalid content", func() {
			It("with items on 3 levels and bad numbering", func() {
				source := `* item 1
					*** item 1.1
					*** item 1.1.1
					** item 1.2
					* item 2`
				expected := types.Document{
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
												types.StringElement{Content: "item 1"},
											},
										},
										&types.GenericList{
											Kind: types.UnorderedListKind,
											Elements: []interface{}{
												&types.UnorderedListElement{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "item 1.1"},
															},
														},
														&types.GenericList{
															Kind: types.UnorderedListKind,
															Elements: []interface{}{
																&types.UnorderedListElement{
																	BulletStyle: types.ThreeAsterisks,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		&types.Paragraph{
																			Elements: []interface{}{
																				types.StringElement{Content: "item 1.1.1"},
																			},
																		},
																	},
																},
															},
														},
													},
												},
												&types.UnorderedListElement{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "item 1.2"},
															},
														},
													},
												},
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
												types.StringElement{Content: "item 2"},
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

			It("invalid list item", func() {
				source := "*an invalid list item"
				expected := types.Document{
					Elements: []interface{}{
						&types.Paragraph{
							Elements: []interface{}{
								types.StringElement{Content: "*an invalid list item"},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})

		Context("list item continuation", func() {

			It("with item continuation - case 1", func() {
				source := `* foo
+
----
a delimited block
----
+
----
another delimited block
----
* bar
`
				expected := types.Document{
					Elements: []interface{}{
						types.UnorderedList{
							Items: []*types.UnorderedListElement{
								{
									BulletStyle: types.OneAsterisk,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												types.StringElement{Content: "foo"},
											},
										},
										types.ListingBlock{
											Elements: []interface{}{
												types.StringElement{
													Content: "a delimited block",
												},
											},
										},
										types.ListingBlock{
											Elements: []interface{}{
												types.StringElement{
													Content: "another delimited block",
												},
											},
										},
									},
								},
								{
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
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with item continuation - case 2", func() {
				source := `.Unordered, complex
* level 1
** level 2
*** level 3
This is a new line inside an unordered list using {plus} symbol.
We can even force content to start on a separate line... +
Amazing, isn't it?
**** level 4
+
The {plus} symbol is on a new line.

***** level 5
`
				expected := types.Document{
					Elements: []interface{}{
						types.UnorderedList{
							Attributes: types.Attributes{
								types.AttrTitle: "Unordered, complex",
							},
							Items: []*types.UnorderedListElement{
								{
									BulletStyle: types.OneAsterisk,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												types.StringElement{Content: "level 1"},
											},
										},
										types.UnorderedList{
											Items: []*types.UnorderedListElement{
												{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "level 2"},
															},
														},
														types.UnorderedList{
															Items: []*types.UnorderedListElement{
																{
																	BulletStyle: types.ThreeAsterisks,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		&types.Paragraph{
																			Elements: []interface{}{
																				types.StringElement{Content: "level 3"},
																				types.StringElement{
																					Content: "This is a new line inside an unordered list using ",
																				},
																				types.PredefinedAttribute{
																					Name: "plus",
																				},
																				types.StringElement{
																					Content: " symbol.",
																				},
																				types.StringElement{
																					Content: "We can even force content to start on a separate line\u2026\u200b",
																				},
																				types.LineBreak{},
																				types.StringElement{
																					Content: "Amazing, isn\u2019t it?",
																				},
																			},
																		},
																		types.UnorderedList{
																			Items: []*types.UnorderedListElement{
																				{
																					BulletStyle: types.FourAsterisks,
																					CheckStyle:  types.NoCheck,
																					Elements: []interface{}{
																						&types.Paragraph{
																							Elements: []interface{}{
																								types.StringElement{Content: "level 4"},
																							},
																						},
																						// the `+` continuation produces the second paragrap below
																						&types.Paragraph{
																							Elements: []interface{}{
																								types.StringElement{
																									Content: "The ",
																								},
																								types.PredefinedAttribute{
																									Name: "plus",
																								},
																								types.StringElement{
																									Content: " symbol is on a new line.",
																								},
																							},
																						},

																						types.UnorderedList{
																							Items: []*types.UnorderedListElement{
																								{
																									BulletStyle: types.FiveAsterisks,
																									CheckStyle:  types.NoCheck,
																									Elements: []interface{}{
																										&types.Paragraph{
																											Elements: []interface{}{
																												types.StringElement{
																													Content: "level 5",
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

			It("without item continuation", func() {
				source := `* foo
----
a delimited block
----
* bar
----
another delimited block
----`
				expected := types.Document{
					Elements: []interface{}{
						types.UnorderedList{
							Items: []*types.UnorderedListElement{
								{
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
							},
						},
						types.ListingBlock{
							Elements: []interface{}{
								types.StringElement{
									Content: "a delimited block",
								},
							},
						},
						types.UnorderedList{
							Items: []*types.UnorderedListElement{
								{
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
						types.ListingBlock{
							Elements: []interface{}{
								types.StringElement{
									Content: "another delimited block",
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})

		Context("attach to ancestor", func() {

			It("attach to grandparent item", func() {
				source := `* grandparent list item
** parent list item
*** child list item


+
paragraph attached to grandparent list item`
				expected := types.Document{
					Elements: []interface{}{
						types.UnorderedList{
							Items: []*types.UnorderedListElement{
								{
									BulletStyle: types.OneAsterisk,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												types.StringElement{Content: "grandparent list item"},
											},
										},
										types.UnorderedList{
											Items: []*types.UnorderedListElement{
												{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "parent list item"},
															},
														},
														types.UnorderedList{
															Items: []*types.UnorderedListElement{
																{
																	BulletStyle: types.ThreeAsterisks,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		&types.Paragraph{
																			Elements: []interface{}{
																				types.StringElement{Content: "child list item"},
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
										&types.Paragraph{
											Elements: []interface{}{
												types.StringElement{Content: "paragraph attached to grandparent list item"},
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

			It("attach to parent item", func() {
				source := `* grandparent list item
** parent list item
*** child list item

+
paragraph attached to parent list item`
				expected := types.Document{
					Elements: []interface{}{
						types.UnorderedList{
							Items: []*types.UnorderedListElement{
								{
									BulletStyle: types.OneAsterisk,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										&types.Paragraph{
											Elements: []interface{}{
												types.StringElement{Content: "grandparent list item"},
											},
										},
										types.UnorderedList{
											Items: []*types.UnorderedListElement{
												{
													BulletStyle: types.TwoAsterisks,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "parent list item"},
															},
														},
														types.UnorderedList{
															Items: []*types.UnorderedListElement{
																{
																	BulletStyle: types.ThreeAsterisks,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		&types.Paragraph{
																			Elements: []interface{}{
																				types.StringElement{Content: "child list item"},
																			},
																		},
																	},
																},
															},
														},
														&types.Paragraph{
															Elements: []interface{}{
																types.StringElement{Content: "paragraph attached to parent list item"},
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
		})
	})
})
