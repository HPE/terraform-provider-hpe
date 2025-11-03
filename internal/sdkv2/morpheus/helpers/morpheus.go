package helpers

import (
	"os"
)

// global stuff here

var UseForce = (os.Getenv("USE_FORCE") == "true")
