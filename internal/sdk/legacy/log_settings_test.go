package morpheus_test

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func TestLogSettings(t *testing.T) {
	client := getTestClient(t)
	req := &morpheus.Request{}
	resp, err := client.GetLogSettings(req)
	assertResponse(t, resp, err)
}
