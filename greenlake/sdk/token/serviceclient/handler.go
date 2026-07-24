// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package serviceclient

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/common"
	httpc "github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/httpclient"
	"github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/iamversion"
	tokenutil "github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/token-util"
)

const retryLimit = 3

// Assert that Handler implements common.TokenChannelInterface
var _ common.TokenChannelInterface = (*Handler)(nil)

//go:generate go run github.com/golang/mock/mockgen -build_flags=-mod=mod -destination=../mocks/IdentityAPI_mocks.go -package=mocks github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/serviceclient IdentityAPI
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
	resultCh            chan common.Result
	exitCh              chan int
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

// NewHandler creates a new handler and returns the common.TokenChannelInterface
// interface. All configuration is supplied via CreateOpt options, keeping the
// handler decoupled from any provider/config framework.
func NewHandler(opts ...CreateOpt) (common.TokenChannelInterface, error) {
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

	// make channels
	h.resultCh = make(chan common.Result)
	h.exitCh = make(chan int)

	// set-up retrieve thread on channel
	h.startRetrieveThread()

	return h, nil
}

// TokenChannels return channels for token retrieve function
func (h *Handler) TokenChannels() (chan common.Result, chan int) { // nolint golint
	return h.resultCh, h.exitCh
}

// startRetrieveThread start the token retrieve thread
// function in an infinite loop, it puts the return value of retrieveToken into h.resultCh by default
// if a signal on exitCh is received the thread exits
func (h *Handler) startRetrieveThread() {
	go func() {
		for {
			select {
			case <-h.exitCh:
				// TODO we need to set-up a context here and cancel it so that the TokenGenerate call is killed
				return
			default:
				h.resultCh <- h.retrieveToken()
			}
		}
	}()
}

// retrieveToken function to get a token
// The token is stashed in the handler.  If its time-to-expiry is <= common.TimeToTokenExpiry then it is
// regenerated.
// If we have to regenerate a token we will retry in the case where the error is retryable up to retryLimit times
// Currently the only error that is retryable is a net Timeout error
func (h *Handler) retrieveToken() common.Result {
	// We use a loop since we may need to retry depending on the error that we get from IAM
	// Reset numRetries
	h.numRetries = 0
	for {
		// Get current time in Unix "epoch" seconds
		now := time.Now().Unix()

		// Generate token if there isn't any
		if h.token == "" {
			token, retry, err := h.generateToken()
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
			token, retry, err := h.generateToken()
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
func (h *Handler) generateToken() (string, bool, error) {
	var token string
	var err error

	// TODO pass a context down to here
	token, err = h.client.GenerateToken(context.Background(), h.tenantID, h.clientID, h.clientSecret, h.iamVersion)

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
