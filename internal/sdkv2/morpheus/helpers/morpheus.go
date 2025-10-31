package helpers

import (
	"os"
)

// global stuff here

var (
	USE_FORCE = (os.Getenv("USE_FORCE") == "true")
)
