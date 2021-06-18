package testsupport

import (
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/parser"
	"github.com/bytesparadise/libasciidoc/pkg/types"
)

// ParseDocumentFragmentGroups parses the actual source with the options
func ParseDocumentFragmentGroups(actual string, options ...interface{}) ([]types.DocumentFragmentGroup, error) {
	r := strings.NewReader(actual)
	c := &rawDocumentParserConfig{
		filename: "test.adoc",
	}
	parserOptions := []parser.Option{}
	for _, o := range options {
		switch set := o.(type) {
		case FilenameOption:
			set(c)
		case parser.Option:
			parserOptions = append(parserOptions, set)
		}
	}
	done := make(chan interface{})
	defer close(done)
	fragmentStream := parser.ScanDocument(r, done, parserOptions...)
	result := []types.DocumentFragmentGroup{}
	for f := range fragmentStream {
		result = append(result, f)
	}
	return result, nil
}

type rawDocumentParserConfig struct {
	filename string
}

func (c *rawDocumentParserConfig) setFilename(f string) {
	c.filename = f
}
