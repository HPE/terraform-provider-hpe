package customtypes_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/HPE/terraform-provider-hpe/utils/customtypes"
)

func TestNewNormalizedLineEndingsStringValue(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  customtypes.NormalizedLineEndingsString
	}{
		{
			name:  "NormalString",
			value: "test",
			want: customtypes.NormalizedLineEndingsString{
				StringValue: basetypes.NewStringValue("test"),
			},
		},
		{
			name:  "CRLFString",
			value: "line1\r\nline2\r\n",
			want: customtypes.NormalizedLineEndingsString{
				StringValue: basetypes.NewStringValue("line1\r\nline2\r\n"),
			},
		},
		{
			name:  "EmptyString",
			value: "",
			want: customtypes.NormalizedLineEndingsString{
				StringValue: basetypes.NewStringValue(""),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := customtypes.NewNormalizedLineEndingsStringValue(tt.value)
			if !got.Equal(tt.want) {
				t.Errorf("NewNormalizedLineEndingsStringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizedLineEndingsString_StringSemanticEquals(t *testing.T) {
	tests := []struct {
		name        string
		value       customtypes.NormalizedLineEndingsString
		newValuable basetypes.StringValuable
		want        bool
		wantErr     bool
	}{
		{
			// The MORPH-12525 case: config CRLF vs API-returned LF.
			name:        "CRLFvsLF",
			value:       customtypes.NewNormalizedLineEndingsStringValue("#!/bin/bash\r\nset -e\r\n"),
			newValuable: customtypes.NewNormalizedLineEndingsStringValue("#!/bin/bash\nset -e\n"),
			want:        true,
			wantErr:     false,
		},
		{
			name:        "CRvsLF",
			value:       customtypes.NewNormalizedLineEndingsStringValue("a\rb\rc"),
			newValuable: customtypes.NewNormalizedLineEndingsStringValue("a\nb\nc"),
			want:        true,
			wantErr:     false,
		},
		{
			name:        "MixedEndings",
			value:       customtypes.NewNormalizedLineEndingsStringValue("a\r\nb\rc\n"),
			newValuable: customtypes.NewNormalizedLineEndingsStringValue("a\nb\nc\n"),
			want:        true,
			wantErr:     false,
		},
		{
			name:        "IdenticalLF",
			value:       customtypes.NewNormalizedLineEndingsStringValue("a\nb"),
			newValuable: customtypes.NewNormalizedLineEndingsStringValue("a\nb"),
			want:        true,
			wantErr:     false,
		},
		{
			// Line endings are folded, but real content differences still matter.
			name:        "DifferentContent",
			value:       customtypes.NewNormalizedLineEndingsStringValue("a\r\nb"),
			newValuable: customtypes.NewNormalizedLineEndingsStringValue("a\nc"),
			want:        false,
			wantErr:     false,
		},
		{
			// Trailing-newline difference is NOT a line-ending style difference.
			name:        "DifferentTrailingNewline",
			value:       customtypes.NewNormalizedLineEndingsStringValue("a\r\nb\r\n"),
			newValuable: customtypes.NewNormalizedLineEndingsStringValue("a\nb"),
			want:        false,
			wantErr:     false,
		},
		{
			name:        "WrongType",
			value:       customtypes.NewNormalizedLineEndingsStringValue("a\nb"),
			newValuable: basetypes.NewStringValue("a\nb"),
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
