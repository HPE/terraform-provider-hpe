// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkconfig

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// The provider learns a network type's shape from the appliance it is talking
// to, at plan and apply time, and falls back to the generated table when it
// cannot. FromAPI is the live half: it normalises a network type as
// GET /api/network-types returns it into the same TypeSpec the harvester
// generates, by the same rules, so the two are interchangeable and a test can
// hold them to that.

// apiNetworkType is the wire shape of a network type, as far as the provider
// reads it. It is decoded from JSON rather than from the SDK's structs because
// the list and show responses use different structs for the same object, and
// because the harvester reads the same JSON.
type apiNetworkType struct {
	Code                   string          `json:"code"`
	Name                   string          `json:"name"`
	Category               string          `json:"category"`
	ExternalType           string          `json:"externalType"`
	Creatable              bool            `json:"creatable"`
	Deletable              bool            `json:"deletable"`
	HasNetworkServer       bool            `json:"hasNetworkServer"`
	HasCidr                bool            `json:"hasCidr"`
	CidrRequired           bool            `json:"cidrRequired"`
	Overlay                bool            `json:"overlay"`
	HasStaticRoutes        bool            `json:"hasStaticRoutes"`
	NameEditable           bool            `json:"nameEditable"`
	CidrEditable           bool            `json:"cidrEditable"`
	GatewayEditable        bool            `json:"gatewayEditable"`
	VlanIDEditable         bool            `json:"vlanIdEditable"`
	DhcpServerEditable     bool            `json:"dhcpServerEditable"`
	DNSEditable            bool            `json:"dnsEditable"`
	StaticOverrideEditable bool            `json:"staticOverrideEditable"`
	NetworkDomainEditable  bool            `json:"networkDomainEditable"`
	CanAssignPool          bool            `json:"canAssignPool"`
	OptionTypes            []apiOptionType `json:"optionTypes"`
}

// apiOptionType is the wire shape of an option type. Booleans that Morpheus
// treats as true when absent are pointers, so absence can be told from false.
type apiOptionType struct {
	FieldName    string `json:"fieldName"`
	FieldContext string `json:"fieldContext"`
	Type         string `json:"type"`
	FieldLabel   string `json:"fieldLabel"`
	Required     bool   `json:"required"`
	Creatable    *bool  `json:"creatable"`
	Editable     *bool  `json:"editable"`
	DefaultValue any    `json:"defaultValue"`
	OptionSource any    `json:"optionSource"`
	MinVal       any    `json:"minVal"`
	MaxVal       any    `json:"maxVal"`
	DisplayOrder int64  `json:"displayOrder"`
}

// FromAPI normalises a network type as the API returns it — an SDK response
// struct for the show or list endpoint, or the decoded JSON object — into a
// TypeSpec. It fails only when the value cannot be read as a network type or
// has no code.
func FromAPI(networkType any) (TypeSpec, error) {
	raw, err := json.Marshal(networkType)
	if err != nil {
		return TypeSpec{}, fmt.Errorf("network type: encode: %w", err)
	}

	var t apiNetworkType
	if err := json.Unmarshal(raw, &t); err != nil {
		return TypeSpec{}, fmt.Errorf("network type: decode: %w", err)
	}

	if t.Code == "" {
		return TypeSpec{}, fmt.Errorf("network type: the response has no code")
	}

	spec := TypeSpec{
		Code:             t.Code,
		Name:             strings.TrimSpace(t.Name),
		Category:         strings.TrimSpace(t.Category),
		ExternalType:     strings.TrimSpace(t.ExternalType),
		Creatable:        t.Creatable,
		Deletable:        t.Deletable,
		HasNetworkServer: t.HasNetworkServer,
		HasCidr:          t.HasCidr,
		CidrRequired:     t.CidrRequired,
		Overlay:          t.Overlay,
		HasStaticRoutes:  t.HasStaticRoutes,
		Flags: TypeFlags{
			Name: t.NameEditable, Cidr: t.CidrEditable, Gateway: t.GatewayEditable, VlanID: t.VlanIDEditable,
			DhcpServer: t.DhcpServerEditable, DNS: t.DNSEditable, StaticOverride: t.StaticOverrideEditable,
			NetworkDomain: t.NetworkDomainEditable, CanAssignPool: t.CanAssignPool,
		},
	}

	type ordered struct {
		spec  FieldSpec
		order int64
	}

	var fields []ordered

	for _, o := range t.OptionTypes {
		if o.FieldName == "" {
			continue
		}

		context := ContextNetwork
		if strings.TrimSpace(o.FieldContext) == string(ContextConfig) {
			context = ContextConfig
		}

		input := o.Type
		if input == "" {
			input = "text"
		}

		fields = append(fields, ordered{order: o.DisplayOrder, spec: FieldSpec{
			Name:         o.FieldName,
			Context:      context,
			Input:        input,
			Label:        strings.TrimSpace(o.FieldLabel),
			Required:     o.Required,
			Creatable:    o.Creatable == nil || *o.Creatable,
			Editable:     o.Editable == nil || *o.Editable,
			Default:      text(o.DefaultValue),
			OptionSource: text(o.OptionSource),
			Min:          text(o.MinVal),
			Max:          text(o.MaxVal),
		}})
	}

	// Network context first, then display order, then name — as generated.
	sort.SliceStable(fields, func(i, j int) bool {
		a, b := fields[i], fields[j]

		if a.spec.Context != b.spec.Context {
			return a.spec.Context == ContextNetwork
		}

		if a.order != b.order {
			return a.order < b.order
		}

		return a.spec.Name < b.spec.Name
	})

	if len(fields) > 0 {
		spec.Fields = make([]FieldSpec, 0, len(fields))
		for _, f := range fields {
			spec.Fields = append(spec.Fields, f.spec)
		}
	}

	return spec, nil
}

// text renders a loosely typed API value as the harvester does: null as "",
// an integral number without a fraction, anything else as its string.
func text(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(x)
	case bool:
		return strconv.FormatBool(x)
	case float64:
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10)
		}

		return strconv.FormatFloat(x, 'f', -1, 64)
	case json.Number:
		return x.String()
	default:
		return strings.TrimSpace(fmt.Sprint(x))
	}
}
