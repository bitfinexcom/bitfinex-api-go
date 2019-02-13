package tests

import (
	"bytes"
	"reflect"
	"testing"
)

func isZeroOfUnderlyingType(x interface{}) bool {
	if x == nil {
		return true
	}
	if _, ok := x.([]string); ok {
		return true
	}
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

func objEq(expected, actual interface{}) bool {

	if expected == nil || actual == nil {
		return expected == actual
	}
	if exp, ok := expected.([]byte); ok {
		act, ok := actual.([]byte)
		if !ok {
			return false
		} else if exp == nil || act == nil {
			return exp == nil && act == nil
		}
		return bytes.Equal(exp, act)
	}
	return reflect.DeepEqual(expected, actual)

}

func assertSlice(t *testing.T, expected, actual interface{}) {
	if objEq(expected, actual) {
		t.Logf("%s OK", reflect.TypeOf(actual))
		return
	}

	actualType := reflect.TypeOf(actual)
	if actualType == nil {
		t.Fatal()
	}
	expectedValue := reflect.ValueOf(expected)
	if expectedValue.IsValid() && expectedValue.Type().ConvertibleTo(actualType) {
		// Attempt comparison after type conversion
		if reflect.DeepEqual(expectedValue.Convert(actualType).Interface(), actual) {
			t.Logf("%s OK", reflect.TypeOf(actual))
			return
		}
	}

	t.Fatalf("FAIL %s: expected %#v, got %#v", reflect.TypeOf(expected), expected, actual)
}

func isPrimitive(exp interface{}) bool {
	t := reflect.TypeOf(exp)
	switch t.Kind() {
	case reflect.Interface:
		return false
	case reflect.Struct:
		return false
	case reflect.Array:
		return false
	case reflect.Func:
		return false
	case reflect.Map:
		return false
	case reflect.Ptr:
		return false
	case reflect.Slice:
		return false
	case reflect.UnsafePointer:
		return false
	default:
		return true
	}
}

// does not work on slices
func assert(t *testing.T, expected interface{}, actual interface{}) {

	prexp := reflect.ValueOf(expected)
	pract := reflect.ValueOf(actual)

	if isPrimitive(actual) {
		if expected != actual {
			t.Fatalf("expected %#v, got %#v", expected, actual)
		}
		t.Logf("OK %s", reflect.TypeOf(expected).Name())
		return
	}

	if pract.IsNil() {
		t.Errorf("nil actual value: %#v", actual)
		t.Fail()
		return
	}

	exp := prexp.Elem()
	act := pract.Elem()

	if !exp.IsValid() {
		t.Errorf("reflected expectation not valid (%#v)", expected)
		t.Fail()
	}

	if exp.Type() != act.Type() {
		t.Errorf("expected type %s, got %s", exp.Type(), act.Type())
		t.Fail()
	}

	for i := 0; i < exp.NumField(); i++ {
		expValueField := exp.Field(i)
		expTypeField := exp.Type().Field(i)

		actValueField := act.Field(i)
		actTypeField := act.Type().Field(i)

		if expTypeField.Name != actTypeField.Name {
			t.Errorf("expected type %s, got %s", expTypeField.Name, actTypeField.Name)
			t.Errorf("%#v", actual)
			t.Fail()
		}
		if isZeroOfUnderlyingType(expValueField.Interface()) {
			continue
		}
		if !isZeroOfUnderlyingType(expValueField.Interface()) && isZeroOfUnderlyingType(actValueField.Interface()) {
			t.Errorf("expected %s, but was empty", expTypeField.Name)
			t.Errorf("%#v", actual)
			t.Fail()
			return
		}
		assert(t, expValueField.Interface(), actValueField.Interface())
		/*
			if expValueField.Interface() != actValueField.Interface() {
				t.Errorf("expected %s %#v, got %#v", expTypeField.Name, expValueField.Interface(), actValueField.Interface())
				t.Fail()
			}
		*/
	}
	if t.Failed() {
		t.Logf("FAIL %s", exp.Type().Name())
		return
	}
	t.Logf("OK %s", exp.Type().Name())
}
