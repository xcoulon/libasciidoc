package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("block images", func() {

	Context("in final documents", func() {

		It("with empty alt", func() {
			source := "image::images/foo.png[]"
			expected := types.Document{
				Elements: []interface{}{
					types.ImageBlock{
						Location: &types.Location{
							Path: []interface{}{
								types.StringElement{Content: "images/foo.png"},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("with attribute alt", func() {
			source := `:alt: the foo.png image
			
image::images/foo.png[{alt}]`
			expected := types.Document{
				Attributes: types.Attributes{
					"alt": "the foo.png image",
				},
				Elements: []interface{}{
					types.ImageBlock{
						Attributes: types.Attributes{
							types.AttrImageAlt: "the foo.png image", // substituted
						},
						Location: &types.Location{
							Path: []interface{}{
								types.StringElement{Content: "images/foo.png"},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("with implicit imagesdir document attribute", func() {
			source := `
:imagesdir: ./path/to/images

image::foo.png[]`
			expected := types.Document{
				Attributes: types.Attributes{
					"imagesdir": "./path/to/images",
				},
				Elements: []interface{}{
					types.ImageBlock{
						Location: &types.Location{
							Path: []interface{}{
								types.StringElement{Content: "./path/to/images/"},
								types.StringElement{Content: "foo.png"},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("with document attribute in URL", func() {
			source := `
:dir: ./path/to/images

image::{dir}/foo.png[]`
			expected := types.Document{
				Attributes: types.Attributes{
					"dir": "./path/to/images",
				},
				Elements: []interface{}{
					types.ImageBlock{
						Location: &types.Location{
							Path: []interface{}{
								types.StringElement{Content: "./path/to/images/foo.png"},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("with implicit imagesdir", func() {
			source := `
:imagesdir: ./path/to/images

image::foo.png[]`
			expected := types.Document{
				Attributes: types.Attributes{
					"imagesdir": "./path/to/images",
				},
				Elements: []interface{}{
					types.ImageBlock{
						Location: &types.Location{
							Path: []interface{}{
								types.StringElement{Content: "./path/to/images/"},
								types.StringElement{Content: "foo.png"},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("with explicit duplicate imagesdir document attribute", func() {
			source := `
:imagesdir: ./path/to/images

image::{imagesdir}/foo.png[]`
			expected := types.Document{
				Attributes: types.Attributes{
					"imagesdir": "./path/to/images",
				},
				Elements: []interface{}{
					types.ImageBlock{
						Location: &types.Location{
							Path: []interface{}{
								types.StringElement{Content: "./path/to/images/"},
								types.StringElement{Content: "./path/to/images/foo.png"},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("2 block images", func() {
			source := `image::images/foo.png[]
image::images/bar.png[]`
			expected := types.Document{
				Elements: []interface{}{
					types.ImageBlock{
						Location: &types.Location{
							Path: []interface{}{
								types.StringElement{Content: "images/foo.png"},
							},
						},
					},
					types.ImageBlock{
						Location: &types.Location{
							Path: []interface{}{
								types.StringElement{Content: "images/bar.png"},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
	})

	Context("errors", func() {

		It("appending inline content", func() {
			source := "a paragraph\nimage::images/foo.png[]"
			expected := types.Document{
				Elements: []interface{}{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{
									Content: "a paragraph",
								},
							},
							{
								types.StringElement{
									Content: "image::images/foo.png[]",
								},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("paragraph with block image with alt and dimensions", func() {
			source := "a foo image::foo.png[foo image, 600, 400] bar"
			expected := types.Document{
				Elements: []interface{}{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "a foo image::foo.png[foo image, 600, 400] bar"},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
	})
})

var _ = Describe("inline images", func() {

	Context("in final documents", func() {

		It("with empty alt only", func() {
			source := "image:images/foo.png[]"
			expected := types.Document{
				Elements: []interface{}{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.InlineImage{
									Location: &types.Location{
										Path: []interface{}{
											types.StringElement{Content: "images/foo.png"},
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

		It("with document attribute in URL", func() {
			source := `
:dir: ./path/to/images

an image:{dir}/foo.png[].`
			expected := types.Document{
				Attributes: types.Attributes{
					"dir": "./path/to/images",
				},
				Elements: []interface{}{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "an "},
								types.InlineImage{
									Location: &types.Location{
										Path: []interface{}{
											types.StringElement{Content: "./path/to/images/foo.png"},
										},
									},
								},
								types.StringElement{Content: "."},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("with implicit imagesdir document attribute", func() {
			source := `
:imagesdir: ./path/to/images

an image:foo.png[].`
			expected := types.Document{
				Attributes: types.Attributes{
					"imagesdir": "./path/to/images",
				},
				Elements: []interface{}{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "an "},
								types.InlineImage{
									Location: &types.Location{
										Path: []interface{}{
											types.StringElement{Content: "./path/to/images/"},
											types.StringElement{Content: "foo.png"},
										},
									},
								},
								types.StringElement{Content: "."},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("with explicit duplicate imagesdir document attribute", func() {
			source := `
:imagesdir: ./path/to/images

an image:{imagesdir}/foo.png[].`
			expected := types.Document{
				Attributes: types.Attributes{
					"imagesdir": "./path/to/images",
				},
				Elements: []interface{}{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "an "},
								types.InlineImage{
									Location: &types.Location{
										Path: []interface{}{
											types.StringElement{Content: "./path/to/images/"},
											types.StringElement{Content: "./path/to/images/foo.png"},
										},
									},
								},
								types.StringElement{Content: "."},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})

		It("with document attribute in URL", func() {
			source := `:path: ./path/to/images

image::{path}/foo.png[]`
			expected := types.Document{
				Attributes: types.Attributes{
					"path": "./path/to/images",
				},
				Elements: []interface{}{
					types.ImageBlock{
						Location: &types.Location{
							Path: []interface{}{
								types.StringElement{Content: "./path/to/images/foo.png"}, // resolved
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
	})

	Context("errors", func() {

		It("appending inline content", func() {
			source := "a paragraph\nimage::images/foo.png[]"
			expected := types.Document{
				Elements: []interface{}{
					types.Paragraph{
						Lines: [][]interface{}{
							{
								types.StringElement{Content: "a paragraph"},
							},
							{
								types.StringElement{Content: "image::images/foo.png[]"},
							},
						},
					},
				},
			}
			Expect(ParseDocument(source)).To(MatchDocument(expected))
		})
	})
})
