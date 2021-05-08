package testsupport

import (
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/parser"
	"github.com/bytesparadise/libasciidoc/pkg/types"
)

// AssembleDocumentFragments assemble the actual source with the options
func AssembleDocumentFragments(actual string) ([]types.DocumentFragment, error) {
	r := strings.NewReader(actual)
	done := make(chan interface{})
	defer close(done)
	fragmentStream := parser.AssembleFragments(done, parser.ParseDocumentFragmentGroups(r, done))
	result := []types.DocumentFragment{}
	for f := range fragmentStream {
		result = append(result, f)
	}
	return result, nil
}
