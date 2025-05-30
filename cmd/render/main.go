package main

import (
	"os"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
)

func main() {
	name := os.Args[1]
	args := os.Args[2:len(os.Args)]

	testhelpers.WriteExample(name, args...)
}
