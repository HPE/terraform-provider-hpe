package customtypes_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/HPE/terraform-provider-hpe/utils/customtypes"
)

func TestNewTrimmedStringValue(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  customtypes.TrimmedString
	}{
		{
			name:  "NormalString",
			value: "test",
			want: customtypes.TrimmedString{
				StringValue: basetypes.NewStringValue("test"),
			},
		},
		{
			name:  "MultilineString",
			value: "\nnewline\ntest\n",
			want: customtypes.TrimmedString{
				StringValue: basetypes.NewStringValue("\nnewline\ntest\n"),
			},
		},
		{
			name:  "EmptyString",
			value: "",
			want: customtypes.TrimmedString{
				StringValue: basetypes.NewStringValue(""),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := customtypes.NewTrimmedStringValue(tt.value)
			if !got.Equal(tt.want) {
				t.Errorf("NewTrimmedStringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrimmedString_StringSemanticEquals(t *testing.T) {
	tests := []struct {
		name        string
		value       customtypes.TrimmedString
		newValuable basetypes.StringValuable
		want        bool
		wantErr     bool
	}{
		{
			name:        "Newlines",
			value:       customtypes.NewTrimmedStringValue("\n\n\nnewline\ntest\n\n\n"),
			newValuable: customtypes.NewTrimmedStringValue("newline\ntest"),
			want:        true,
			wantErr:     false,
		},
		{
			name:        "ReturnNewline",
			value:       customtypes.NewTrimmedStringValue("\r\nnewline\r\ntest\r\n\r\n"),
			newValuable: customtypes.NewTrimmedStringValue("newline\r\ntest"),
			want:        true,
			wantErr:     false,
		},
		{
			name:        "Spaces",
			value:       customtypes.NewTrimmedStringValue("  newline\r\ntest  "),
			newValuable: customtypes.NewTrimmedStringValue("newline\r\ntest"),
			want:        true,
			wantErr:     false,
		},
		{
			name:        "NoWhitespace",
			value:       customtypes.NewTrimmedStringValue("newline\ntest"),
			newValuable: customtypes.NewTrimmedStringValue("newline\ntest"),
			want:        true,
			wantErr:     false,
		},
		{
			name:        "StringValue",
			value:       customtypes.NewTrimmedStringValue("newline\ntest"),
			newValuable: basetypes.NewStringValue("test"),
			want:        false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.value.StringSemanticEquals(context.Background(), tt.newValuable)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("StringSemanticEquals() failed: %v", gotErr)
				}

				return
			}

			if tt.wantErr {
				t.Fatal("StringSemanticEquals() succeeded unexpectedly")
			}

			if got != tt.want {
				t.Errorf("StringSemanticEquals() = %v, want %v", got, tt.want)
			}
		})
	}
}
