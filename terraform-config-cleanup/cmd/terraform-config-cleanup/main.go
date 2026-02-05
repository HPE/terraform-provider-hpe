package main

import (
	"fmt"
	"os"

	"github.com/HPE/terraform-config-cleanup/pkg/config"
	"github.com/HPE/terraform-config-cleanup/pkg/parser"
	"github.com/HPE/terraform-config-cleanup/pkg/transform"
	"github.com/HPE/terraform-config-cleanup/pkg/writer"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "terraform-config-cleanup",
		Usage: "Clean up Terraform configuration generated via terraform import",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Path to directory containing YAML cleanup configs",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Usage:   "Input Terraform file to clean",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file path",
			},
			&cli.BoolFlag{
				Name:  "in-place",
				Usage: "Modify input file in place",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "Show what would be changed without modifying files",
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	configDir := c.String("config")
	inputFile := c.String("input")
	outputFile := c.String("output")
	inPlace := c.Bool("in-place")
	dryRun := c.Bool("dry-run")

	if inputFile == "" {
		return fmt.Errorf("--input is required")
	}

	if outputFile == "" && !inPlace && !dryRun {
		return fmt.Errorf("either --output, --in-place, or --dry-run must be specified")
	}

	// Determine output path
	outPath := outputFile
	if inPlace {
		outPath = inputFile
	}

	// 1. Load cleanup configurations
	fmt.Printf("Loading cleanup configs from: %s\n", configDir)
	configs, err := config.LoadConfigs(configDir)
	if err != nil {
		return fmt.Errorf("failed to load configs: %w", err)
	}
	fmt.Printf("Loaded %d cleanup configuration(s)\n", len(configs))

	// 2. Parse the input Terraform file
	fmt.Printf("Parsing input file: %s\n", inputFile)
	tfFile, err := parser.ParseFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to parse input file: %w", err)
	}
	fmt.Printf("Found %d resource(s)\n", len(tfFile.Resources))

	// 3. Apply transformations
	fmt.Println("Applying transformations...")
	for _, resource := range tfFile.Resources {
		cfg, exists := configs[resource.Type]
		if !exists {
			fmt.Printf("  - Skipping %s.%s (no cleanup config)\n", resource.Type, resource.Name)
			continue
		}

		fmt.Printf("  - Transforming %s.%s\n", resource.Type, resource.Name)
		if err := transform.ApplyTransformations(resource, cfg); err != nil {
			return fmt.Errorf("failed to transform %s.%s: %w", resource.Type, resource.Name, err)
		}
	}

	// 4. Write output (unless dry run)
	if !dryRun {
		fmt.Printf("Writing output to: %s\n", outPath)
		if err := writer.WriteFile(tfFile, outPath); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		fmt.Printf("\nSuccessfully cleaned up configuration and wrote to: %s\n", outPath)
	} else {
		fmt.Println("\nDry run - showing preview:")
		// In dry run mode, show what would be written
		tmpFile := "/tmp/terraform-config-cleanup-preview.tf"
		if err := writer.WriteFile(tfFile, tmpFile); err == nil {
			data, _ := os.ReadFile(tmpFile)
			fmt.Println(string(data))
			os.Remove(tmpFile)
		}
		fmt.Println("\nDry run completed - no changes were made")
	}

	return nil
}
