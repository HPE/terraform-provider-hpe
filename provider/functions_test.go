package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestUnitSHA512CryptMux(t *testing.T) {
	ctx := context.Background()
	factory := testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider())["hpe"]
	server, err := factory()
	if err != nil {
		t.Fatal(err)
	}
	schema, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for _, diagnostic := range schema.Diagnostics {
		if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
			t.Fatal(diagnostic.Summary)
		}
	}
	definitions, err := server.GetFunctions(ctx, &tfprotov6.GetFunctionsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	definition := definitions.Functions["morpheus_sha512_crypt"]
	if definition == nil || len(definition.Parameters) != 3 || !definition.Return.Type.Equal(tftypes.String) {
		t.Fatal("missing or incorrect namespaced function definition")
	}
	for _, parameter := range definition.Parameters {
		if parameter.AllowNullValue || parameter.AllowUnknownValues {
			t.Fatal("null and unknown inputs must be handled by the framework")
		}
	}
	cases := []struct {
		name     string
		password any
		salt     any
		rounds   any
		unknown  bool
		invalid  bool
	}{
		{"known", "password", "saltsalt", int64(5000), false, false},
		{"unknown_password", tftypes.UnknownValue, "saltsalt", int64(5000), true, false},
		{"unknown_salt", "password", tftypes.UnknownValue, int64(5000), true, false},
		{"unknown_rounds", "password", "saltsalt", tftypes.UnknownValue, true, false},
		{"null_password", nil, "saltsalt", int64(5000), false, true},
		{"null_salt", "password", nil, int64(5000), false, true},
		{"null_rounds", "password", "saltsalt", nil, false, true},
		{"invalid_salt", "password", "$", int64(5000), false, true},
		{"invalid_rounds", "password", "saltsalt", int64(999), false, true},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			values := []tftypes.Value{
				tftypes.NewValue(tftypes.String, test.password),
				tftypes.NewValue(tftypes.String, test.salt),
				tftypes.NewValue(tftypes.Number, test.rounds),
			}
			arguments := make([]*tfprotov6.DynamicValue, len(values))
			for index, value := range values {
				dynamic, err := tfprotov6.NewDynamicValue(value.Type(), value)
				if err != nil {
					t.Fatal(err)
				}
				arguments[index] = &dynamic
			}
			response, err := server.CallFunction(ctx, &tfprotov6.CallFunctionRequest{
				Name: "morpheus_sha512_crypt", Arguments: arguments,
			})
			if err != nil {
				t.Fatal(err)
			}
			if test.invalid {
				if response.Error == nil {
					t.Fatal("expected function error")
				}

				return
			}
			if response.Error != nil {
				t.Fatal(response.Error)
			}
			result, err := response.Result.Unmarshal(tftypes.String)
			if err != nil {
				t.Fatal(err)
			}
			if test.unknown {
				if result.IsKnown() {
					t.Fatal("unknown input must yield unknown result")
				}

				return
			}
			var actual string
			if err := result.As(&actual); err != nil {
				t.Fatal(err)
			}
			if actual != "$6$saltsalt$qFmFH.bQmmtXzyBY0s9v7Oicd2z4XSIecDzlB5KiA2/jctKu9YterLp8wwnSq.qc.eoxqOmSuNp2xS0ktL3nh/" {
				t.Fatal("muxed function differs from OpenSSL vector")
			}
		})
	}
}
