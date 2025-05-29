package testhelpers

import (
	"os"
	"regexp"
	"testing"
)

func ReadExample(t *testing.T, name, rgx, replace string) string {
	t.Helper()

	bytes, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}

	rg := regexp.MustCompile(rgx)

	example := rg.ReplaceAllString(string(bytes), replace)

	return example
}
