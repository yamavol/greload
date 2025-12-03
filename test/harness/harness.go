// test assertion library
//
// import only from _test files
package harness

import (
	"runtime"
	"testing"
)

func IsEqual[T comparable](t *testing.T, actual T, expect T, reason string) {
	if actual != expect {
		_, f, l, _ := runtime.Caller(1)
		t.Errorf("%s\n"+
			"  expect: %v\n"+
			"  actual: %v\n"+
			"  in file (%s:%d)\n",
			reason, expect, actual, f, l,
		)
	}
}

func IsNotEqual[T comparable](t *testing.T, actual T, expectNot T, reason string) {
	if actual == expectNot {
		_, f, l, _ := runtime.Caller(1)
		t.Errorf("%s\n"+
			"  Expected not equal but received: %v\n"+
			"  in file (%s:%d)\n",
			reason, actual, f, l,
		)
	}
}

func IsTrue(t *testing.T, condition bool, reason string) {
	if !condition {
		_, f, l, _ := runtime.Caller(1)
		t.Errorf("%s\n"+
			"  Expected true but received false\n"+
			"  in file (%s:%d)\n",
			reason, f, l,
		)
	}
}

func IsFalse(t *testing.T, condition bool, reason string) {
	if condition {
		_, f, l, _ := runtime.Caller(1)
		t.Errorf("%s\n"+
			"  Expected false but received true\n"+
			"  in file (%s:%d)\n",
			reason, f, l,
		)
	}
}

func IsNil[T comparable](t *testing.T, value T, reason string) {
	var zero T
	if value != zero {
		_, f, l, _ := runtime.Caller(1)
		t.Errorf("%s\n"+
			"  Expected nil but received non-nil\n"+
			"  in file (%s:%d)\n",
			reason, f, l,
		)
	}
}

func IsNotNil[T comparable](t *testing.T, value T, reason string) {
	var zero T
	if value == zero {
		_, f, l, _ := runtime.Caller(1)
		t.Errorf("%s\n"+
			"  Expected non-nil but received nil\n"+
			"  in file (%s:%d)\n",
			reason, f, l,
		)
	}
}
