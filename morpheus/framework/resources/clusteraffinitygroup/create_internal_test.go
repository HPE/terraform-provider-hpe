// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup

import (
	"errors"
	"net/http"
	"testing"
)

// classifyRefusal tells Morpheus 9.1's refusal of a memberless group -- the
// rule-shape error, or the bare 403 that 9.1.0 renders it as -- from any other
// failure. The version gate on the 403 is applied by the caller.
func TestUnitClassifyRefusal(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		err  error
		code int
		want refusalKind
	}{
		{
			name: "the rule-shape error, as 9.1 renders it when servers is sent empty",
			err:  errors.New(`HTTP 400: {"success":false,"msg":"An affinity rule requires either a Servers list (VM-to-VM rule) or both a VM Group and a Host Group (VM-to-Host rule)","errors":{"servers":"affinityGroup.ruleShape.invalid"}}`),
			code: http.StatusBadRequest,
			want: refusalRuleShape,
		},
		{
			name: "the bare 403 9.1.0 answers when servers is absent",
			err:  errors.New(`HTTP 403: {"success":false,"msg":"You do not have permissions to access this api endpoint"}`),
			code: http.StatusForbidden,
			want: refusalForbidden,
		},
		{
			name: "a duplicate name is neither",
			err:  errors.New(`HTTP 400: {"success":false,"msg":"name must be unique","errors":{"name":"affinityGroup.name.unique"}}`),
			code: http.StatusBadRequest,
			want: refusalOther,
		},
		{
			name: "a server error is neither",
			err:  errors.New("HTTP 500: Internal Server Error"),
			code: http.StatusInternalServerError,
			want: refusalOther,
		},
		{
			name: "no response at all is neither",
			err:  errors.New("dial tcp: connection refused"),
			code: 0,
			want: refusalOther,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			var resp *http.Response
			if c.code != 0 {
				resp = &http.Response{StatusCode: c.code}
			}

			if got := classifyRefusal(c.err, resp); got != c.want {
				t.Fatalf("classifyRefusal() = %v, want %v", got, c.want)
			}
		})
	}
}
