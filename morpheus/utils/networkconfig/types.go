// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkconfig

// Context is where Morpheus reads an option type's value from and stores it:
// the network object itself, or its free-form config map.
type Context string

const (
	// ContextNetwork fields are first-class properties of the network, sent at
	// the top level of the request and, for most of them, returned by the API.
	ContextNetwork Context = "network"
	// ContextConfig fields live in the network's config map: sent under
	// "config" and returned under "config", masked where secret.
	ContextConfig Context = "config"
)

// FieldSpec is one option type of a network type, as the appliance declares it.
type FieldSpec struct {
	// Name is the field name Morpheus reads the value under, e.g.
	// "hvmVlanTrunks", or the dotted path of a related object, e.g. "zonePool.id".
	Name string
	// Context is where the value lives.
	Context Context
	// Input is the Morpheus input type: "text", "number", "checkbox", "select",
	// "typeahead", "textarea" or "hidden".
	Input string
	// Label is the field's UI label, for messages.
	Label string
	// Required reports whether Morpheus requires a value on create.
	Required bool
	// Creatable reports whether Morpheus applies the field on create.
	Creatable bool
	// Editable reports whether Morpheus applies a new value on update. A change
	// to a field that is not editable is accepted and ignored without error,
	// so the provider must replace the network instead.
	Editable bool
	// Default is the value Morpheus applies when the field is omitted, as the
	// UI would send it ("on"/"off" for a checkbox).
	Default string
	// OptionSource names the option source that lists the field's values.
	OptionSource string
	// Min and Max bound a number field, when the type declares them.
	Min, Max string
}

// TypeFlags are a network type's own switches for the core fields Morpheus
// still edits after create. A false flag means an update that changes the
// field is accepted and ignored, so a change requires replacement.
type TypeFlags struct {
	Name, Cidr, Gateway, VlanID, DhcpServer, DNS, StaticOverride, NetworkDomain bool
	// CanAssignPool reports whether a network pool can be assigned to the type.
	CanAssignPool bool
}

// TypeSpec is one network type as the appliance declares it.
type TypeSpec struct {
	// Code is the network type's code, e.g. "ovsPortGroup".
	Code string
	// Name is the type's display name, e.g. "OVS Port Group".
	Name string
	// Category groups related types, e.g. "ovs", "nsxt"; empty for many.
	Category string
	// ExternalType is the type's external classification; empty for many.
	ExternalType string
	// Creatable reports whether networks of the type can be created through
	// the API rather than only discovered by a cloud sync.
	Creatable bool
	// Deletable reports whether networks of the type can be deleted.
	Deletable bool
	// HasNetworkServer reports whether a network of the type belongs to a
	// network server (integration), which Morpheus 9.1 requires on create.
	HasNetworkServer bool
	// HasCidr and CidrRequired describe the type's use of a CIDR.
	HasCidr, CidrRequired bool
	// Overlay reports whether the type is an overlay network.
	Overlay bool
	// HasStaticRoutes reports whether the type carries static routes.
	HasStaticRoutes bool
	// Flags are the type's editability switches for the core fields.
	Flags TypeFlags
	// Fields are the type's option types, network context first, then config
	// context, each in display order.
	Fields []FieldSpec
	// Live is true when the spec was read from the appliance the provider is
	// talking to, and false for an entry of the generated table. Only a live
	// spec says what that appliance's release and plugins actually declare.
	Live bool
}

// Lookup returns the generated table's entry for a network type code. The
// resource prefers the live shape from the appliance (FromAPI) and uses this
// when the appliance cannot be read; the block-declaration tests use it as
// their fixture.
func Lookup(code string) (TypeSpec, bool) {
	t, ok := knownTypes[code]

	return t, ok
}

// Codes returns every network type code in the generated table, unordered.
func Codes() []string {
	codes := make([]string, 0, len(knownTypes))
	for code := range knownTypes {
		codes = append(codes, code)
	}

	return codes
}

// Field returns the type's option type named name, in either context.
func (t TypeSpec) Field(name string) (FieldSpec, bool) {
	for _, f := range t.Fields {
		if f.Name == name {
			return f, true
		}
	}

	return FieldSpec{}, false
}

// ConfigFields returns the type's config-context option types, in order.
func (t TypeSpec) ConfigFields() []FieldSpec {
	return t.fields(ContextConfig)
}

// NetworkFields returns the type's network-context option types, in order.
func (t TypeSpec) NetworkFields() []FieldSpec {
	return t.fields(ContextNetwork)
}

func (t TypeSpec) fields(c Context) []FieldSpec {
	var out []FieldSpec

	for _, f := range t.Fields {
		if f.Context == c {
			out = append(out, f)
		}
	}

	return out
}
