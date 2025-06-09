package testhelpers

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

func CreateEnvironment(t *testing.T) (
	*sdk.ListEnvironments200ResponseAllOfEnvironmentsInner, error,
) {
	t.Helper()

	name := fmt.Sprintf("testacc-%s-%s", t.Name(), rand.Text())

	addEnvironment := sdk.NewAddEnvironmentsRequestEnvironmentWithDefaults()
	addEnvironment.SetName(name)
	addEnvironment.SetCode(strings.ToLower(name))

	addEnvironmentReq := sdk.NewAddEnvironmentsRequest(*addEnvironment)

	ctx := context.TODO()

	client := newClient(ctx, t)

	e, hresp, err := client.EnvironmentsAPI.AddEnvironments(ctx).AddEnvironmentsRequest(
		*addEnvironmentReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("POST failed for Environment %w", err)
	}

	environment := e.GetEnvironment()

	return &environment, nil
}

func DeleteEnvironment(t *testing.T, id int64) error {
	t.Helper()

	ctx := context.TODO()

	client := newClient(ctx, t)

	_, resp, err := client.EnvironmentsAPI.DeleteEnvironments(ctx, id).Execute()
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE failed for Environment %d: %v", id, err)
	}

	for range 6 {
		_, resp, _ := client.EnvironmentsAPI.GetEnvironments(ctx, id).Execute()
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}

		t.Log("Waiting for Environment to be deleted")
		time.Sleep(time.Second * 10)
	}

	return fmt.Errorf("DELETE failed for Environment %d: %v", id, err)
}
