package parser_test

import (
	"fmt"
	"io"

	"github.com/bytesparadise/libasciidoc/pkg/parser"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("document scanning", func() {

	It("should collect 1 fragment with single line", func() {
		source := `a line`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					types.RawLine("a line"),
				},
			},
		}
		Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should collect 1 fragment with multiple lines", func() {
		source := `a line
another line`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					types.RawLine("a line"),
					types.RawLine("another line"),
				},
			},
		}
		Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should collect 2 fragments with non-empty first line", func() {
		source := `a line
		
another line`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					types.RawLine("a line"),
					&types.BlankLine{},
				},
			},
			{
				LineOffset: 3,
				Elements: []interface{}{
					types.RawLine("another line"),
				},
			},
		}
		Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should collect 2 fragments with empty first line", func() {
		source := `
a line
		
another line`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					&types.BlankLine{},
				},
			},
			{
				LineOffset: 2,
				Elements: []interface{}{
					types.RawLine("a line"),
					&types.BlankLine{},
				},
			},
			{
				LineOffset: 4,
				Elements: []interface{}{
					types.RawLine("another line"),
				},
			},
		}
		Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should collect 1 listing block with single line", func() {
		source := `----
a line
----`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					&types.BlockDelimiter{
						Kind: types.Listing,
					},
					types.RawLine("a line"),
					&types.BlockDelimiter{
						Kind: types.Listing,
					},
				},
			},
		}
		Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should collect 1 listing block with multiple rawlines only", func() {
		source := `----
a line

****
not a sidebar block
****
----
`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					&types.BlockDelimiter{
						Kind: types.Listing,
					},
					types.RawLine("a line"),
					types.RawLine(""),
					types.RawLine("****"),
					types.RawLine("not a sidebar block"),
					types.RawLine("****"),
					&types.BlockDelimiter{
						Kind: types.Listing,
					},
				},
			},
		}
		Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should collect 1 listing block with multiple lines and content afterwards", func() {
		source := `----
a line

another line
----


a paragraph
on
3 lines.

`
		expected := []types.DocumentFragment{
			{
				LineOffset: 1,
				Elements: []interface{}{
					&types.BlockDelimiter{
						Kind: types.Listing,
					},
					types.RawLine("a line"),
					types.RawLine(""),
					types.RawLine("another line"),
					&types.BlockDelimiter{
						Kind: types.Listing,
					},
					&types.BlankLine{},
				},
			},
			{
				LineOffset: 7,
				Elements: []interface{}{
					&types.BlankLine{},
				},
			},
			{
				LineOffset: 8,
				Elements: []interface{}{
					types.RawLine("a paragraph"),
					types.RawLine("on"),
					types.RawLine("3 lines."),
					&types.BlankLine{},
				},
			},
		}
		Expect(ParseDocumentFragments(source)).To(MatchDocumentFragments(expected))
	})

	It("should not get an error when reading failed", func() {
		// given
		done := make(chan interface{})
		r := MockReader{}
		// when
		fragmentStream := parser.ScanDocument(r, done)
		// then
		// simplified example: expect a single fragment with an error
		fragment := <-fragmentStream
		Expect(fragment.Error).To(MatchError("mock error"))
		Expect(fragment.Elements).To(BeNil())
	})
})

type MockReader struct{}

var _ io.Reader = MockReader{}

func (r MockReader) Read(_ []byte) (n int, err error) {
	return 0, fmt.Errorf("mock error")
}
