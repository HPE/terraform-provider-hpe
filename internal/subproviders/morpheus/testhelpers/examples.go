// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"bytes"
	"fmt"
	"os"
	"path"
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

func WriteExample(fn string, args ...string) {
	text, err := renderExample(fn, args...)
	if err != nil {
		panic(err)
	}

	absPath, err := filepath.Abs(fn)
	if err != nil {
		panic(err)
	}

	pathParts := strings.Split(absPath, string(os.PathSeparator))
	if len(pathParts) < 3 {
		panic("Not enough path elements: " + absPath)
	}

	name := pathParts[len(pathParts)-2]
	kind := pathParts[len(pathParts)-3]

	exampleDir := map[string]string{
		"datasources": "examples/data-sources/hpe_morpheus_" + name,
		"resources":   "examples/resources/hpe_morpheus_" + name,
	}

	pathParts = pathParts[:len(pathParts)-1]

	var rootDir string

	for {
		rootDir = "/" + path.Join(pathParts...)
		if _, err := os.Stat(rootDir + "/.git"); err == nil {
			break
		}
		pathParts = pathParts[:len(pathParts)-1]
	}

	fn = strings.TrimSuffix(fn, ".tmpl")

	dest := path.Join(rootDir, exampleDir[kind], fn)

	err = os.WriteFile(dest, []byte(text), 0o644)
	if err != nil {
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
