// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"text/template"
)

func ReadExample(t *testing.T, name, rgx, replace string) (string, error) {
	t.Helper()

	bytes, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}

	rg := regexp.MustCompile(rgx)

	example := rg.ReplaceAllString(string(bytes), replace)

	return example, nil
}

func RenderExample(t *testing.T, name string, args ...string) (string, error) {
	t.Helper()

	example, err := renderExample(name, args...)
	if err != nil {
		return "", err
	}

	return example, nil
}

// MergeOverrides applies overrides on top of defaults, ignoring any override
// whose value is empty.
//
// Test fixtures routinely build overrides straight from os.Getenv. Assigning
// unconditionally lets an unset variable replace a working default with "",
// which the template then renders as nothing -- turning `network_server_id = 1`
// into `network_server_id = `, an HCL parse error that reads as a provider bug.
// Skipping empty overrides keeps the default in play; a test that genuinely
// requires the variable should check it and t.Skip explicitly.
func MergeOverrides(defaults, overrides map[string]string) map[string]string {
	merged := make(map[string]string, len(defaults))
	for key, value := range defaults {
		merged[key] = value
	}

	for key, value := range overrides {
		if value == "" {
			continue
		}

		merged[key] = value
	}

	return merged
}

// RenderArgs flattens a map into the "Key" "value" argument pairs that
// RenderExample and WriteExample expect.
func RenderArgs(values map[string]string) []string {
	args := make([]string, 0, len(values)*2)
	for key, value := range values {
		args = append(args, key, value)
	}

	return args
}

func WriteExample(name string, args ...string) {
	text, err := renderExample(name, args...)
	if err != nil {
		panic(err)
	}

	name = strings.TrimSuffix(name, ".tmpl")

	err = os.WriteFile(name, []byte(text), 0o600)
	if err != nil {
		panic(err)
	}
}

func WriteExampleToDir(dest, fn string, args ...string) {
	text, err := renderExample(fn, args...)
	if err != nil {
		panic(err)
	}

	absPath, err := filepath.Abs(fn)
	if err != nil {
		panic(err)
	}
	absPath = filepath.Dir(absPath)

	for {
		if _, err := os.Stat(filepath.Join(absPath, ".git")); err == nil {
			break
		}
		absPath = filepath.Dir(absPath)
	}

	dest = filepath.Join(absPath, dest)

	if err := os.MkdirAll(filepath.Dir(dest), 0o700); err != nil {
		panic(err)
	}

	if err = os.WriteFile(dest, []byte(text), 0o600); err != nil {
		panic(err)
	}
}

func renderExample(name string, args ...string) (string, error) {
	if len(args)%2 != 0 {
		return "", fmt.Errorf(`arguments must be space separated pairs in the format "Key" "value"`)
	}

	bs, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}

	tmpl := template.New(name)
	tmpl, err = tmpl.Parse(string(bs))
	if err != nil {
		return "", fmt.Errorf("unable to parse template %q: %w", bs, err)
	}

	data := make(map[string]string)
	for i := 0; i < len(args)-1; i += 2 {
		data[args[i]] = args[i+1]
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("unable to execute template: %w", err)
	}

	return buf.String(), nil
}
