// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package main

import (
	"context"
	"flag"
	"log"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false,
		"set to true to run the provider with debugger support",
	)
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "github.com/hpe/hpe",
		Debug:   debug,
	}

	p := provider.New(
		version,
		// subprovider1.New(),
		// subprovider2.New(),
		// .
		// .
		// .
	)

	err := providerserver.Serve(context.Background(), p, opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
