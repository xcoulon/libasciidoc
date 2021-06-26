package parser

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	log "github.com/sirupsen/logrus"
)

func SplitElements(done <-chan interface{}, fragmentStream <-chan types.DocumentFragment) <-chan types.DocumentFragment {
	resultStream := make(chan types.DocumentFragment)
	go func() {
		defer close(resultStream)
		for fragment := range fragmentStream {
			for _, f := range splitElements(fragment) {
				select {
				case <-done:
					log.WithField("pipeline_task", "split_elements").Debug("received 'done' signal")
					return
				case resultStream <- f:
				}
			}
		}
		log.WithField("pipeline_task", "split_elements").Debug("done processing upstream content")
	}()
	return resultStream
}

func splitElements(f types.DocumentFragment) []types.DocumentFragment {
	result := make([]types.DocumentFragment, len(f.Elements))
	for i, element := range f.Elements {
		result[i] = types.NewDocumentFragment(f.LineOffset, []interface{}{element})
	}
	return result
}
