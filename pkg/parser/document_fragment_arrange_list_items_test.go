package parser

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("arrange list items", func() {

	It("mixed lists - complex case 1", func() {
		// - unordered 1
		// 1. ordered 1.1
		// 	a. ordered 1.1.a
		// 	b. ordered 1.1.b
		// 	c. ordered 1.1.c
		// 2. ordered 1.2
		// 	i)  ordered 1.2.i
		// 	ii) ordered 1.2.ii
		// 3. ordered 1.3
		// 4. ordered 1.4
		// - unordered 2
		// * unordered 2.1
		actual := &types.ListItemBucket{
			Elements: []interface{}{
				types.UnorderedListElement{
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
					},
				},
				types.OrderedListElement{
					Style: types.Arabic,
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "ordered 1.1"},
								},
							},
						},
					},
				},
				types.OrderedListElement{
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
				types.OrderedListElement{
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
				types.OrderedListElement{
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
				types.OrderedListElement{
					Style: types.Arabic,
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "ordered 1.2"},
								},
							},
						},
					},
				},
				types.OrderedListElement{
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
				types.OrderedListElement{
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
				types.OrderedListElement{
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
				types.OrderedListElement{
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
				types.UnorderedListElement{
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
					},
				},
				types.UnorderedListElement{
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
		}
		expected := []interface{}{
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
		}
		Expect(arrangeListItems(actual)).To(Equal(expected))
	})

	It("labeled list with rich terms", func() {
		actual := &types.ListItemBucket{
			Elements: []interface{}{
				types.LabeledListElement{
					Term: []interface{}{
						types.StringElement{
							Content: "`foo` term",
						},
					},
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "description 1"},
								},
							},
						},
					},
				},
				types.LabeledListElement{
					Term: []interface{}{
						types.StringElement{
							Content: "`bar` term",
						},
					},
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "description 2"},
								},
							},
						},
					},
				},
				types.LabeledListElement{
					Term: []interface{}{
						types.StringElement{
							Content: "icon:caution[]",
						},
					},
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "description 3"},
								},
							},
						},
					},
				},
			},
		}
		expected := []interface{}{
			types.LabeledList{
				Items: []*types.LabeledListElement{
					{
						Term: []interface{}{
							&types.QuotedText{
								Kind: types.SingleQuoteMonospace,
								Elements: []interface{}{
									types.StringElement{
										Content: "foo",
									},
								},
							},
							types.StringElement{
								Content: " term",
							},
						},
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "description 1"},
									},
								},
							},
							types.LabeledList{
								Items: []*types.LabeledListElement{
									{
										Term: []interface{}{
											&types.QuotedText{
												Kind: types.SingleQuoteMonospace,
												Elements: []interface{}{
													types.StringElement{
														Content: "bar",
													},
												},
											},
											types.StringElement{
												Content: " term",
											},
										},
										Elements: []interface{}{
											types.Paragraph{
												Lines: [][]interface{}{
													{
														types.StringElement{Content: "description 2"},
													},
												},
											},
										},
									},
									{
										Term: []interface{}{
											types.Icon{
												Class: "caution",
											},
										},
										Elements: []interface{}{
											types.Paragraph{
												Lines: [][]interface{}{
													{
														types.StringElement{Content: "description 3"},
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
		Expect(arrangeListItems(actual)).To(Equal(expected))
	})

	It("callout list items and a block afterwards", func() {
		actual := &types.ListItemBucket{
			Elements: []interface{}{
				types.CalloutListElement{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "description 1"},
								},
							},
						},
					},
				},
				types.CalloutListElement{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "description 2"},
								},
							},
						},
					},
				},
				types.ExampleBlock{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "foo",
									},
								},
							},
						},
					},
				},
			},
		}
		expected := []interface{}{
			types.CalloutList{
				Items: []*types.CalloutListElement{
					{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "description 1"},
									},
								},
							},
						},
					},
					{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "description 2"},
									},
								},
							},
						},
					},
				},
			},
			types.ExampleBlock{
				Elements: []interface{}{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{
									Content: "foo",
								},
							},
						},
					},
				},
			},
		}
		Expect(arrangeListItems(actual)).To(Equal(expected))
	})

	It("unordered list items and continued list item attached to grandparent", func() {
		// * grandparent list item
		// ** parent list item
		// *** child list item
		//
		//
		// +
		// paragraph attached to parent list item

		actual := &types.ListItemBucket{
			Elements: []interface{}{
				types.UnorderedListElement{
					BulletStyle: types.OneAsterisk,
					CheckStyle:  types.NoCheck,
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "grandparent list item"},
								},
							},
						},
					},
				},
				types.UnorderedListElement{
					BulletStyle: types.TwoAsterisks,
					CheckStyle:  types.NoCheck,
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "parent list item"},
								},
							},
						},
					},
				},
				types.UnorderedListElement{
					BulletStyle: types.ThreeAsterisks,
					CheckStyle:  types.NoCheck,
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "child list item"},
								},
							},
						},
					},
				},
				types.BlankLine{},
				types.BlankLine{},
				types.ContinuedListItemElement{
					Element: types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "paragraph attached to grandparent list item"},
							},
						},
					},
				},
			},
		}
		expected := []interface{}{
			types.UnorderedList{
				Items: []*types.UnorderedListElement{
					{
						BulletStyle: types.OneAsterisk,
						CheckStyle:  types.NoCheck,
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "grandparent list item"},
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
														types.StringElement{Content: "parent list item"},
													},
												},
											},
											types.UnorderedList{
												Items: []*types.UnorderedListElement{
													{
														BulletStyle: types.ThreeAsterisks,
														CheckStyle:  types.NoCheck,
														Elements: []interface{}{
															types.Paragraph{
																Lines: [][]interface{}{
																	{
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
							},
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "paragraph attached to grandparent list item"},
									},
								},
							},
						},
					},
				},
			},
		}
		Expect(arrangeListItems(actual)).To(Equal(expected))
	})

	It("unordered list items and continued list item attached to parent", func() {
		// * grandparent list item
		// ** parent list item
		// *** child list item
		//
		// +
		// paragraph attached to parent list item

		actual := &types.ListItemBucket{
			Elements: []interface{}{
				types.UnorderedListElement{
					BulletStyle: types.OneAsterisk,
					CheckStyle:  types.NoCheck,
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "grandparent list item"},
								},
							},
						},
					},
				},
				types.UnorderedListElement{
					BulletStyle: types.TwoAsterisks,
					CheckStyle:  types.NoCheck,
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "parent list item"},
								},
							},
						},
					},
				},
				types.UnorderedListElement{
					BulletStyle: types.ThreeAsterisks,
					CheckStyle:  types.NoCheck,
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "child list item"},
								},
							},
						},
					},
				},
				types.BlankLine{},
				types.ContinuedListItemElement{
					Element: types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "paragraph attached to parent list item"},
							},
						},
					},
				},
			},
		}
		expected := []interface{}{
			types.UnorderedList{
				Items: []*types.UnorderedListElement{
					{
						BulletStyle: types.OneAsterisk,
						CheckStyle:  types.NoCheck,
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "grandparent list item"},
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
														types.StringElement{Content: "parent list item"},
													},
												},
											},
											types.UnorderedList{
												Items: []*types.UnorderedListElement{
													{
														BulletStyle: types.ThreeAsterisks,
														CheckStyle:  types.NoCheck,
														Elements: []interface{}{
															types.Paragraph{
																Lines: [][]interface{}{
																	{
																		types.StringElement{Content: "child list item"},
																	},
																},
															},
														},
													},
												},
											},
											types.Paragraph{
												Lines: [][]interface{}{
													{
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
		Expect(arrangeListItems(actual)).To(Equal(expected))
	})

	It("unordered list items and continued list item attached to parent", func() {
		// * grandparent list item
		// ** parent list item
		// *** child list item
		//
		// +
		// paragraph attached to parent list item

		actual := &types.ListItemBucket{
			Elements: []interface{}{
				types.UnorderedListElement{
					BulletStyle: types.OneAsterisk,
					CheckStyle:  types.NoCheck,
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "grandparent list item"},
								},
							},
						},
					},
				},
				types.UnorderedListElement{
					BulletStyle: types.TwoAsterisks,
					CheckStyle:  types.NoCheck,
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "parent list item"},
								},
							},
						},
					},
				},
				types.UnorderedListElement{
					BulletStyle: types.ThreeAsterisks,
					CheckStyle:  types.NoCheck,
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "child list item"},
								},
							},
						},
					},
				},
				types.BlankLine{},
				types.ContinuedListItemElement{
					Element: types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "paragraph attached to parent list item"},
							},
						},
					},
				},
			},
		}
		expected := []interface{}{
			types.UnorderedList{
				Items: []*types.UnorderedListElement{
					{
						BulletStyle: types.OneAsterisk,
						CheckStyle:  types.NoCheck,
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "grandparent list item"},
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
														types.StringElement{Content: "parent list item"},
													},
												},
											},
											types.UnorderedList{
												Items: []*types.UnorderedListElement{
													{
														BulletStyle: types.ThreeAsterisks,
														CheckStyle:  types.NoCheck,
														Elements: []interface{}{
															types.Paragraph{
																Lines: [][]interface{}{
																	{
																		types.StringElement{Content: "child list item"},
																	},
																},
															},
														},
													},
												},
											},
											types.Paragraph{
												Lines: [][]interface{}{
													{
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
		Expect(arrangeListItems(actual)).To(Equal(expected))
	})

	It("unordered list items and continued list item attached to parent and grandparent", func() {
		// * grandparent list item
		// ** parent list item
		// *** child list item
		//
		// +
		// paragraph attached to parent list item
		//
		// +
		// paragraph attached to grandparent list item

		actual := &types.ListItemBucket{
			Elements: []interface{}{
				types.UnorderedListElement{
					BulletStyle: types.OneAsterisk,
					CheckStyle:  types.NoCheck,
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "grandparent list item"},
								},
							},
						},
					},
				},
				types.UnorderedListElement{
					BulletStyle: types.TwoAsterisks,
					CheckStyle:  types.NoCheck,
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "parent list item"},
								},
							},
						},
					},
				},
				types.UnorderedListElement{
					BulletStyle: types.ThreeAsterisks,
					CheckStyle:  types.NoCheck,
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{Content: "child list item"},
								},
							},
						},
					},
				},
				types.BlankLine{},
				types.ContinuedListItemElement{
					Element: types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "paragraph attached to parent list item"},
							},
						},
					},
				},
				types.BlankLine{},
				types.ContinuedListItemElement{
					Element: types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "paragraph attached to grandparent list item"},
							},
						},
					},
				},
			},
		}
		expected := []interface{}{
			types.UnorderedList{
				Items: []*types.UnorderedListElement{
					{
						BulletStyle: types.OneAsterisk,
						CheckStyle:  types.NoCheck,
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "grandparent list item"},
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
														types.StringElement{Content: "parent list item"},
													},
												},
											},
											types.UnorderedList{
												Items: []*types.UnorderedListElement{
													{
														BulletStyle: types.ThreeAsterisks,
														CheckStyle:  types.NoCheck,
														Elements: []interface{}{
															types.Paragraph{
																Lines: [][]interface{}{
																	{
																		types.StringElement{Content: "child list item"},
																	},
																},
															},
														},
													},
												},
											},
											types.Paragraph{
												Lines: [][]interface{}{
													{
														types.StringElement{Content: "paragraph attached to parent list item"},
													},
												},
											},
										},
									},
								},
							},
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "paragraph attached to grandparent list item"},
									},
								},
							},
						},
					},
				},
			},
		}
		Expect(arrangeListItems(actual)).To(Equal(expected))
	})

})
