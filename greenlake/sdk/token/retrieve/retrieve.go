// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package retrieve

import (
	"context"

	"github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/common"
)

// TokenRetrieveFuncCtx type of function to retrieve a token passing-in a context
type TokenRetrieveFuncCtx func(ctx context.Context) (string, error)

// NewTokenRetrieveFunc takes a common.TokenChannelInterface as an input and returns a
// TokenRetrieveFuncCtx.  Exit from loop if a token is received on resCh, or if the
// context passed-in is cancelled.  On cancellation of context a signal is sent
// on exitCh to tell the token Handler retrieve thread to exit
func NewTokenRetrieveFunc(channelInterface common.TokenChannelInterface) TokenRetrieveFuncCtx {
	resCh, exitCh := channelInterface.TokenChannels()

	return func(ctx context.Context) (string, error) {
		for {
			select {
			case tok := <-resCh:
				return tok.Token, tok.Err
			case <-ctx.Done():
				exitCh <- 1

				return "", nil
			}
		}
	}
}
