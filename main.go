// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

	morpheus "github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2Morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/provider"
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

	// sdkv2 Morpheus provider server
	legacyMorpheus, err := tf5to6server.UpgradeServer(context.Background(), sdkv2Morpheus.Provider().GRPCProvider)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Combine framework and sdkv2 providers
	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(p()),
		func() tfprotov6.ProviderServer { return legacyMorpheus },
	}

	muxServer, err := tf6muxserver.NewMuxServer(context.Background(), providers...)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := tf6server.Serve(
		"registry.terraform.io/HPE/hpe",
		muxServer.ProviderServer,
		opts...,
	); err != nil {
		log.Fatal(err.Error())
	}
}
