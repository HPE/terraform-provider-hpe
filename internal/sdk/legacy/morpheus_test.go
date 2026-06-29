// Test Package for the morpheus-go-sdk package.
// This file defines global assertion methods for use in your tests.

package morpheus_test

import (
	_ "fmt"
	_ "io"
	_ "os"
	"reflect"
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Morpehus SDK Testing Unexported methods
//___________________________________

func assertNil(t *testing.T, v interface{}) {
	if !isNil(v) {
		t.Errorf("[%v] was expected to be nil", v)
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	if isNil(v) {
		t.Errorf("[%v] was expected to be non-nil", v)
	}
}

func assertType(t *testing.T, typ, v interface{}) {
	if reflect.DeepEqual(reflect.TypeOf(typ), reflect.TypeOf(v)) {
		t.Errorf("Expected type %t, got %t", typ, v)
	}
}

func assertError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error occurred [%v]", err)
	}
}

func assertErrorNotNil(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Error expected to be non-nil")
	}
}

func assertEqual(t *testing.T, e, g interface{}) (r bool) {
	if !equal(e, g) {
		t.Errorf("Expected [%v], got [%v]", e, g)
	}

	return
}

func assertNotEqual(t *testing.T, e, g interface{}) (r bool) {
	if equal(e, g) {
		t.Errorf("Expected [%v], got [%v]", e, g)
	} else {
		r = true
	}

	return
}

func equal(expected, got interface{}) bool {
	return reflect.DeepEqual(expected, got)
}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	kind := rv.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && rv.IsNil() {
		return true
	}

	return false
}

func assertResponse(t *testing.T, resp *morpheus.Response, err error) {
	assertError(t, err)
	if resp == nil {
		t.Errorf("API Response was nil")
	} else {
		if resp.Success != true {
			t.Errorf("API Response Error [%v]", err)
			logFailure(t, resp, err)
		}
	}

}

func assertBadResponse(t *testing.T, resp *morpheus.Response, err error) {
	assertErrorNotNil(t, err)
	if resp == nil {
		t.Errorf("API Response was nil")
	} else {
		if resp.Success != false {
			logResponse(t, resp)
			t.Errorf("API Request should have failed and did not [%v]", err)
		}
	}
}

func assertResponseStatusCode(t *testing.T, resp *morpheus.Response, statusCode int) {
	if resp.StatusCode != statusCode {
		logResponse(t, resp)
		t.Errorf("Expected API response status to be %d and got [%d]", statusCode, resp.StatusCode)
	}
}

//todo: make this pretty
// just use logResponse always probably

func logFailure(t *testing.T, resp *morpheus.Response, err error) {
	// logResponse(t, resp)
	t.Logf("API FAILURE: %v", resp)
	if err != nil {
		t.Logf("ERROR: %v", err)
	}
}

func logResponse(t *testing.T, resp *morpheus.Response) {
	t.Logf("API RESPONSE: %v", resp)
	// t.Logf("Response Status: %v", resp.RestyResponse.Status())
	// t.Logf("Response Time: %v", resp.RestyResponse.Time())
	// t.Logf("Response Headers: %v", resp.RestyResponse.Header())
	// t.Logf("Response Cookies: %v", resp.RestyResponse.Cookies())
	// t.Logf("Response Body: %v", resp)
}
