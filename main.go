// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"

	"github.com/HPE/terraform-provider-hpe/internal/framework/provider"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
)

var version = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false,
		"set to true to run the provider with debugger support",
	)
	flag.Parse()

	p := provider.New(
		version,
		morpheus.New(),
		// subprovider2.New(),
		// subprovider3.New(),
		// .
		// .
		// .
	)

	var opts []tf6server.ServeOpt
	if debug {
		homeDir, _ := os.UserHomeDir()
		config := filepath.Join(homeDir, ".config", "terraform-provider-hpe")
		configFile := filepath.Join(config, "debug.env")

		if err := os.MkdirAll(config, os.ModePerm); err != nil {
			log.Fatal("could not create a debug config folder: ", err)
		}

		opts = append(opts,
			tf6server.WithManagedDebug(),
			tf6server.WithManagedDebugEnvFilePath(configFile),
		)
	}

	if err := tf6server.Serve(
		"registry.terraform.io/HPE/hpe",
		providerserver.NewProtocol6(p()),
		opts...,
	); err != nil {
		log.Fatal(err.Error())
	}
}
