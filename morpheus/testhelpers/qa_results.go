package testhelpers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

type testResult struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

var (
	testResults         = make(map[string]testResult)
	m                   sync.Mutex
	_, recordingEnabled = os.LookupEnv("RECORD_TEST_RESULTS")
)

func RecordResult(t *testing.T) {
	if !recordingEnabled {
		return
	}

	m.Lock()
	if t.Failed() {
		testResults[t.Name()] = testResult{
			Status: "Failed",
			Error:  "Test " + t.Name() + "failed.",
		}
	} else if t.Skipped() {
		testResults[t.Name()] = testResult{
			Status: "Skipped",
			Error:  "",
		}
	} else {
		testResults[t.Name()] = testResult{
			Status: "Passed",
			Error:  "",
		}
	}
	m.Unlock()
}

func WriteMergedResults() {
	if !recordingEnabled {
		return
	}

	rootOutputDir := filepath.Join("/tmp", "test_output")
	outputFile := filepath.Join(rootOutputDir, acctest.RandomWithPrefix("test-result")+".result")

	// Marshal the merged results
	output, err := json.MarshalIndent(testResults, "", "  ")
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
