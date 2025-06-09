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

func FindProjectRootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		logger.Printf("Error getting current directory: %v", err)

		return "."
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		} else if !os.IsNotExist(err) {
			logger.Printf("Error checking go.mod in %s: %v", dir, err)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			logger.Println("Reached filesystem root, go.mod not found")

			break
		}
		dir = parent
	}

	return "."
}

func RecordResult(t *testing.T) {
	if os.Getenv("RECORD_TEST_RESULTS") != "true" {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			TestResults[t.Name()] = TestResult{
				Status: "Failed",
				Error:  "panic: " + toString(r) + "\n\nstack trace:\n" + stack,
			}
			panic(r)
		} else if t.Failed() {
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
	rootOutputDir := filepath.Join(FindProjectRootDir(), "test_output")
	outputFile := filepath.Join(rootOutputDir, "result.json")

	existing := map[string]TestResult{}
	data, err := os.ReadFile(outputFile)
	if err == nil {
		if err := json.Unmarshal(data, &existing); err != nil {
			logger.Printf("Error unmarshaling result.json: %v", err)
		}
	} else if !os.IsNotExist(err) {
		logger.Printf("Error reading result.json: %v", err)
	}

	for k, v := range TestResults {
		existing[k] = v
	}

	output, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		logger.Printf("Error marshaling merged results: %v", err)

		return
	}

	if err := os.MkdirAll(rootOutputDir, 0o755); err != nil {
		logger.Printf("Error creating directory %s: %v", rootOutputDir, err)

		return
	}

	if err := os.WriteFile(outputFile, output, 0o600); err != nil {
		logger.Printf("Error writing result file: %v", err)
	}
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}

	return "unknown panic type"
}
