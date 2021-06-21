package testsupport

import (
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/parser"
	"github.com/bytesparadise/libasciidoc/pkg/types"
)

// AssembleDocumentFragments assemble the actual source with the options
func AssembleDocumentFragments(actual string) ([]types.DocumentFragmentGroup, error) {
	r := strings.NewReader(actual)
	done := make(chan interface{})
	defer close(done)
	fragmentGroupStream := parser.AssembleFragments(done, parser.ScanDocument(r, done))
	result := []types.DocumentFragmentGroup{}
	for f := range fragmentGroupStream {
		result = append(result, f)
	}
	return result, nil
}
