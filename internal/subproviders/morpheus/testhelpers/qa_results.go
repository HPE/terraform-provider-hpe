package testhelpers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"testing"
)

type TestResult struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

var TestResults = make(map[string]TestResult)

func RecordResult(t *testing.T) {
	if os.Getenv("RECORD_TEST_RESULTS") != "true" {
		return
	}

	if t.Failed() {
		TestResults[t.Name()] = TestResult{
			Status: "Failed",
			Error:  "Test " + t.Name() + "failed.",
		}
	} else if t.Skipped() {
		TestResults[t.Name()] = TestResult{
			Status: "Skipped",
			Error:  "",
		}
	} else {
		TestResults[t.Name()] = TestResult{
			Status: "Passed",
			Error:  "",
		}
	}
}

func WriteMergedResults() {
	rootOutputDir := filepath.Join("/tmp", "test_output")
	outputFile := filepath.Join(rootOutputDir, "result.json")

	existing := map[string]TestResult{}

	// Try to read the existing file
	data, err := os.ReadFile(outputFile)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Error reading existing results file: %v", err)
			os.Exit(1)
		}
	} else {
		// File exists, parse the JSON data
		if err := json.Unmarshal(data, &existing); err != nil {
			log.Printf("Error parsing JSON from existing results file: %v", err)
			os.Exit(1)
		}
	}

	// Merge new results into existing map
	for k, v := range TestResults {
		existing[k] = v
	}

	// Marshal the merged results
	output, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		log.Printf("Error marshalling merged results: %v", err)
		os.Exit(1)
	}

	// Ensure the output directory exists
	if err := os.MkdirAll(rootOutputDir, 0o755); err != nil {
		log.Printf("Error creating output directory: %v", err)
		os.Exit(1)
	}

	// Write the merged results to the file
	if err := os.WriteFile(outputFile, output, 0o600); err != nil {
		log.Printf("Error writing merged results to file: %v", err)
		os.Exit(1)
	}
}
