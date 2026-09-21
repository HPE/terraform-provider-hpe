// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package optionsource

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// parseData turns a source's "data" into values.
//
// The controller renders whatever the source method returns, so the shape is
// per source. Two are handled:
//
//   - an array of options (most sources). Header rows — option-group labels
//     the UI shows, marked isGroup: true and carrying no value — are dropped.
//   - an object whose members include lists (zoneNetworkOptions: networks,
//     networkGroups, networkSubnets, networkTypes, plus scalars). collection
//     names the list to read; without it the error lists what is available,
//     so an object-shaped source never yields a silent empty result.
//
// Within an option, "value" is the value; when there is none, "id" stands in
// (zoneNetworkOptions items are {id: "network-79", name: ...}) and is also kept
// in Attributes. Every other key is stringified into Attributes.
func parseData(source string, data json.RawMessage, collection string) ([]Value, error) {
	trimmed := bytes.TrimSpace(data)

	switch {
	case len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")):
		if collection != "" {
			return nil, &CollectionError{Source: source, Collection: collection}
		}

		return []Value{}, nil

	case trimmed[0] == '[':
		if collection != "" {
			return nil, &ShapeError{Source: source, Collection: collection}
		}

		return parseArray(source, trimmed)

	case trimmed[0] == '{':
		var members map[string]json.RawMessage
		if err := json.Unmarshal(trimmed, &members); err != nil {
			return nil, fmt.Errorf("option source %q: decoding response object: %w", source, err)
		}

		available := make([]string, 0, len(members))

		for k, v := range members {
			if t := bytes.TrimSpace(v); len(t) > 0 && t[0] == '[' {
				available = append(available, k)
			}
		}

		sort.Strings(available)

		list, ok := members[collection]
		if collection == "" || !ok || len(bytes.TrimSpace(list)) == 0 || bytes.TrimSpace(list)[0] != '[' {
			return nil, &CollectionError{Source: source, Collection: collection, Available: available}
		}

		return parseArray(source, list)

	default:
		return nil, fmt.Errorf("option source %q: unexpected response data: %s", source, excerpt(trimmed))
	}
}

// parseArray decodes a list of options. Items are usually objects; a bare
// string is accepted as an option whose name and value are the string.
func parseArray(source string, list json.RawMessage) ([]Value, error) {
	dec := json.NewDecoder(bytes.NewReader(list))
	dec.UseNumber() // keep 1702 as "1702", not 1.702e+03

	var items []any
	if err := dec.Decode(&items); err != nil {
		return nil, fmt.Errorf("option source %q: decoding response list: %w", source, err)
	}

	values := make([]Value, 0, len(items))

	for _, item := range items {
		switch it := item.(type) {
		case string:
			values = append(values, Value{Name: it, Value: it, Attributes: map[string]string{}})
		case map[string]any:
			if isHeaderRow(it) {
				continue
			}

			values = append(values, valueFromObject(it))
		default:
			// Numbers, booleans and nested lists are not options; skip them.
		}
	}

	return values, nil
}

// isHeaderRow reports whether an item is an option-group label rather than an
// option. The UI marks these isGroup: true; they have a name and no value.
func isHeaderRow(item map[string]any) bool {
	switch g := item["isGroup"].(type) {
	case bool:
		return g
	case string:
		return strings.EqualFold(g, "true")
	default:
		return false
	}
}

// valueFromObject maps one option object to a Value.
func valueFromObject(item map[string]any) Value {
	v := Value{Attributes: make(map[string]string, len(item))}

	if s, ok := stringify(item["name"]); ok {
		v.Name = s
	}

	if s, ok := stringify(item["value"]); ok {
		v.Value = s
	} else if s, ok := stringify(item["id"]); ok {
		v.Value = s
	}

	for k, raw := range item {
		if k == "name" || k == "value" {
			continue
		}

		if s, ok := stringify(raw); ok {
			v.Attributes[k] = s
		}
	}

	return v
}

// stringify renders a decoded JSON value as the string the data source exposes.
// Strings pass through; numbers keep their textual form; booleans are
// "true"/"false"; objects and lists are compact JSON; null has no string.
func stringify(raw any) (string, bool) {
	switch x := raw.(type) {
	case nil:
		return "", false
	case string:
		return x, true
	case json.Number:
		return x.String(), true
	case bool:
		if x {
			return "true", true
		}

		return "false", true
	case float64:
		// Only reached if a caller decoded without UseNumber.
		return strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintf("%f", x), "000000"), "."), true
	default:
		b, err := json.Marshal(x)
		if err != nil {
			return "", false
		}

		return string(b), true
	}
}
