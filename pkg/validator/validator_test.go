package validator

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("document validator", func() {

	Context("article", func() {

		It("should not report problems", func() {
			// given
			doc := types.Document{
				Attributes:        types.DocumentAttributes{},
				ElementReferences: types.ElementReferences{},
				Footnotes:         []types.Footnote{},
				Elements: []interface{}{
					types.Section{
						Attributes: types.ElementAttributes{},
						Level:      0,
						Title: []interface{}{
							types.StringElement{
								Content: "foo",
							},
						},
					},
				},
			}

			// when
			problems := Validate(&doc)

			// then
			Expect(problems).To(BeEmpty()) // no problem found
		})
	})

	Context("manpage", func() {

		It("should not report problems", func() {
			// given
			doc := types.Document{
				Attributes: types.DocumentAttributes{
					types.AttrDocType: "manpage",
				},
				ElementReferences: types.ElementReferences{},
				Footnotes:         []types.Footnote{},
				Elements: []interface{}{
					types.Section{
						Attributes: types.ElementAttributes{},
						Level:      0,
						Title: []interface{}{
							types.StringElement{
								Content: "foo",
							},
						},
						Elements: []interface{}{
							types.Section{
								Attributes: types.ElementAttributes{},
								Level:      1,
								Title: []interface{}{
									types.StringElement{
										Content: "Name",
									},
								},
								Elements: []interface{}{
									types.Paragraph{
										Attributes: types.ElementAttributes{},
										Lines: [][]interface{}{
											{
												types.StringElement{
													Content: "a single paragraph to describe the program",
												},
											},
										},
									},
								},
							},
							types.Section{
								Attributes: types.ElementAttributes{},
								Level:      1,
								Title: []interface{}{
									types.StringElement{
										Content: "Synopsis",
									},
								},
								Elements: []interface{}{},
							},
						},
					},
				},
			}

			// when
			problems := Validate(&doc)

			// then
			Expect(problems).To(BeEmpty())                                                            // no problem found
			Expect(doc.Attributes.GetAsStringWithDefault(types.AttrDocType, "")).To(Equal("manpage")) // unchanged
		})

		Context("should report problems", func() {

			It("missing header - invalid level", func() {
				// given
				doc := types.Document{
					Attributes: types.DocumentAttributes{
						types.AttrDocType: "manpage",
					},
					ElementReferences: types.ElementReferences{},
					Footnotes:         []types.Footnote{},
					Elements: []interface{}{
						types.Section{
							Attributes: types.ElementAttributes{},
							Level:      1, // invalid level
							Title: []interface{}{
								types.StringElement{
									Content: "foo",
								},
							},
							Elements: []interface{}{
								types.Section{
									Attributes: types.ElementAttributes{},
									Level:      1,
									Title: []interface{}{
										types.StringElement{
											Content: "Name",
										},
									},
									Elements: []interface{}{
										types.Paragraph{
											Attributes: types.ElementAttributes{},
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "a single paragraph to describe the program",
													},
												},
											},
										},
									},
								},
								types.Section{
									Attributes: types.ElementAttributes{},
									Level:      1,
									Title: []interface{}{
										types.StringElement{
											Content: "Synopsis",
										},
									},
									Elements: []interface{}{},
								},
							},
						},
					},
				}

				// when
				problems := Validate(&doc)

				// then
				Expect(problems).To(ContainElement(Problem{
					Severity: Error,
					Message:  "manpage document is missing a header",
				}))
				Expect(doc.Attributes.GetAsStringWithDefault(types.AttrDocType, "")).To(Equal("article")) // changed
			})

			It("missing name section - invalid level", func() {
				// given
				doc := types.Document{
					Attributes: types.DocumentAttributes{
						types.AttrDocType: "manpage",
					},
					ElementReferences: types.ElementReferences{},
					Footnotes:         []types.Footnote{},
					Elements: []interface{}{
						types.Section{
							Attributes: types.ElementAttributes{},
							Level:      0,
							Title: []interface{}{
								types.StringElement{
									Content: "foo",
								},
							},
							Elements: []interface{}{
								types.Section{
									Attributes: types.ElementAttributes{},
									Level:      2, // invalid level
									Title: []interface{}{
										types.StringElement{
											Content: "Name",
										},
									},
									Elements: []interface{}{
										types.Paragraph{
											Attributes: types.ElementAttributes{},
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "a single paragraph to describe the program",
													},
												},
											},
										},
									},
								},
								types.Section{
									Attributes: types.ElementAttributes{},
									Level:      1,
									Title: []interface{}{
										types.StringElement{
											Content: "Synopsis",
										},
									},
									Elements: []interface{}{},
								},
							},
						},
					},
				}

				// when
				problems := Validate(&doc)

				// then
				Expect(problems).To(ContainElement(Problem{
					Severity: Error,
					Message:  "manpage document is missing the 'Name' section'",
				}))
				Expect(doc.Attributes.GetAsStringWithDefault(types.AttrDocType, "")).To(Equal("article")) // changed
			})

			It("missing name section - invalid title", func() {
				// given
				doc := types.Document{
					Attributes: types.DocumentAttributes{
						types.AttrDocType: "manpage",
					},
					ElementReferences: types.ElementReferences{},
					Footnotes:         []types.Footnote{},
					Elements: []interface{}{
						types.Section{
							Attributes: types.ElementAttributes{},
							Level:      0,
							Title: []interface{}{
								types.StringElement{
									Content: "foo",
								},
							},
							Elements: []interface{}{
								types.Section{
									Attributes: types.ElementAttributes{},
									Level:      1,
									Title: []interface{}{
										types.StringElement{
											Content: "bar", // invalid title
										},
									},
									Elements: []interface{}{
										types.Paragraph{
											Attributes: types.ElementAttributes{},
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "a single paragraph to describe the program",
													},
												},
											},
										},
									},
								},
								types.Section{
									Attributes: types.ElementAttributes{},
									Level:      1,
									Title: []interface{}{
										types.StringElement{
											Content: "Synopsis",
										},
									},
									Elements: []interface{}{},
								},
							},
						},
					},
				}

				// when
				problems := Validate(&doc)

				// then
				Expect(problems).To(ContainElement(Problem{
					Severity: Error,
					Message:  "manpage document is missing the 'Name' section'",
				}))
				Expect(doc.Attributes.GetAsStringWithDefault(types.AttrDocType, "")).To(Equal("article")) // changed
			})

			It("missing name section - invalid elements", func() {
				// given
				doc := types.Document{
					Attributes: types.DocumentAttributes{
						types.AttrDocType: "manpage",
					},
					ElementReferences: types.ElementReferences{},
					Footnotes:         []types.Footnote{},
					Elements: []interface{}{
						types.Section{
							Attributes: types.ElementAttributes{},
							Level:      0,
							Title: []interface{}{
								types.StringElement{
									Content: "foo",
								},
							},
							Elements: []interface{}{
								types.Section{
									Attributes: types.ElementAttributes{},
									Level:      1,
									Title: []interface{}{
										types.StringElement{
											Content: "Name",
										},
									},
									Elements: []interface{}{}, // invalid length
								},
								types.Section{
									Attributes: types.ElementAttributes{},
									Level:      1,
									Title: []interface{}{
										types.StringElement{
											Content: "Synopsis",
										},
									},
									Elements: []interface{}{},
								},
							},
						},
					},
				}

				// when
				problems := Validate(&doc)

				// then
				Expect(problems).To(ContainElement(Problem{
					Severity: Error,
					Message:  "'Name' section' should contain a single paragraph",
				}))
				Expect(doc.Attributes.GetAsStringWithDefault(types.AttrDocType, "")).To(Equal("article")) // changed
			})

			It("missing synopsis section - invalid level", func() {
				// given
				doc := types.Document{
					Attributes: types.DocumentAttributes{
						types.AttrDocType: "manpage",
					},
					ElementReferences: types.ElementReferences{},
					Footnotes:         []types.Footnote{},
					Elements: []interface{}{
						types.Section{
							Attributes: types.ElementAttributes{},
							Level:      0,
							Title: []interface{}{
								types.StringElement{
									Content: "foo",
								},
							},
							Elements: []interface{}{
								types.Section{
									Attributes: types.ElementAttributes{},
									Level:      1,
									Title: []interface{}{
										types.StringElement{
											Content: "Name",
										},
									},
									Elements: []interface{}{
										types.Paragraph{
											Attributes: types.ElementAttributes{},
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "a single paragraph to describe the program",
													},
												},
											},
										},
									},
								},
								types.Section{
									Attributes: types.ElementAttributes{},
									Level:      2, // invalid level
									Title: []interface{}{
										types.StringElement{
											Content: "Synopsis",
										},
									},
									Elements: []interface{}{},
								},
							},
						},
					},
				}

				// when
				problems := Validate(&doc)

				// then
				Expect(problems).To(ContainElement(Problem{
					Severity: Error,
					Message:  "manpage document is missing the 'Synopsis' section'",
				}))
				Expect(doc.Attributes.GetAsStringWithDefault(types.AttrDocType, "")).To(Equal("article")) // changed
			})

			It("missing synopsis section - invalid title", func() {
				// given
				doc := types.Document{
					Attributes: types.DocumentAttributes{
						types.AttrDocType: "manpage",
					},
					ElementReferences: types.ElementReferences{},
					Footnotes:         []types.Footnote{},
					Elements: []interface{}{
						types.Section{
							Attributes: types.ElementAttributes{},
							Level:      0,
							Title: []interface{}{
								types.StringElement{
									Content: "foo",
								},
							},
							Elements: []interface{}{
								types.Section{
									Attributes: types.ElementAttributes{},
									Level:      1,
									Title: []interface{}{
										types.StringElement{
											Content: "Name",
										},
									},
									Elements: []interface{}{
										types.Paragraph{
											Attributes: types.ElementAttributes{},
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "a single paragraph to describe the program",
													},
												},
											},
										},
									},
								},
								types.Section{
									Attributes: types.ElementAttributes{},
									Level:      1,
									Title: []interface{}{
										types.StringElement{
											Content: "bar", // invalid title
										},
									},
									Elements: []interface{}{},
								},
							},
						},
					},
				}

				// when
				problems := Validate(&doc)

				// then
				Expect(problems).To(ContainElement(Problem{
					Severity: Error,
					Message:  "manpage document is missing the 'Synopsis' section'",
				}))
				Expect(doc.Attributes.GetAsStringWithDefault(types.AttrDocType, "")).To(Equal("article")) // changed
			})
		})
	})

})
