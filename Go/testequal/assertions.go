//go:build !solution

package testequal

import "fmt"

func compare(expected, actual interface{}) bool {
	switch expected := expected.(type) {
	case int:
		if actual, ok := actual.(int); ok && expected == actual {
			return true
		}
	case uint:
		if actual, ok := actual.(uint); ok && expected == actual {
			return true
		}
	case int8:
		if actual, ok := actual.(int8); ok && expected == actual {
			return true
		}
	case uint8:
		if actual, ok := actual.(uint8); ok && expected == actual {
			return true
		}
	case int16:
		if actual, ok := actual.(int16); ok && expected == actual {
			return true
		}
	case uint16:
		if actual, ok := actual.(uint16); ok && expected == actual {
			return true
		}
	case int32:
		if actual, ok := actual.(int32); ok && expected == actual {
			return true
		}
	case uint32:
		if actual, ok := actual.(uint32); ok && expected == actual {
			return true
		}
	case int64:
		if actual, ok := actual.(int64); ok && expected == actual {
			return true
		}
	case uint64:
		if actual, ok := actual.(uint64); ok && expected == actual {
			return true
		}
	case float64:
		if actual, ok := actual.(float64); ok && expected == actual {
			return true
		}
	case string:
		if actual, ok := actual.(string); ok && expected == actual {
			return true
		}
	case nil:
		if actual == nil {
			return true
		}
	case []int:
		if actual, ok := actual.([]int); ok {
			if actual == nil && expected == nil {
				return true
			}
			if expected == nil || actual == nil {
				return false
			}
			if len(expected) != len(actual) {
				return false
			}
			for i := range expected {
				if expected[i] != actual[i] {
					return false
				}
			}
			return true
		}
	case map[string]string:
		if actual, ok := actual.(map[string]string); ok {
			if actual == nil && expected == nil {
				return true
			}
			if expected == nil || actual == nil {
				return false
			}
			if len(expected) != len(actual) {
				return false
			}
			for k, v := range expected {
				if av, ok := actual[k]; !ok || av != v {
					return false
				}
			}
			return true
		}
	case []byte:
		if actual, ok := actual.([]byte); ok {
			if actual == nil && expected == nil {
				return true
			}
			if expected == nil || actual == nil {
				return false
			}
			if len(expected) != len(actual) {
				return false
			}
			for i := range expected {
				if expected[i] != actual[i] {
					return false
				}
			}
			return true
		}
	case struct{}:
		return false
	}
	return false
}

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	if compare(expected, actual) {
		return true
	}
	t.Helper()
	if msgAndArgs != nil {
		t.Errorf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	} else {
		t.Errorf("")
	}

	return false
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	if !compare(expected, actual) {
		return true
	}
	t.Helper()
	if msgAndArgs != nil {
		t.Errorf(fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...))
	} else {
		t.Errorf("")
	}
	return false
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	if !compare(expected, actual) {
		t.Helper()
		if msgAndArgs != nil {
			t.Errorf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		} else {
			t.Errorf("")
		}
		t.FailNow()
	}
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	if compare(expected, actual) {
		t.Helper()
		if msgAndArgs != nil {
			t.Errorf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		} else {
			t.Errorf("")
		}
		t.FailNow()
	}
}
