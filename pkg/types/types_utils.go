package types

import (
	"fmt"
	"strings"
)

// Merge merge string elements together
func Merge(elements ...interface{}) []interface{} {
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debugf("merging %s", spew.Sdump(elements))
	// }
	result := make([]interface{}, 0, len(elements))
	buf := &strings.Builder{}
	for _, element := range elements {
		if element == nil {
			continue
		}
		switch element := element.(type) {
		case string:
			buf.WriteString(element)
		case []byte:
			buf.Write(element)
		case *StringElement:
			buf.WriteString(element.Content)
		case []interface{}:
			if len(element) > 0 {
				f := Merge(element...)
				result, buf = appendBuffer(result, buf)
				result = Merge(append(result, f...)...)
			}
		default:
			// log.Debugf("Merging with 'default' case an element of type %[1]T", element)
			result, buf = appendBuffer(result, buf)
			result = append(result, element)
		}
	}
	// if buf was filled because some text was found
	result, _ = appendBuffer(result, buf)
	return result
}

// AllNilEntries returns true if all the entries in the given `elements` are `nil`
func AllNilEntries(elements []interface{}) bool {
	for _, e := range elements {
		switch e := e.(type) {
		case []interface{}: // empty slice if not `nil` since it has a type
			if !AllNilEntries(e) {
				return false
			}
		default:
			if e != nil {
				return false
			}
		}
	}
	return true
}

// appendBuffer appends the content of the given buffer to the given array of elements,
// and returns a new buffer, or returns the given arguments if the buffer was empty
func appendBuffer(elements []interface{}, buf *strings.Builder) ([]interface{}, *strings.Builder) {
	if buf.Len() > 0 {
		s, _ := NewStringElement(buf.String())
		return append(elements, s), &strings.Builder{}
	}
	return elements, buf
}

// ReduceOption an option to apply on the reduced content when it is a `string`
type ReduceOption func(string) string

// Reduce merges and returns a string if the given elements only contain a single StringElement
// (ie, return its `Content`), otherwise return the given elements or empty string if the elements
// is `nil` or an empty `[]interface{}`
func Reduce(elements interface{}, opts ...ReduceOption) interface{} {
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debugf("reducing %s", spew.Sdump(elements))
	// }
	switch e := elements.(type) {
	case []interface{}:
		e = Merge(e...)
		switch len(e) {
		case 0: // if empty, return nil
			return nil
		case 1:
			if e, ok := e[0].(*StringElement); ok {
				c := e.Content
				for _, apply := range opts {
					c = apply(c)
				}
				return c
			}
			return e
		default:
			return e
		}
	case string:
		for _, apply := range opts {
			e = apply(e)
		}
		return e
	default:
		return elements
	}
}

// applyFunc a function to apply on the result of the `apply` function below, before returning
type applyFunc func(s string) string

// Apply applies the given funcs to transform the given input
func Apply(source string, fs ...applyFunc) string {
	result := source
	for _, f := range fs {
		result = f(result)
	}
	return result
}

func stringify(element interface{}) string {
	switch element := element.(type) {
	case []interface{}:
		result := strings.Builder{}
		for _, e := range element {
			result.WriteString(stringify(e))
		}
		return result.String()
	case string:
		return element
	case *StringElement:
		return element.Content
	case AttributeSubstitution:
		return "{" + element.Name + "}"
	default:
		return fmt.Sprintf("%v", element) // "best-effort" here
	}
}
