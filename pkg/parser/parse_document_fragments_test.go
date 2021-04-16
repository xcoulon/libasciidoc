package parser_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/parser"
	"github.com/bytesparadise/libasciidoc/pkg/types"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("document fragments parsing", func() {

	It("should collect fragments", func() {
		// given
		done := make(chan interface{})
		r := strings.NewReader(`a line`)
		// when
		fragmentStream := parser.ParseDocumentFragments(r, done)
		// then
		// simplified example: expect a single valid fragment
		fragment := <-fragmentStream
		Expect(fragment).To(Equal(types.DocumentFragment{
			LineOffset: 1,
			Content: []interface{}{
				types.RawLine("a line"),
			},
		}))
	})

	It("should get an error when reading", func() {
		// given
		done := make(chan interface{})
		r := MockReader{}
		// when
		fragmentStream := parser.ParseDocumentFragments(r, done)
		// then
		// simplified example: expect a single fragment with an error
		fragment := <-fragmentStream
		Expect(fragment.Error).To(MatchError("mock error"))
		Expect(fragment.Content).To(BeNil())
	})

	It("should get an error when parsing", func() {
		// given
		done := make(chan interface{})
		r := strings.NewReader(`= Title`)
		// when
		fragmentStream := parser.ParseDocumentFragments(r, done, parser.MaxExpressions(1))
		// then
		// simplified example: expect a single fragment with an error
		fragment := <-fragmentStream
		Expect(fragment.Error).To(MatchError("1:1 (0): rule DocumentFragmentElement: max number of expresssions parsed"))
		Expect(fragment.Content).To(BeNil())
	})
})

type MockReader struct{}

var _ io.Reader = MockReader{}

func (r MockReader) Read(_ []byte) (n int, err error) {
	return 0, fmt.Errorf("mock error")
}
