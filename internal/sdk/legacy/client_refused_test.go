// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// A Morpheus appliance that answers every request with a fixed status and
// body, for testing how the client judges the response.
func fakeAppliance(t *testing.T, status int, body string) *Client {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(srv.URL)
	c.SetAccessToken("t", "", 0, "write")

	return c
}

func TestUnitExecuteJudgesTheBodyNotJustTheStatus(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		status      int
		body        string
		wantSuccess bool
		wantErr     string // "" for nil
		wantRefused string
	}{
		{"200 success:true", 200, `{"success":true}`, true, "", ""},
		{"200 no success field", 200, `{"user":{"id":1}}`, true, "", ""},
		{"200 empty body", 200, ``, true, "", ""},
		{"200 success as a number is not the flag", 200, `{"success":3}`, true, "", ""},
		{"200 JSON array", 200, `[1,2]`, true, "", ""},
		{
			"200 success:false with msg", 200, `{"success":false,"msg":"Network is in use"}`,
			false, "200 (OK): Morpheus refused the request: Network is in use", "Network is in use",
		},
		{
			"200 success:false with errors", 200, `{"success":false,"errors":{"b":"bad","a":"worse"}}`,
			false, "Morpheus refused the request: a: worse; b: bad", "a: worse; b: bad",
		},
		{
			"200 success:false, no reason", 200, `{"success":false}`,
			false, "Morpheus refused the request: no reason given", "no reason given",
		},
		{"404", 404, `{"success":false,"msg":"not found"}`, false, "404 (Not Found)", ""},
		{"500", 500, `boom`, false, "500 (Internal Server Error)", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := fakeAppliance(t, tc.status, tc.body)

			resp, err := c.Delete(&Request{Path: "/api/things/1"})
			if resp == nil {
				t.Fatal("expected a response")
			}

			if resp.Success != tc.wantSuccess {
				t.Fatalf("Success = %v, want %v (err=%v)", resp.Success, tc.wantSuccess, err)
			}

			switch {
			case tc.wantErr == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case tc.wantErr != "" && err == nil:
				t.Fatalf("expected an error containing %q", tc.wantErr)
			case tc.wantErr != "" && !strings.Contains(err.Error(), tc.wantErr):
				t.Fatalf("error %q does not contain %q", err.Error(), tc.wantErr)
			}

			if resp.RefusedMessage != tc.wantRefused {
				t.Fatalf("RefusedMessage = %q, want %q", resp.RefusedMessage, tc.wantRefused)
			}
		})
	}
}

// The refusal also applies to the other verbs: it is the response view that
// is shared, not the delete controller.
func TestUnitExecuteRefusalAppliesToEveryVerb(t *testing.T) {
	t.Parallel()

	c := fakeAppliance(t, 200, `{"success":false,"msg":"nope"}`)

	for name, call := range map[string]func(*Request) (*Response, error){
		"GET": c.Get, "POST": c.Post, "PUT": c.Put, "DELETE": c.Delete,
	} {
		resp, err := call(&Request{Path: "/api/things", Body: map[string]interface{}{}})
		if err == nil || resp == nil || resp.Success {
			t.Errorf("%s: a 200 with success:false must be a failure (err=%v)", name, err)
		}
	}
}

// A typed Result still decodes on a refusal, so callers that inspect
// StandardResult see the same message.
func TestUnitExecuteRefusalKeepsTheDecodedResult(t *testing.T) {
	t.Parallel()

	c := fakeAppliance(t, 200, `{"success":false,"msg":"in use","errors":{"x":"y"}}`)

	resp, err := c.Delete(&Request{Path: "/api/things/1", Result: &StandardResult{}})
	if err == nil {
		t.Fatal("expected an error")
	}

	r, ok := resp.Result.(*StandardResult)
	if !ok || r.Success || r.Message != "in use" || r.Errors["x"] != "y" {
		t.Fatalf("decoded result not as expected: %+v", resp.Result)
	}
}
