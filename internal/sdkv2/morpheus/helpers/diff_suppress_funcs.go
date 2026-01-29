package helpers

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SuppressEquivalentJSONDiffs(k, old, new string, d *schema.ResourceData) bool {
	// First try simple whitespace-trimmed comparison
	if strings.TrimSpace(old) == strings.TrimSpace(new) {
		return true
	}

	// Trim whitespace before attempting JSON parsing
	// This handles cases where JSON has leading/trailing whitespace
	old = strings.TrimSpace(old)
	new = strings.TrimSpace(new)

	ob := bytes.NewBufferString("")
	if err := json.Compact(ob, []byte(old)); err != nil {
		return false
	}

	nb := bytes.NewBufferString("")
	if err := json.Compact(nb, []byte(new)); err != nil {
		return false
	}

	return jsonBytesEqual(ob.Bytes(), nb.Bytes())
}

func jsonBytesEqual(b1, b2 []byte) bool {
	var o1 any
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	var o2 any
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}
