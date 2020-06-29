package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("tables", func() {

	It("1-line table with 2 cells", func() {
		source := `|===
| *foo* foo  | _bar_  
|===
`
		expected := types.Table{
			Lines: []types.TableLine{
				{
					Cells: [][]interface{}{
						{
							types.QuotedText{
								Kind: types.Bold,
								Elements: []interface{}{
									types.StringElement{
										Content: "foo",
									},
								},
							},
							types.StringElement{
								Content: " foo  ",
							},
						},
						{
							types.QuotedText{
								Kind: types.Italic,
								Elements: []interface{}{
									types.StringElement{
										Content: "bar",
									},
								},
							},
							types.StringElement{
								Content: "  ",
							},
						},
					},
				},
			},
		}
		Expect(ParseDocumentBlock(source)).To(Equal(expected))
	})

	It("1-line table with 3 cells", func() {
		source := `|===
| *foo* foo  | _bar_  | baz
|===`
		expected := types.Table{
			Lines: []types.TableLine{
				{
					Cells: [][]interface{}{
						{
							types.QuotedText{
								Kind: types.Bold,
								Elements: []interface{}{
									types.StringElement{
										Content: "foo",
									},
								},
							},
							types.StringElement{
								Content: " foo  ",
							},
						},
						{
							types.QuotedText{
								Kind: types.Italic,
								Elements: []interface{}{
									types.StringElement{
										Content: "bar",
									},
								},
							},
							types.StringElement{
								Content: "  ",
							},
						},
						{
							types.StringElement{
								Content: "baz",
							},
						},
					},
				},
			},
		}
		Expect(ParseDocumentBlock(source)).To(Equal(expected))
	})

	It("table with title, headers and 1 line per cell", func() {
		source := `.table title
|===
|heading 1 |heading 2

|row 1, column 1
|row 1, column 2

|row 2, column 1
|row 2, column 2
|===`
		expected := types.Table{
			Attributes: types.Attributes{
				types.AttrTitle: "table title",
			},
			Header: types.TableLine{
				Cells: [][]interface{}{
					{
						types.StringElement{
							Content: "heading 1 ",
						},
					},
					{
						types.StringElement{
							Content: "heading 2",
						},
					},
				},
			},

			Lines: []types.TableLine{
				{
					Cells: [][]interface{}{
						{
							types.StringElement{
								Content: "row 1, column 1",
							},
						},
						{
							types.StringElement{
								Content: "row 1, column 2",
							},
						},
					},
				},
				{
					Cells: [][]interface{}{
						{
							types.StringElement{
								Content: "row 2, column 1",
							},
						},
						{
							types.StringElement{
								Content: "row 2, column 2",
							},
						},
					},
				},
			},
		}
		Expect(ParseDocumentBlock(source)).To(Equal(expected))
	})

	It("table with title, headers, id and multiple roles, stretch", func() {
		source := `.table title
[#anchor.role1%autowidth.stretch]
|===
|heading 1 |heading 2

|row 1, column 1
|row 1, column 2

|row 2, column 1
|row 2, column 2
|===`
		expected := types.Table{
			Attributes: types.Attributes{
				types.AttrTitle:    "table title",
				types.AttrOptions:  map[string]bool{"autowidth": true},
				types.AttrRole:     []string{"role1", "stretch"},
				types.AttrID:       "anchor",
				types.AttrCustomID: true,
			},
			Header: types.TableLine{
				Cells: [][]interface{}{
					{
						types.StringElement{
							Content: "heading 1 ",
						},
					},
					{
						types.StringElement{
							Content: "heading 2",
						},
					},
				},
			},

			Lines: []types.TableLine{
				{
					Cells: [][]interface{}{
						{
							types.StringElement{
								Content: "row 1, column 1",
							},
						},
						{
							types.StringElement{
								Content: "row 1, column 2",
							},
						},
					},
				},
				{
					Cells: [][]interface{}{
						{
							types.StringElement{
								Content: "row 2, column 1",
							},
						},
						{
							types.StringElement{
								Content: "row 2, column 2",
							},
						},
					},
				},
			},
		}
		Expect(ParseDocumentBlock(source)).To(Equal(expected))
	})

	It("empty table ", func() {
		source := `|===
|===`
		expected := types.Table{
			Lines: []types.TableLine{},
		}
		Expect(ParseDocumentBlock(source)).To(Equal(expected))
	})
})
