package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("mixed lists", func() {

	Context("final documents", func() {

		Context("valid mixed lists", func() {

			It("ordered list with nested unordered lists", func() {
				source := `. Item 1
* Item A
* Item B
. Item 2
* Item C
* Item D`
				expected := types.Document{
					Elements: []interface{}{
						types.OrderedList{
							Items: []*types.OrderedListElement{
								{
									Style: types.Arabic,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{Content: "Item 1"},
												},
											},
										},
										types.UnorderedList{
											Items: []*types.UnorderedListElement{
												{
													BulletStyle: types.OneAsterisk,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "Item A"},
																},
															},
														},
													},
												},
												{
													BulletStyle: types.OneAsterisk,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "Item B"},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								{
									Style: types.Arabic,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{Content: "Item 2"},
												},
											},
										},
										types.UnorderedList{
											Items: []*types.UnorderedListElement{
												{
													BulletStyle: types.OneAsterisk,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "Item C"},
																},
															},
														},
													},
												},
												{
													BulletStyle: types.OneAsterisk,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "Item D"},
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

			It("unordered list item and order list item with roman numbering", func() {
				source := `- unordered list item
 II) ordered list item`
				expected := types.Document{
					Elements: []interface{}{
						types.UnorderedList{
							Items: []*types.UnorderedListElement{
								{
									BulletStyle: types.Dash,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "unordered list item",
													},
												},
											},
										},
										types.OrderedList{
											Items: []*types.OrderedListElement{
												{
													Style: types.UpperRoman,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{
																		Content: "ordered list item",
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
		})

		Context("complex cases", func() {

			It("complex case 1 - mixed lists", func() {
				source := `- unordered 1
	1. ordered 1.1
		a. ordered 1.1.a
		b. ordered 1.1.b
		c. ordered 1.1.c
	2. ordered 1.2
		i)  ordered 1.2.i
		ii) ordered 1.2.ii
	3. ordered 1.3
	4. ordered 1.4
	- unordered 2
	* unordered 2.1`
				expected := types.Document{
					Elements: []interface{}{
						types.UnorderedList{
							Items: []*types.UnorderedListElement{
								{
									BulletStyle: types.Dash,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{Content: "unordered 1"},
												},
											},
										},
										types.OrderedList{
											Items: []*types.OrderedListElement{
												{
													Style: types.Arabic,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "ordered 1.1"},
																},
															},
														},
														types.OrderedList{
															Items: []*types.OrderedListElement{
																{
																	Style: types.LowerAlpha,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{Content: "ordered 1.1.a"},
																				},
																			},
																		},
																	},
																},
																{
																	Style: types.LowerAlpha,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{Content: "ordered 1.1.b"},
																				},
																			},
																		},
																	},
																},
																{
																	Style: types.LowerAlpha,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{Content: "ordered 1.1.c"},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
												{
													Style: types.Arabic,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "ordered 1.2"},
																},
															},
														},
														types.OrderedList{
															Items: []*types.OrderedListElement{
																{
																	Style: types.LowerRoman,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{Content: "ordered 1.2.i"},
																				},
																			},
																		},
																	},
																},
																{
																	Style: types.LowerRoman,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{Content: "ordered 1.2.ii"},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
												{
													Style: types.Arabic,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "ordered 1.3"},
																},
															},
														},
													},
												},
												{
													Style: types.Arabic,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "ordered 1.4"},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								{
									BulletStyle: types.Dash,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{Content: "unordered 2"},
												},
											},
										},
										types.UnorderedList{
											Items: []*types.UnorderedListElement{
												{
													BulletStyle: types.OneAsterisk,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "unordered 2.1"},
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

			It("complex case 2 - mixed lists", func() {
				source := `- unordered 1
1. ordered 1.1
a. ordered 1.1.a
b. ordered 1.1.b
c. ordered 1.1.c
2. ordered 1.2
i)  ordered 1.2.i
ii) ordered 1.2.ii
3. ordered 1.3
4. ordered 1.4
- unordered 2
* unordered 2.1
** unordered 2.1.1
	with some
	extra lines.
** unordered 2.1.2
* unordered 2.2
- unordered 3
. ordered 3.1
. ordered 3.2
[upperroman]
.. ordered 3.2.I
.. ordered 3.2.II
. ordered 3.3`
				expected := types.Document{
					Elements: []interface{}{
						types.UnorderedList{
							Items: []*types.UnorderedListElement{
								{
									BulletStyle: types.Dash,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{Content: "unordered 1"},
												},
											},
										},
										types.OrderedList{
											Items: []*types.OrderedListElement{
												{
													Style: types.Arabic,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "ordered 1.1"},
																},
															},
														},
														types.OrderedList{
															Items: []*types.OrderedListElement{
																{
																	Style: types.LowerAlpha,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{Content: "ordered 1.1.a"},
																				},
																			},
																		},
																	},
																},
																{
																	Style: types.LowerAlpha,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{Content: "ordered 1.1.b"},
																				},
																			},
																		},
																	},
																},
																{
																	Style: types.LowerAlpha,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{Content: "ordered 1.1.c"},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
												{
													Style: types.Arabic,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "ordered 1.2"},
																},
															},
														},
														types.OrderedList{
															Items: []*types.OrderedListElement{
																{
																	Style: types.LowerRoman,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{Content: "ordered 1.2.i"},
																				},
																			},
																		},
																	},
																},
																{
																	Style: types.LowerRoman,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{Content: "ordered 1.2.ii"},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
												{
													Style: types.Arabic,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "ordered 1.3"},
																},
															},
														},
													},
												},
												{
													Style: types.Arabic,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "ordered 1.4"},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								{
									BulletStyle: types.Dash,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{Content: "unordered 2"},
												},
											},
										},
										types.UnorderedList{
											Items: []*types.UnorderedListElement{
												{
													BulletStyle: types.OneAsterisk,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "unordered 2.1"},
																},
															},
														},
														types.UnorderedList{
															Items: []*types.UnorderedListElement{
																{
																	BulletStyle: types.TwoAsterisks,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{Content: "unordered 2.1.1"},
																				},
																				{
																					types.StringElement{Content: "with some"}, // heading tabs are trimmed
																				},
																				{
																					types.StringElement{Content: "extra lines."}, // heading tabs are trimmed
																				},
																			},
																		},
																	},
																},
																{
																	BulletStyle: types.TwoAsterisks,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{Content: "unordered 2.1.2"},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
												{
													BulletStyle: types.OneAsterisk,
													CheckStyle:  types.NoCheck,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "unordered 2.2"},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								{
									BulletStyle: types.Dash,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{Content: "unordered 3"},
												},
											},
										},
										types.OrderedList{
											Items: []*types.OrderedListElement{
												{
													Style: types.Arabic,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "ordered 3.1"},
																},
															},
														},
													},
												},
												{
													Style: types.Arabic,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "ordered 3.2"},
																},
															},
														},
														types.OrderedList{
															Attributes: types.Attributes{
																types.AttrStyle: "upperroman",
															},
															Items: []*types.OrderedListElement{
																{
																	Style: types.LowerAlpha,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{Content: "ordered 3.2.I"},
																				},
																			},
																		},
																	},
																},
																{
																	Style: types.LowerAlpha,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{Content: "ordered 3.2.II"},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
												{
													Style: types.Arabic,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{Content: "ordered 3.3"},
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

			It("complex case 4 - mixed lists", func() {
				source := `.Mixed
Operating Systems::
  . Fedora
    * Desktop`
				expected := types.Document{
					Elements: []interface{}{
						types.LabeledList{
							Attributes: types.Attributes{
								types.AttrTitle: "Mixed",
							},
							Items: []*types.LabeledListElement{
								{
									Term: []interface{}{
										types.StringElement{
											Content: "Operating Systems",
										},
									},
									Elements: []interface{}{
										types.OrderedList{
											Items: []*types.OrderedListElement{
												{
													Style: types.Arabic,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{
																		Content: "Fedora",
																	},
																},
															},
														},
														types.UnorderedList{
															Items: []*types.UnorderedListElement{
																{
																	BulletStyle: types.OneAsterisk,
																	CheckStyle:  types.NoCheck,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{
																						Content: "Desktop",
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

			It("complex case 5 - mixed lists and a paragraph", func() {
				source := `.Mixed
Operating Systems::
  Linux:::
    . Fedora
      * Desktop
    . Ubuntu
      * Desktop
      * Server
  BSD:::
    . FreeBSD
    . NetBSD

Cloud Providers::
  PaaS:::
    . OpenShift
    . CloudBees
  IaaS:::
    . Amazon EC2
	. Rackspace
	
a paragraph
`
				expected := types.Document{
					Elements: []interface{}{
						types.LabeledList{
							Attributes: types.Attributes{
								types.AttrTitle: "Mixed",
							},
							Items: []*types.LabeledListElement{
								{
									Term: []interface{}{
										types.StringElement{
											Content: "Operating Systems",
										},
									},
									Elements: []interface{}{
										types.LabeledList{
											Items: []*types.LabeledListElement{
												{
													Term: []interface{}{
														types.StringElement{
															Content: "Linux",
														},
													},
													Elements: []interface{}{
														types.OrderedList{
															Items: []*types.OrderedListElement{
																{
																	Style: types.Arabic,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{
																						Content: "Fedora",
																					},
																				},
																			},
																		},
																		types.UnorderedList{
																			Items: []*types.UnorderedListElement{
																				{
																					BulletStyle: types.OneAsterisk,
																					CheckStyle:  types.NoCheck,
																					Elements: []interface{}{
																						types.Paragraph{
																							Lines: [][]interface{}{
																								{
																									types.StringElement{
																										Content: "Desktop",
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
																{
																	Style: types.Arabic,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{
																						Content: "Ubuntu",
																					},
																				},
																			},
																		},
																		types.UnorderedList{
																			Items: []*types.UnorderedListElement{
																				{
																					BulletStyle: types.OneAsterisk,
																					CheckStyle:  types.NoCheck,
																					Elements: []interface{}{
																						types.Paragraph{
																							Lines: [][]interface{}{
																								{
																									types.StringElement{
																										Content: "Desktop",
																									},
																								},
																							},
																						},
																					},
																				},
																				{
																					BulletStyle: types.OneAsterisk,
																					CheckStyle:  types.NoCheck,
																					Elements: []interface{}{
																						types.Paragraph{
																							Lines: [][]interface{}{
																								{
																									types.StringElement{
																										Content: "Server",
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
												{
													Term: []interface{}{
														types.StringElement{
															Content: "BSD",
														},
													},
													Elements: []interface{}{
														types.OrderedList{
															Items: []*types.OrderedListElement{
																{
																	Style: types.Arabic,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{
																						Content: "FreeBSD",
																					},
																				},
																			},
																		},
																	},
																},
																{
																	Style: types.Arabic,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{
																						Content: "NetBSD",
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
								{
									Term: []interface{}{
										types.StringElement{
											Content: "Cloud Providers",
										},
									},
									Elements: []interface{}{
										types.LabeledList{
											Items: []*types.LabeledListElement{
												{
													Term: []interface{}{
														types.StringElement{
															Content: "PaaS",
														},
													},
													Elements: []interface{}{
														types.OrderedList{
															Items: []*types.OrderedListElement{
																{
																	Style: types.Arabic,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{
																						Content: "OpenShift",
																					},
																				},
																			},
																		},
																	},
																},
																{
																	Style: types.Arabic,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{
																						Content: "CloudBees",
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
												{
													Term: []interface{}{
														types.StringElement{
															Content: "IaaS",
														},
													},
													Elements: []interface{}{
														types.OrderedList{
															Items: []*types.OrderedListElement{
																{
																	Style: types.Arabic,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{
																						Content: "Amazon EC2",
																					},
																				},
																			},
																		},
																	},
																},
																{
																	Style: types.Arabic,
																	Elements: []interface{}{
																		types.Paragraph{
																			Lines: [][]interface{}{
																				{
																					types.StringElement{
																						Content: "Rackspace",
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
						// types.BlankLine{},
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "a paragraph",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})

		Context("distinct list blocks", func() {

			It("same list without attributes", func() {
				source := `[lowerroman, start=5]
	. Five
	.. a
	. Six`
				expected := types.Document{
					Elements: []interface{}{ // a single ordered list
						types.OrderedList{
							Attributes: types.Attributes{
								types.AttrStyle: "lowerroman",
								types.AttrStart: "5",
							},
							Items: []*types.OrderedListElement{
								{
									Style: types.Arabic, // will be overridden during rendering
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "Five",
													},
												},
											},
										},
										types.OrderedList{
											Items: []*types.OrderedListElement{
												{
													Style: types.LowerAlpha,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{
																		Content: "a",
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
								{
									Style: types.Arabic, // will be overridden during rendering
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "Six",
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

			It("same list with custom number style on sublist", func() {
				// need to be aligned on first column of file
				source := `[lowerroman, start=5]
. Five
[upperalpha]
.. a
.. b
. Six`
				expected := types.Document{
					Elements: []interface{}{ // a single ordered list
						types.OrderedList{
							Attributes: types.Attributes{
								types.AttrStyle: "lowerroman",
								types.AttrStart: "5",
							},
							Items: []*types.OrderedListElement{
								{
									Style: types.Arabic, // will be overridden during rendering
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "Five",
													},
												},
											},
										},
										types.OrderedList{
											Attributes: types.Attributes{
												types.AttrStyle: "upperalpha",
											},
											Items: []*types.OrderedListElement{
												{
													Style: types.LowerAlpha, // will be overridden during rendering
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{
																		Content: "a",
																	},
																},
															},
														},
													},
												},
												{
													Style: types.LowerAlpha, // will be overridden during rendering
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{
																		Content: "b",
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
								{
									Style: types.Arabic, // will be overridden during rendering
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "Six",
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

			It("distinct lists with blankline and item attribute - case 1", func() {
				source := `[lowerroman, start=5]
. Five

[upperalpha]
.. a
. Six`
				expected := types.Document{
					Elements: []interface{}{ // a single ordered list
						types.OrderedList{
							Attributes: types.Attributes{
								types.AttrStyle: "lowerroman",
								types.AttrStart: "5",
							},
							Items: []*types.OrderedListElement{
								{
									Style: types.Arabic,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "Five",
													},
												},
											},
										},
									},
								},
							},
						},
						types.OrderedList{
							Attributes: types.Attributes{
								types.AttrStyle: "upperalpha",
							},
							Items: []*types.OrderedListElement{
								{
									Style: types.LowerAlpha,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "a",
													},
												},
											},
										},
										types.OrderedList{
											Items: []*types.OrderedListElement{
												{
													Style: types.Arabic,
													Elements: []interface{}{
														types.Paragraph{
															Lines: [][]interface{}{
																{
																	types.StringElement{
																		Content: "Six",
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

			It("distinct lists with blankline and item attribute - case 2", func() {

				source := `.Checklist
- [*] checked
-     normal list item

.Ordered, basic
. Step 1
. Step 2`
				expected := types.Document{
					Elements: []interface{}{
						types.UnorderedList{
							Attributes: types.Attributes{
								types.AttrTitle: "Checklist",
							},
							Items: []*types.UnorderedListElement{
								{
									BulletStyle: types.Dash,
									CheckStyle:  types.Checked,
									Elements: []interface{}{
										types.Paragraph{
											Attributes: types.Attributes{
												types.AttrCheckStyle: types.Checked,
											},
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "checked",
													},
												},
											},
										},
									},
								},
								{
									BulletStyle: types.Dash,
									CheckStyle:  types.NoCheck,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "normal list item",
													},
												},
											},
										},
									},
								},
							},
						},
						types.OrderedList{
							Attributes: types.Attributes{
								types.AttrTitle: "Ordered, basic",
							},
							Items: []*types.OrderedListElement{
								{
									Style: types.Arabic,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "Step 1",
													},
												},
											},
										},
									},
								},
								{
									Style: types.Arabic,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "Step 2",
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

			It("same list with single comment line inside", func() {
				source := `. a
// -
. b`
				expected := types.Document{
					Elements: []interface{}{
						types.OrderedList{
							Items: []*types.OrderedListElement{
								{
									Style: types.Arabic,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "a",
													},
												},
											},
										},
									},
								},
								{
									Style: types.Arabic,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "b",
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

			It("same list with multiple comment lines inside", func() {
				source := `. a
// -
// -
// -
. b`
				expected := types.Document{
					Elements: []interface{}{
						types.OrderedList{
							Items: []*types.OrderedListElement{
								{
									Style: types.Arabic,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "a",
													},
												},
											},
										},
									},
								},
								{
									Style: types.Arabic,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "b",
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

			It("distinct lists separated by single comment line", func() {
				source := `. a
	
// -
. b`
				expected := types.Document{
					Elements: []interface{}{
						types.OrderedList{
							Items: []*types.OrderedListElement{
								{
									Style: types.Arabic,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "a",
													},
												},
											},
										},
									},
								},
							},
						},
						types.OrderedList{
							Items: []*types.OrderedListElement{
								{
									Style: types.Arabic,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "b",
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

			It("distinct lists separated by multiple comment lines", func() {
				source := `. a
	
// -
// -
// -
. b`
				expected := types.Document{
					Elements: []interface{}{
						types.OrderedList{
							Items: []*types.OrderedListElement{
								{
									Style: types.Arabic,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "a",
													},
												},
											},
										},
									},
								},
							},
						},
						types.OrderedList{
							Items: []*types.OrderedListElement{
								{
									Style: types.Arabic,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "b",
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
