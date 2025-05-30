package testhelpers

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"text/template"
)

func RenderExample(t *testing.T, name string, args ...string) string {
	t.Helper()

	if len(args)%2 != 0 {
		t.Fatal(`arguments must be space separated pairs in the format "Key" "value"`)
	}

	bs, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}

	tmpl := template.New(name)
	tmpl, err = tmpl.Parse(string(bs))
	if err != nil {
		t.Fatalf("unable to parse template %q: %s", bs, err.Error())
	}

	data := make(map[string]string)
	for i := 0; i < len(args)-1; i += 2 {
		data[args[i]] = args[i+1]
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		t.Fatalf("unable to execute template: %s", err.Error())
	}

	return buf.String()
}

func RenderExample(t *testing.T, name string, args ...string) string {
	t.Helper()

	example, err := renderExample(name, args...)
	if err != nil {
		t.Fatal(err)
	}

	return example
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
