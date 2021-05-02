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

var _ = Describe("document fragment groups parsing", func() {

	It("should collect 1 fragment", func() {
		source := `a line`
		expected := []types.DocumentFragmentGroup{
			{
				LineOffset: 1,
				Content: []interface{}{
					types.RawLine("a line"),
				},
			},
		}
		Expect(ParseDocumentFragmentGroups(source)).To(MatchDocumentFragmentGroups(expected))
	})

	It("should collect 2 fragments with non-empty first line", func() {
		source := `a line
		
another line`
		expected := []types.DocumentFragmentGroup{
			{
				LineOffset: 1,
				Content: []interface{}{
					types.RawLine("a line"),
					types.BlankLine{},
				},
			},
			{
				LineOffset: 3,
				Content: []interface{}{
					types.RawLine("another line"),
				},
			},
		}
		Expect(ParseDocumentFragmentGroups(source)).To(MatchDocumentFragmentGroups(expected))
	})

	It("should collect 2 fragments with empty first line", func() {
		source := `
a line
		
another line`
		expected := []types.DocumentFragmentGroup{
			{
				LineOffset: 1,
				Content: []interface{}{
					types.BlankLine{},
				},
			},
			{
				LineOffset: 2,
				Content: []interface{}{
					types.RawLine("a line"),
					types.BlankLine{},
				},
			},
			{
				LineOffset: 4,
				Content: []interface{}{
					types.RawLine("another line"),
				},
			},
		}
		Expect(ParseDocumentFragmentGroups(source)).To(MatchDocumentFragmentGroups(expected))
	})

	It("should not get an error when reading failed", func() {
		// given
		done := make(chan interface{})
		r := MockReader{}
		// when
		fragmentStream := parser.ParseDocumentFragmentGroups(r, done)
		// then
		// simplified example: expect a single fragment with an error
		fragment := <-fragmentStream
		Expect(fragment.Error).To(MatchError("mock error"))
		Expect(fragment.Content).To(BeNil())
	})
})

type MockReader struct{}

var _ io.Reader = MockReader{}

func (r MockReader) Read(_ []byte) (n int, err error) {
	return 0, fmt.Errorf("mock error")
}
