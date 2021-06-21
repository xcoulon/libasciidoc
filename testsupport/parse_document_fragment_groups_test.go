package testsupport_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("parse document fragment groups", func() {

	expected := []types.DocumentFragment{
		{
			LineOffset: 1,
			Content: []interface{}{
				types.RawLine("hello, world!"),
			},
		},
	}

	It("should match", func() {
		// given
		actual := "hello, world!"
		// when
		result, err := testsupport.ParseDocumentFragmentGroups(actual)
		// then
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal(expected))
	})

	It("should not match", func() {
		// given
		actual := "foo"
		// when
		result, err := testsupport.ParseDocumentFragmentGroups(actual)
		// then
		Expect(err).ToNot(HaveOccurred())
		Expect(result).NotTo(Equal(expected))
	})
})
