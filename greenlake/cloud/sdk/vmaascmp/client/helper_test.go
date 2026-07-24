// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"testing"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
)

const mockHost = "mockhost"

func getDefaultHeaders() map[string]string {
	return map[string]string{
		"Accept": consts.ContentType,
	}
}

func Test_parseVersion(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    int
		wantErr bool
	}{
		{
			name:    "Normal test case 1: version = 8",
			args:    "8.0.0",
			want:    80000,
			wantErr: false,
		},
		{
			name:    "Normal test case 2: version = 5.3.12",
			args:    "5.3.12",
			want:    50312,
			wantErr: false,
		},
		{
			name:    "Normal test case 2: version = 99.88.77",
			args:    "99.88.77",
			want:    998877,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVersion(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
