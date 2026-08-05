// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package iamversion defines the GreenLake IAM token-exchange flavours shared
// by the token SDK. The token client is version-agnostic: callers pass the
// desired Version at call time, so the same client serves both GLCS and GLP
// exchanges.
package iamversion

// Version identifies a GreenLake IAM token-exchange flavour.
type Version string

const (
	// GLCS is the IAM version for GreenLake Cloud Services.
	GLCS Version = "glcs"
	// GLP is the IAM version for the GreenLake Platform.
	GLP Version = "glp"
)
