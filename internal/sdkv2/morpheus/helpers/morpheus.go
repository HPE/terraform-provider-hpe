// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package helpers

import (
	"os"
)

// global stuff here

var UseForce = (os.Getenv("MORPHEUS_USE_FORCE") == "true")
