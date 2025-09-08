// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
)

const usage = `
Usage: 
go run render <type/name> <example_path> [args...]

Example:
go run render resources/hpe_morpheus_group group.tf.tmpl Name 'Test'
  `

func main() {
	if len(os.Args) < 2 || !strings.Contains(os.Args[1], "/") {
		fmt.Println(usage)
		os.Exit(1)
	}

	dest := os.Args[1]
	fn := os.Args[2]

	var args []string
	if len(os.Args) > 3 {
		args = os.Args[3:len(os.Args)]
	}

	testhelpers.WriteExample(dest, fn, args...)
}
