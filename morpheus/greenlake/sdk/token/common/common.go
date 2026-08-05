// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package common

const (
	TokenRetrieveFunctionKey = "tokenRetrieveFunc"
	// TimeToTokenExpiry is seconds in int64, not time.Second
	// This constant should be used in all handler code
	TimeToTokenExpiry = 120
)

// Result the result struct sent back on the resultCh of a token Handler
type Result struct {
	Token string
	Err   error
}

// TokenChannelInterface the interface that is implemented by a token Handler
// This interface is used in retrieve.NewTokenRetrieveFunc
type TokenChannelInterface interface {
	TokenChannels() (chan Result, chan int)
}
