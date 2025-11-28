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

	fn = strings.TrimSuffix(fn, ".tmpl")
	dest = filepath.Join(absPath, "examples", dest, fn)

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
