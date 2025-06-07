package testhelpers

import (
	"encoding/json"
	"fmt"
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

func FindProjectRootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "." // fallback
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			// Found go.mod here, this is project root
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root without finding go.mod
			break
		}
		dir = parent
	}

	// Fallback if no go.mod found
	return "."
}

func RecordResult(t *testing.T) {
	if os.Getenv("RECORD_TEST_RESULTS") != "true" {
		return
	}
	if t.Skipped() {
		TestResults[t.Name()] = struct {
			Status string `json:"status"`
			Error  string `json:"error"`
		}{"Skipped", ""}
		return
	}

	if r := recover(); r != nil {
		stack := string(debug.Stack())
		TestResults[t.Name()] = struct {
			Status string `json:"status"`
			Error  string `json:"error"`
		}{
			"Failed",
			fmt.Sprintf("panic: %v\n\nstack trace:\n%s", r, stack),
		}
		panic(r)
	} else if t.Failed() {
		stack := string(debug.Stack())
		TestResults[t.Name()] = struct {
			Status string `json:"status"`
			Error  string `json:"error"`
		}{
			"Failed",
			fmt.Sprintf("test failed\n\nstack trace:\n%s", stack),
		}
	} else {
		TestResults[t.Name()] = struct {
			Status string `json:"status"`
			Error  string `json:"error"`
		}{"Passed", ""}
	}
}

func WriteMergedResults() {
	rootOutputDir := filepath.Join(FindProjectRootDir(), "test_output")
	outputFile := filepath.Join(rootOutputDir, "result.json")

	// Load existing results if any
	existing := map[string]TestResult{}
	if data, err := os.ReadFile(outputFile); err == nil {
		_ = json.Unmarshal(data, &existing)
	}

	// Merge new results
	for k, v := range TestResults {
		existing[k] = v
	}

	// Write back
	output, err := json.MarshalIndent(existing, "", "  ")
	if err == nil {
		if err := os.MkdirAll(rootOutputDir, 0755); err == nil {
			_ = os.WriteFile(outputFile, output, 0644)
		}
	}
	fmt.Println("Test results written")
}
