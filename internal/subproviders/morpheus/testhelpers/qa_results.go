package testhelpers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime/debug"
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
		stack := string(debug.Stack())
		TestResults[t.Name()] = TestResult{
			Status: "Failed",
			Error:  "test failed\n\nstack trace:\n" + stack,
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
			panic(err) // Panic only on errors other than file not existing
		}
	} else {
		// File exists, parse the JSON data
		if err := json.Unmarshal(data, &existing); err != nil {
			panic(err)
		}
	}

	// Merge new results into existing map
	for k, v := range TestResults {
		existing[k] = v
	}

	// Marshal the merged results
	output, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		panic(err)
	}

	// Ensure the output directory exists
	if err := os.MkdirAll(rootOutputDir, 0o755); err != nil {
		panic(err)
	}

	// Write the merged results to the file
	if err := os.WriteFile(outputFile, output, 0o600); err != nil {
		panic(err)
	}
}
