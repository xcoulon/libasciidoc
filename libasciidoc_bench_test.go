package libasciidoc_test

import (
	"strings"
	"testing"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/bytesparadise/libasciidoc/testsupport"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: unexclude this bench func
func XBenchmarkRenderRealDocument(b *testing.B) {
	filename := "./test/bench/mocking.adoc"
	for i := 0; i < b.N; i++ {
		_, err := testsupport.RenderHTML5Document(filename)
		require.NoError(b, err)
	}
}

func BenchmarkParseBasicDocument(b *testing.B) {
	content := `== Lorem Ipsum
	
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, 
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, 
sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. 
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. 
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, 
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, 
sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. 
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.`

	for i := 0; i < b.N; i++ {
		_, err := testsupport.ParseDocument(content)
		require.NoError(b, err)
	}
}

func BenchmarkParseLongDocument(b *testing.B) {
	content := strings.Builder{}
	for i := 0; i < 50; i++ {
		content.WriteString(`== Lorem Ipsum
	
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, 
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, 
sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. 
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. 
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, 
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, 
sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. 
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.

`)
	}
	for i := 0; i < b.N; i++ {
		_, err := testsupport.ParseDocument(content.String())
		require.NoError(b, err)
	}
}

func TestParseBasicDocument(t *testing.T) {
	source := `== Lorem Ipsum
	
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, 
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, 
sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. 
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. 
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, 
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, 
sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. 
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit *amet*.`

	expected := types.Document{
		Elements: []interface{}{
			&types.Section{
				Level: 1,
				Title: []interface{}{
					types.StringElement{
						Content: "Lorem Ipsum",
					},
				},
			},
			&types.BlankLine{},
			&types.Paragraph{
				Elements: []interface{}{
					types.StringElement{
						Content: `Lorem ipsum dolor sit amet, consetetur sadipscing elitr, 
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, 
sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. 
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. 
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, 
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, 
sed diam voluptua. 
At vero eos et accusam et justo duo dolores et ea rebum. 
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit `,
					},
					&types.QuotedText{
						Kind: types.SingleQuoteBold,
						Elements: []interface{}{
							types.StringElement{
								Content: "amet",
							},
						},
					},
					types.StringElement{
						Content: ".",
					},
				},
			},
		},
	}
	result, err := testsupport.ParseDocument(source)
	require.NoError(t, err)
	assert.Equal(t, expected, result)

}
