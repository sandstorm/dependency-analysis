package parsing

import (
	"reflect"
	"testing"
)

func AssertEquals(t *testing.T, expected interface{}, actual interface{}) {
	if reflect.DeepEqual(actual, expected) {
		return
	}
	t.Errorf("expected %v (type %v), received %v (type %v)",
		expected, reflect.TypeOf(expected),
		actual, reflect.TypeOf(actual))

}
