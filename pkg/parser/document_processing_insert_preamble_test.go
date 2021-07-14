package parser

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("preambles", func() {

	sectionATitle := []interface{}{
		&types.StringElement{Content: "Section A"},
	}

	sectionBTitle := []interface{}{
		&types.StringElement{Content: "Section B"},
	}

	It("doc without sections", func() {
		// given
		doc := &types.Document{
			Attributes: types.Attributes{
				types.AttrTitle: "foo",
			},
			ElementReferences: types.ElementReferences{
				"_section_a": sectionATitle,
				"_section_b": sectionBTitle,
			},
			Elements: []interface{}{
				types.Paragraph{
					Lines: [][]interface{}{
						{
							&types.StringElement{Content: "a short paragraph"},
						},
					},
				},
				&types.BlankLine{},
				types.Paragraph{
					Lines: [][]interface{}{
						{
							&types.StringElement{Content: "another short paragraph"},
						},
					},
				},
			},
		}
		expected := &types.Document{
			Attributes: types.Attributes{
				types.AttrTitle: "foo",
			},
			ElementReferences: types.ElementReferences{
				"_section_a": sectionATitle,
				"_section_b": sectionBTitle,
			},
			Elements: []interface{}{
				types.Paragraph{
					Lines: [][]interface{}{
						{
							&types.StringElement{Content: "a short paragraph"},
						},
					},
				},
				&types.BlankLine{},
				types.Paragraph{
					Lines: [][]interface{}{
						{
							&types.StringElement{Content: "another short paragraph"},
						},
					},
				},
			},
		}
		// when
		doc.InsertPreamble()
		// then
		Expect(doc).To(Equal(expected))
	})

	It("doc with 1-paragraph preamble", func() {
		// given
		doc := &types.Document{
			Attributes: types.Attributes{
				types.AttrTitle: "foo",
			},
			ElementReferences: types.ElementReferences{
				"_section_a": sectionATitle,
				"_section_b": sectionBTitle,
			},
			Elements: []interface{}{
				types.Paragraph{
					Lines: [][]interface{}{
						{
							&types.StringElement{Content: "a short paragraph"},
						},
					},
				},
				&types.BlankLine{},
				types.Section{
					Level: 1,
					Title: sectionATitle,
					Attributes: types.Attributes{
						types.AttrID: "_section_a",
					},

					Elements: []interface{}{},
				},
				types.Section{
					Level: 1,
					Title: sectionBTitle,
					Attributes: types.Attributes{
						types.AttrID: "_section_b",
					},
					Elements: []interface{}{},
				},
			},
		}
		expected := &types.Document{
			Attributes: types.Attributes{
				types.AttrTitle: "foo",
			},
			ElementReferences: types.ElementReferences{
				"_section_a": sectionATitle,
				"_section_b": sectionBTitle,
			},
			Elements: []interface{}{
				&types.Preamble{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a short paragraph"}},
							},
						},
						&types.BlankLine{},
					},
				},
				types.Section{
					Level: 1,
					Title: sectionATitle,
					Attributes: types.Attributes{
						types.AttrID: "_section_a",
					},

					Elements: []interface{}{},
				},
				types.Section{
					Level: 1,
					Title: sectionBTitle,
					Attributes: types.Attributes{
						types.AttrID: "_section_b",
					},
					Elements: []interface{}{},
				},
			},
		}
		// when
		doc.InsertPreamble()
		// then
		Expect(doc).To(Equal(expected))
	})

	It("doc with 2-paragraph preamble", func() {
		// given
		doc := &types.Document{
			Attributes: types.Attributes{
				types.AttrTitle: "foo",
			},
			ElementReferences: types.ElementReferences{
				"_section_a": sectionATitle,
				"_section_b": sectionBTitle,
			},
			Elements: []interface{}{
				types.Paragraph{
					Lines: [][]interface{}{
						{
							&types.StringElement{Content: "a short paragraph"},
						},
					},
				},
				&types.BlankLine{},
				types.Paragraph{
					Lines: [][]interface{}{
						{
							&types.StringElement{Content: "another short paragraph"},
						},
					},
				},
				&types.BlankLine{},
				types.Section{
					Level: 1,
					Title: sectionATitle,
					Attributes: types.Attributes{
						types.AttrID: "_section_a",
					},

					Elements: []interface{}{},
				},
				types.Section{
					Level: 1,
					Title: sectionBTitle,
					Attributes: types.Attributes{
						types.AttrID: "_section_b",
					},
					Elements: []interface{}{},
				},
			},
		}
		expected := &types.Document{
			Attributes: types.Attributes{
				types.AttrTitle: "foo",
			},
			ElementReferences: types.ElementReferences{
				"_section_a": sectionATitle,
				"_section_b": sectionBTitle,
			},
			Elements: []interface{}{
				&types.Preamble{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a short paragraph"}},
							},
						},
						&types.BlankLine{},
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "another short paragraph"}},
							},
						},
						&types.BlankLine{},
					},
				},
				types.Section{
					Level: 1,
					Title: sectionATitle,
					Attributes: types.Attributes{
						types.AttrID: "_section_a",
					},

					Elements: []interface{}{},
				},
				types.Section{
					Level: 1,
					Title: sectionBTitle,
					Attributes: types.Attributes{
						types.AttrID: "_section_b",
					},
					Elements: []interface{}{},
				},
			},
		}
		// when
		doc.InsertPreamble()
		// then
		Expect(doc).To(Equal(expected))
	})

})
