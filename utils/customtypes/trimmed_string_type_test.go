package customtypes_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/utils/customtypes"
)

func TestTrimmedStringType_ValueFromTerraform(t *testing.T) {
	tests := []struct {
		name    string
		in      tftypes.Value
		want    attr.Value
		wantErr bool
	}{
		{
			name:    "String",
			in:      tftypes.NewValue(tftypes.String, "test string"),
			want:    customtypes.NewTrimmedStringValue("test string"),
			wantErr: false,
		},
		{
			name:    "StringWithNewlines",
			in:      tftypes.NewValue(tftypes.String, "\r\n\r\nnewline\r\ntest\r\n\r\n"),
			want:    customtypes.NewTrimmedStringValue("\r\n\r\nnewline\r\ntest\r\n\r\n"),
			wantErr: false,
		},
		{
			name:    "IncorrectTypeMap",
			in:      tftypes.NewValue(tftypes.Map{}, map[string]tftypes.Value{}),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "IncorrectTypeNumber",
			in:      tftypes.NewValue(tftypes.Number, 1),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tr customtypes.TrimmedStringType
			got, gotErr := tr.ValueFromTerraform(context.Background(), tt.in)
			fmt.Println(gotErr)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValueFromTerraform() failed: %v", gotErr)
				}

				return
			}

			if tt.wantErr {
				t.Fatal("ValueFromTerraform() succeeded unexpectedly")
			}

			if got != tt.want {
				t.Errorf("ValueFromTerraform() = %v(%[1]T), want %v(%[2]T)", got, tt.want)
			}
		})
	}
}

func TestTrimmedStringType_Equal(t *testing.T) {
	tests := []struct {
		name  string
		value customtypes.TrimmedStringType
		in    attr.Type
		want  bool
	}{
		{
			name: "Equal",
			value: customtypes.TrimmedStringType{
				StringType: basetypes.StringType{},
			},
			in: customtypes.TrimmedStringType{
				StringType: basetypes.StringType{},
			},
			want: true,
		},
		{
			name: "NotEqualBasetypeString",
			value: customtypes.TrimmedStringType{
				StringType: basetypes.StringType{},
			},
			in:   basetypes.StringType{},
			want: false,
		},
		{
			name: "NotEqualBasetypeInt64",
			value: customtypes.TrimmedStringType{
				StringType: basetypes.StringType{},
			},
			in:   basetypes.Int64Type{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.value.Equal(tt.in)
			if got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
