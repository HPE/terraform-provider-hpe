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

// TestEnvironment is a simplified environment struct for test usage.
type TestEnvironment struct {
	ID   int64
	Name string
	Code string
}

func CreateEnvironment(t *testing.T) (*TestEnvironment, error) {
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

	return &TestEnvironment{
		ID:   environment.GetId(),
		Name: environment.GetName(),
		Code: environment.GetCode(),
	}, nil
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

// GetID returns the ID of the TestEnvironment.
func (e *TestEnvironment) GetID() int64 {
	if e == nil {
		var ret int64

		return ret
	}

	return e.ID
}

// GetName returns the Name of the TestEnvironment.
func (e *TestEnvironment) GetName() string {
	if e == nil {
		var ret string

		return ret
	}

	return e.Name
}

// GetCode returns the Code of the TestEnvironment.
func (e *TestEnvironment) GetCode() string {
	if e == nil {
		var ret string

		return ret
	}

	return e.Code
}
