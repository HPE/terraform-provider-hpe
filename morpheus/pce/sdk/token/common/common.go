// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package common

const (
	// TimeToTokenExpiry is seconds in int64, not time.Second
	// This constant should be used in all handler code
	TimeToTokenExpiry = 120
)

// Result the result of a token retrieval
type Result struct {
	Token string
	Err   error
}
