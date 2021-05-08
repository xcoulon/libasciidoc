package testsupport

import (
	"fmt"
	"reflect"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/davecgh/go-spew/spew"
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"
	log "github.com/sirupsen/logrus"
)

// MatchDocumentFragmentGroups a custom matcher to verify that a document matches the given expectation
// Similar to the standard `Equal` matcher, but display a diff when the values don't match
func MatchDocumentFragments(expected []types.DocumentFragment) gomegatypes.GomegaMatcher {
	return &documentFragmentsMatcher{
		expected: expected,
	}
}

type documentFragmentsMatcher struct {
	expected []types.DocumentFragment
	diffs    string
}

func (m *documentFragmentsMatcher) Match(actual interface{}) (success bool, err error) {
	if _, ok := actual.([]types.DocumentFragment); !ok {
		return false, errors.Errorf("MatchDocumentFragments matcher expects an array of types.DocumentFragmentGroup (actual: %T)", actual)
	}
	if !reflect.DeepEqual(m.expected, actual) {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debugf("actual raw document:\n%s", spew.Sdump(actual))
			log.Debugf("expected raw document:\n%s", spew.Sdump(m.expected))
		}
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(spew.Sdump(actual), spew.Sdump(m.expected), true)
		m.diffs = dmp.DiffPrettyText(diffs)
		return false, nil
	}
	return true, nil
}

func (m *documentFragmentsMatcher) FailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("expected document fragments to match:\n%s", m.diffs)
}

func (m *documentFragmentsMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("expected document fragments not to match:\n%s", m.diffs)
}