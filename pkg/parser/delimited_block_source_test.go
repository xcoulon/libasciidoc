package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("source blocks", func() {

	Context("in final documents", func() {

		Context("delimited block", func() {

			It("with source attribute only", func() {
				source := `[source]
----
require 'sinatra'

get '/hi' do
  "Hello World!"
end
----`
				expected := types.Document{
					Elements: []interface{}{
						types.ListingBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Source,
							},
							Elements: []interface{}{
								types.StringElement{
									Content: "require 'sinatra'",
								},
								types.StringElement{
									Content: "get '/hi' do",
								},
								types.StringElement{
									Content: "  \"Hello World!\"",
								},
								types.StringElement{
									Content: "end",
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with title, source and languages attributes", func() {
				source := `[source,ruby]
.Source block title
----
require 'sinatra'

get '/hi' do
  "Hello World!"
end
----`
				expected := types.Document{
					Elements: []interface{}{
						types.ListingBlock{
							Attributes: types.Attributes{
								types.AttrStyle:    types.Source,
								types.AttrLanguage: "ruby",
								types.AttrTitle:    "Source block title",
							},
							Elements: []interface{}{
								types.StringElement{
									Content: "require 'sinatra'",
								},
								types.StringElement{
									Content: "get '/hi' do",
								},
								types.StringElement{
									Content: "  \"Hello World!\"",
								},
								types.StringElement{
									Content: "end",
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with id, title, source and languages attributes", func() {
				source := `[#id-for-source-block]
[source,ruby]
.app.rb
----
require 'sinatra'

get '/hi' do
  "Hello World!"
end
----`
				expected := types.Document{
					Elements: []interface{}{
						types.ListingBlock{
							Attributes: types.Attributes{
								types.AttrStyle:    types.Source,
								types.AttrLanguage: "ruby",
								types.AttrID:       "id-for-source-block",
								types.AttrTitle:    "app.rb",
							},
							Elements: []interface{}{
								types.StringElement{
									Content: "require 'sinatra'",
								},
								types.StringElement{
									Content: "get '/hi' do",
								},
								types.StringElement{
									Content: "  \"Hello World!\"",
								},
								types.StringElement{
									Content: "end",
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with callout and admonition block afterwards", func() {
				source := `[source]
----
const cookies = "cookies" <1>
----
<1> a constant

[NOTE]
====
a note
====`

				expected := types.Document{
					Elements: []interface{}{
						types.ListingBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Source,
							},
							Elements: []interface{}{
								types.StringElement{
									Content: `const cookies = "cookies" `,
								},
								types.Callout{
									Ref: 1,
								},
							},
						},
						types.CalloutList{
							Items: []types.CalloutListItem{
								{
									Ref: 1,
									Elements: []interface{}{
										types.Paragraph{
											Lines: [][]interface{}{
												{
													types.StringElement{
														Content: "a constant",
													},
												},
											},
										},
									},
								},
							},
						},
						types.ExampleBlock{
							Attributes: types.Attributes{
								types.AttrStyle: types.Note,
							},
							Elements: []interface{}{
								types.Paragraph{
									Lines: [][]interface{}{
										{
											types.StringElement{
												Content: "a note",
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
