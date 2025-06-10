package testhelpers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"
)

type TestResult struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

var (
	TestResults = make(map[string]TestResult)
	logger      = log.New(os.Stderr, "[testhelpers] ", log.LstdFlags)
)

func RecordResult(t *testing.T) {
	if os.Getenv("RECORD_TEST_RESULTS") != "true" {
		return
	}

	defer func() {
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
	}()
}

func WriteMergedResults() {
	rootOutputDir := filepath.Join("/tmp", "test_output")
	outputFile := filepath.Join(rootOutputDir, "result.json")

	existing := map[string]TestResult{}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Printf("result.json not found at %s; initializing fresh result set", outputFile)
			panic(err)
		}
		logger.Printf("Error reading result.json: %v", err)
		panic(err)
	}

	err = json.Unmarshal(data, &existing)
	if err != nil {
		logger.Printf("Error unmarshaling result.json: %v", err)
		panic(err)
	}

	for k, v := range TestResults {
		existing[k] = v
	}

	output, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		logger.Printf("Error marshaling merged results: %v", err)
		panic(err)
	}

	if err := os.MkdirAll(rootOutputDir, 0o755); err != nil {
		logger.Printf("Error creating directory %s: %v", rootOutputDir, err)
		panic(err)
	}

	if err := os.WriteFile(outputFile, output, 0o600); err != nil {
		logger.Printf("Error writing result file: %v", err)
		panic(err)
	}
}
