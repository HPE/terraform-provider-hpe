// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

const usage = `
Usage: 
go run render <type/name> <example_path> [args...]

Examples:

Render to the same directory as the template:

go run render group.tf.tmpl Name 'Test'

Render to examples/resources/hpe_morpheus_group/group.tf:

go run render -out resources/hpe_morpheus_group group.tf.tmpl Name 'Test'
  `

func main() {
	out := flag.String("out", "", "Destination relative to .git directory")
	flag.Parse()

	osArgs := flag.Args()

	if len(osArgs) < 1 {
		fmt.Println(usage)
		os.Exit(1)
	}

	fn := osArgs[0]

	var args []string
	if len(osArgs) > 1 {
		args = osArgs[1:]
	}

	if *out != "" {
		testhelpers.WriteExampleToDir(*out, fn, args...)
	} else {
		testhelpers.WriteExample(fn, args...)
	}
}
