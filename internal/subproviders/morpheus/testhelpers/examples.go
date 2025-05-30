package testhelpers

import (
	"bytes"
	"os"
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
