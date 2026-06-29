package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestGetApplianceSettings(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.GetApplianceSettings(req)
	assertResponse(t, resp, err)
}
