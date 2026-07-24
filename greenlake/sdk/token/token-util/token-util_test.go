// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package tokenutil

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	hpeglErrors "github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/errors"
)

var errLimitExceeded = hpeglErrors.MakeErrInternalError(hpeglErrors.ErrorResponse{
	ErrorCode: "ErrGenerateTokenRetryLimitExceeded",
	Message:   "Retry limit exceeded"})

//nolint:scopelint
func TestDecodeAccessToken(t *testing.T) {
	type args struct {
		rawToken string
	}
	tests := []struct {
		name    string
		args    args
		want    Token
		wantErr bool
	}{
		{
			name: "Decode access token",
			args: args{
				rawToken: "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJGd1BrQlJUbWY2Rm54R045SFp2VVhmOGhGa0xiOVI4" +
					"am9XSWVORzNSNlJZIn0.eyJqdGkiOiJkZjRlYjg1OS0zOTA2LTRmNDktOGRlZi0wMzA5MjJiZTRkMDYiLCJleHAiOjE1NTc5Nj" +
					"I3NjQsIm5iZiI6MCwiaWF0IjoxNTU3OTQ4MzY0LCJpc3MiOiJodHRwczovL2lkcC1icm9rZXIuZGV2LmhwZWRldm9wcy5uZXQv" +
					"YXV0aC9yZWFsbXMvbWFzdGVyIiwiYXVkIjpbIm1hc3Rlci1yZWFsbSIsImFjY291bnQiXSwic3ViIjoiMTBhYzIxZDQtNTg5Yi" +
					"00M2NhLWE0MzQtYWQxOTk5ZjM5NTBhIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoib25lc3BoZXJlLWF1dGgiLCJub25jZSI6IjEy" +
					"MzQ1IiwiYXV0aF90aW1lIjoxNTU3OTQ4MzYzLCJzZXNzaW9uX3N0YXRlIjoiOGI1YTY1ZGYtODkwOS00ZmE3LWJmNWEtYWQ2ZT" +
					"M5M2Y0ZDVlIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyJodHRwczovL2lkcC1icm9rZXIuZGV2LmhwZWRldm9wcy5u" +
					"ZXQiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbImNyZWF0ZS1yZWFsbSIsIm9mZmxpbmVfYWNjZXNzIiwiYWRtaW4iLCJ1bW" +
					"FfYXV0aG9yaXphdGlvbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7Im1hc3Rlci1yZWFsbSI6eyJyb2xlcyI6WyJ2aWV3LWlkZW50" +
					"aXR5LXByb3ZpZGVycyIsInZpZXctcmVhbG0iLCJtYW5hZ2UtaWRlbnRpdHktcHJvdmlkZXJzIiwiaW1wZXJzb25hdGlvbiIsIm" +
					"NyZWF0ZS1jbGllbnQiLCJtYW5hZ2UtdXNlcnMiLCJxdWVyeS1yZWFsbXMiLCJ2aWV3LWF1dGhvcml6YXRpb24iLCJxdWVyeS1j" +
					"bGllbnRzIiwicXVlcnktdXNlcnMiLCJtYW5hZ2UtZXZlbnRzIiwibWFuYWdlLXJlYWxtIiwidmlldy1ldmVudHMiLCJ2aWV3LX" +
					"VzZXJzIiwidmlldy1jbGllbnRzIiwibWFuYWdlLWF1dGhvcml6YXRpb24iLCJtYW5hZ2UtY2xpZW50cyIsInF1ZXJ5LWdyb3Vw" +
					"cyJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm" +
					"9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgdGVuYW50SWQgcHJvZmlsZSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJu" +
					"YW1lIjoiQWJoaXNoZWsgU3JpdmFzdGF2YSIsInRlbmFudElkIjoiOHI0MC1sZDlwIiwicHJlZmVycmVkX3VzZXJuYW1lIjoic3" +
					"JpdmFzdGF2YS5hYmhpc2hla0BocGUuY29tIiwiZ2l2ZW5fbmFtZSI6IkFiaGlzaGVrIiwiZmFtaWx5X25hbWUiOiJTcml2YXN0" +
					"YXZhIiwiZW1haWwiOiJzcml2YXN0YXZhLmFiaGlzaGVrQGhwZS5jb20ifQ.FP40XhVuO_-1ZsVH4FKzK6HKUig7F1Ahwl2c9RQ" +
					"CDGp2WxA3xkCOoA_5lFwsVJatCT1f9vRw0LuTDRCIo-8bLv374X8F8V1rThK5ReUBom5-0ul4rruWSL13VTfhgkRhtrjNIwp00" +
					"IgitVO_mO_hpjusOqQk3uWglpkI1zbFrP5kXaOmy6qj_RLnxS3NCOSdvieEs_r_YN5TuzE75T6OP3_AxFUKpbnOLIH_5TTQtQk" +
					"oZJfGaIge195FCT5i1o6MqCP3xB0_zfyIbP86lhgyTfyow1SaiDf2uEvtKSrZtmxgJjhAGLWPClSklGX4sky6mnfqNd-ReF5LZ" +
					"rqEClUAeQ",
			},
			want: Token{
				Issuer:          "https://idp-broker.dev.hpedevops.net/auth/realms/master",
				Subject:         "users/10ac21d4-589b-43ca-a434-ad1999f3950a",
				TenantID:        "8r40-ld9p",
				Expiry:          1557962764,
				Type:            "Bearer",
				IssuedAt:        1557948364,
				Nonce:           "12345",
				AuthorizedParty: "onesphere-auth",
			},
			wantErr: false,
		},
		{
			name: "Decode user access token from okta",
			args: args{
				rawToken: "eyJraWQiOiI3Ykc1T1JxY0NLbkZZYV9zYXVIM0JFV2VuaC1kTmt5WWlQd0NKSm9velU4IiwiYWxnIjoiUlMyNTYifQ." +
					"eyJ2ZXIiOjEsImp0aSI6IkFULm9Od3k4RWpheW9WRU10dm9YemJNMi1sdXhwT1FVM1BqcTZYTl9lV3FGWkUiLCJpc3MiOiJod" +
					"HRwczovL2hwZS1ncmVlbmxha2Uub2t0YXByZXZpZXcuY29tL29hdXRoMi9kZWZhdWx0IiwiYXVkIjoiYXBpOi8vZGVmYXVsdC" +
					"IsImlhdCI6MTU2NjQwMzk1MCwiZXhwIjoxNTY2NDA3NTUwLCJjaWQiOiIwb2FtczQ2OGZ3eHhZcDZ6UjBoNyIsInVpZCI6IjA" +
					"wdW4xem9xcHFSRDEwWDVLMGg3Iiwic2NwIjpbIm9wZW5pZCIsInByb2ZpbGUiLCJlbWFpbCJdLCJzdWIiOiJyeWFuLmJyYW5k" +
					"dEBocGUuY29tIiwidGVuYW50SWQiOiJocGUtZ3JlZW5sYWtlLWludGcifQ.dwo2de1BSam5kME8bqQnQWBys6fGfRfGR05z1l" +
					"fJSq1BTIuXGPcxpkmKOyjOFzPR1VoteQp1ZwC02Qj9WyYabgGI6LZ5odPv01GT9hAVUc0RijFPrp9w666LDPv_1LAyzgeccfAl" +
					"5MZNeVDk7OdAZX7UiSj_mLnDAJv9Dzln-CMJA3JgA5GjZVZo0mdvfzD6mCqNL594za_ouZt5_Pp4TPoAAzGkYSjyv8LJWT0Zrs" +
					"ruzPKlS1UTJR6wUSxlqCSQnjZtXQtE5UugQmxdV1F_SC6HxikabnCu-6zLTPgJANUQ1XXKSE2brEgkYIx4suTv6IzejGz1Ccnw" +
					"O__RLEJ2gQ",
			},
			want: Token{
				Issuer:   "https://hpe-greenlake.oktapreview.com/oauth2/default",
				Subject:  "users/00un1zoqpqRD10X5K0h7",
				ClientID: "0oams468fwxxYp6zR0h7",
				UserID:   "00un1zoqpqRD10X5K0h7",
				TenantID: "hpe-greenlake-intg",
				Expiry:   1566407550,
				IssuedAt: 1566403950,
			},
			wantErr: false,
		},
		{
			name: "Decode client access token from okta",
			args: args{
				rawToken: "eyJraWQiOiI3Ykc1T1JxY0NLbkZZYV9zYXVIM0JFV2VuaC1kTmt5WWlQd0NKSm9velU4IiwiYWxnIjoiUlMyNTYifQ." +
					"eyJ2ZXIiOjEsImp0aSI6IkFULlFtZjQ0djhjODZiSjE5UDliMlhFcjczS0RrOGY0Nzdfd3hDZE94SF82bjgiLCJpc3MiOiJo" +
					"dHRwczovL2hwZS1ncmVlbmxha2Uub2t0YXByZXZpZXcuY29tL29hdXRoMi9kZWZhdWx0IiwiYXVkIjoiYXBpOi8vZGVmYXVs" +
					"dCIsImlhdCI6MTU2NjQxMjc0OCwiZXhwIjoxNTY2NDE2MzQ4LCJjaWQiOiIwb2FuMnh2a2k1UXF5aHNjdTBoNyIsInNjcCI6Wy" +
					"J0ZW5hbnQtc2NvcGUiXSwic3ViIjoiMG9hbjJ4dmtpNVFxeWhzY3UwaDciLCJ0ZW5hbnRJZCI6ImhwZS1ncmVlbmxha2UtaW50ZyJ9." +
					"lUogzirlS8RKe0RZKzubUG-gG-HfZIfiTfMf_GXgvsvozlRaqDIPKlzwYasuXmaJVV83VhFrXsKqVimIcA6lIjjKdETzLXMY3O35" +
					"HajlpLUlzwwupvTyzUmUPzalNzl85Fpdqqc42E3HA0PIesJm7X3GFMwy0OFzhZlf0PyAQ2gEikg-LQCRXiWEeE_INj1OQtACYbMJg" +
					"I4vkWk_PqYFev4_tkEwctABy-1jg_ecHbSwYsAsgKfVALl7-tPwO-cs-zyRtzTh_ltPfW89kD4OFj6XQppZKig3yhdhSvIcNUcuQS" +
					"OksLTh0IxldW7J_ct2tYx3lJ1BnGpwenlfyBKCIw",
			},
			want: Token{
				Issuer:   "https://hpe-greenlake.oktapreview.com/oauth2/default",
				Subject:  "clients/0oan2xvki5Qqyhscu0h7",
				ClientID: "0oan2xvki5Qqyhscu0h7",
				TenantID: "hpe-greenlake-intg",
				Expiry:   1566416348,
				IssuedAt: 1566412748,
			},
			wantErr: false,
		},
		{
			name: "Decode access token wrong token format",
			args: args{
				rawToken: "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJGd1BrQlJUbWY2Rm54R045SFp2VVhmOGhGa0xiOVI",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeAccessToken(tt.args.rawToken)
			if tt.wantErr {
				require.Error(t, err, "Error was expected, but decoding worked")
			} else {
				require.NoError(t, err, "Unexpected error")
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

// nolint: tparallel
func TestDoRetries(t *testing.T) {
	t.Parallel()
	totalRetries := 0
	testcases := []struct {
		name           string
		ctx            context.Context
		call           func(ctx context.Context) (*http.Request, *http.Response, error)
		responseStatus int
		err            error
	}{
		{
			name: "status 500",
			ctx:  context.Background(),
			call: func(ctx context.Context) (*http.Request, *http.Response, error) {
				totalRetries++

				return nil, &http.Response{StatusCode: http.StatusInternalServerError}, nil
			},
			responseStatus: http.StatusInternalServerError,
		},
		{
			name: "status 429",
			ctx:  context.Background(),
			call: func(ctx context.Context) (*http.Request, *http.Response, error) {
				totalRetries++

				return nil, &http.Response{StatusCode: http.StatusTooManyRequests}, nil
			},
			responseStatus: http.StatusTooManyRequests,
		},
		{
			name: "status 502",
			ctx:  context.Background(),
			call: func(ctx context.Context) (*http.Request, *http.Response, error) {
				totalRetries++

				return nil, &http.Response{StatusCode: http.StatusBadGateway}, nil
			},
			responseStatus: http.StatusBadGateway,
		},
		{
			name: "status 403 no retry",
			ctx:  context.Background(),
			call: func(ctx context.Context) (*http.Request, *http.Response, error) {
				return nil, &http.Response{StatusCode: http.StatusForbidden}, nil
			},
			responseStatus: http.StatusForbidden,
		},
		{
			name: "Deadline exceeded",
			ctx:  context.Background(),
			call: func(ctx context.Context) (*http.Request, *http.Response, error) {
				totalRetries++
				req := &http.Request{}
				req = req.WithContext(ctx)
				select {
				case <-ctx.Done():
					return req, nil, context.DeadlineExceeded
				case <-time.After(5 * time.Second): // this is greater than the deadline
					return req, nil, nil
				}
			},
			err: errLimitExceeded,
		},
		{
			name: "no url",
			ctx:  context.Background(),
			call: func(ctx context.Context) (*http.Request, *http.Response, error) {
				return nil, nil, errors.New("http: nil Request.URL")
			},
			err: errors.New("http: nil Request.URL"),
		},
	}

	for _, testcase := range testcases {
		tc := testcase
		t.Run(tc.name, func(t *testing.T) {
			cancelFuncs := make([]context.CancelFunc, 0)
			resp, err := DoRetries(tc.ctx, &cancelFuncs, tc.call, 2) // nolint: bodyclose
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
				if tc.err == errLimitExceeded {
					assert.Equal(t, 2, totalRetries)
					assert.Equal(t, 2, len(cancelFuncs))
				} else {
					assert.Equal(t, 0, totalRetries)
					assert.Equal(t, 1, len(cancelFuncs))
				}

				totalRetries = 0

			} else {
				assert.Equal(t, tc.responseStatus, resp.StatusCode)

				// only 429, 500 and 502 status codes should retry
				if tc.responseStatus == http.StatusForbidden {
					assert.Equal(t, 0, totalRetries)
					assert.Equal(t, 1, len(cancelFuncs))
				} else {
					assert.Equal(t, 2, totalRetries)
					assert.Equal(t, 2, len(cancelFuncs))
				}

				totalRetries = 0
			}
		})
	}
}
