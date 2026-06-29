package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestGetMonitoringSettings(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.GetMonitoringSettings(req)
	assertResponse(t, resp, err)
}
