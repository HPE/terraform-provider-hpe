// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package serviceclient

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/common"
	httpc "github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/httpclient"
	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/iamversion"
	tokenutil "github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/token-util"
)

const retryLimit = 3

//go:generate go run github.com/golang/mock/mockgen -build_flags=-mod=mod -destination=../mocks/IdentityAPI_mocks.go -package=mocks github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/serviceclient IdentityAPI
type IdentityAPI interface {
	GenerateToken(context.Context, string, string, string, string) (string, error)
}

// Handler the handler for service-client creds
type Handler struct {
	iamServiceURL       string
	token               string
	passedInToken       string
	tenantID            string
	clientID            string
	clientSecret        string
	iamVersion          string
	vendedServiceClient bool
	numRetries          int
	client              IdentityAPI
}

// CreateOpt - function option definition
type CreateOpt func(h *Handler)

// WithIdentityAPI override the IdentityAPI in Handler
func WithIdentityAPI(i IdentityAPI) CreateOpt {
	return func(h *Handler) {
		h.client = i
	}
}

// WithIAMServiceURL sets the GreenLake IAM issuer/service URL used to generate
// access tokens.
func WithIAMServiceURL(url string) CreateOpt {
	return func(h *Handler) {
		h.iamServiceURL = url
	}
}

// WithClientID sets the API client ID used for authentication.
func WithClientID(clientID string) CreateOpt {
	return func(h *Handler) {
		h.clientID = clientID
	}
}

// WithClientSecret sets the API client secret used for authentication.
func WithClientSecret(clientSecret string) CreateOpt {
	return func(h *Handler) {
		h.clientSecret = clientSecret
	}
}

// WithIAMToken sets a pre-generated IAM token. When supplied, token generation
// from credentials is skipped and this token is returned as-is.
func WithIAMToken(token string) CreateOpt {
	return func(h *Handler) {
		h.passedInToken = token
	}
}

// WithIAMVersion sets the GreenLake IAM flavour (GLCS/GLP) used for the token
// exchange. Defaults to GLCS when not supplied.
func WithIAMVersion(v iamversion.Version) CreateOpt {
	return func(h *Handler) {
		h.iamVersion = string(v)
	}
}

// WithTenantID sets the tenant ID. This is only consumed by the legacy
// non-vended identity-token path and is retained for legacy code support.
func WithTenantID(tenantID string) CreateOpt {
	return func(h *Handler) {
		h.tenantID = tenantID
	}
}

// WithAPIVendedServiceClient overrides whether the API client is API-vended.
// Handlers default to vended (true); this is retained for legacy code support.
func WithAPIVendedServiceClient(vended bool) CreateOpt {
	return func(h *Handler) {
		h.vendedServiceClient = vended
	}
}

// NewHandler creates a new token handler. All configuration is supplied via
// CreateOpt options, keeping the handler decoupled from any provider/config
// framework.
func NewHandler(opts ...CreateOpt) (*Handler, error) {
	h := new(Handler)

	// Defaults; may be overridden via opts (retained for legacy code support).
	h.vendedServiceClient = true
	h.iamVersion = string(iamversion.GLCS)

	// run overrides
	for _, opt := range opts {
		if opt != nil {
			opt(h)
		}
	}

	// Build the real client only if one wasn't injected via an opt
	// (e.g. WithIdentityAPI in tests). Opts are applied first so that
	// WithAPIVendedServiceClient is reflected in the client.
	if h.client == nil {
		h.client = httpc.New(h.iamServiceURL, h.vendedServiceClient, h.passedInToken)
	}

	return h, nil
}

// Token returns an IAM token for the handler's configuration, generating one if
// the handler does not already hold a usable one.
//
// The work is done synchronously on the calling goroutine, so a caller needing
// a single token pays for exactly one exchange and leaves nothing running
// behind it. The handler is not safe for concurrent use.
func (h *Handler) Token(ctx context.Context) (string, error) {
	res := h.retrieveToken(ctx)

	return res.Token, res.Err
}

// retrieveToken function to get a token
// The token is stashed in the handler.  If its time-to-expiry is <= common.TimeToTokenExpiry then it is
// regenerated.
// If we have to regenerate a token we will retry in the case where the error is retryable up to retryLimit times
// Currently the only error that is retryable is a net Timeout error
func (h *Handler) retrieveToken(ctx context.Context) common.Result {
	// We use a loop since we may need to retry depending on the error that we get from IAM
	// Reset numRetries
	h.numRetries = 0
	for {
		// Stop before spending a retry on a call that has been cancelled.
		if err := ctx.Err(); err != nil {
			return common.Result{
				Token: "",
				Err:   err,
			}
		}

		// Get current time in Unix "epoch" seconds
		now := time.Now().Unix()

		// Generate token if there isn't any
		if h.token == "" {
			token, retry, err := h.generateToken(ctx)
			if retry {
				continue
			}

			if err != nil {
				return common.Result{
					Token: "",
					Err:   err,
				}
			}

			h.token = token
		}

		// Decode token
		tokenDetails, err := tokenutil.DecodeAccessToken(h.token)
		if err != nil {
			return common.Result{
				Token: "",
				Err:   err,
			}
		}

		// If token is about to expire in TimeToTokenExpiry seconds or less generate a new one
		if tokenDetails.Expiry-now <= common.TimeToTokenExpiry {
			token, retry, err := h.generateToken(ctx)
			if retry {
				continue
			}

			if err != nil {
				return common.Result{
					Token: "",
					Err:   err,
				}
			}

			h.token = token
		}

		return common.Result{
			Token: h.token,
			Err:   nil,
		}
	}
}

// generateToken simple function to call the API client's GenerateToken
func (h *Handler) generateToken(ctx context.Context) (string, bool, error) {
	var token string
	var err error

	token, err = h.client.GenerateToken(ctx, h.tenantID, h.clientID, h.clientSecret, h.iamVersion)

	// If this is a retryable error check to see if we've reached our retryLimit or not, if we can retry again
	// return true
	if err != nil && isErrRetryable(err) {
		h.numRetries++

		return token, h.numRetries <= retryLimit, err
	}

	return token, false, err
}

// isErrRetryable checks if an error is retryable, currently limited to net Timeout errors
func isErrRetryable(err error) bool {
	var t net.Error
	if errors.As(err, &t) && t.Timeout() {
		return true
	}

	return false
}
