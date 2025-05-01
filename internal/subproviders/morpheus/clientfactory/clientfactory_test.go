package clientfactory_test

import (
	"context"
	"crypto/x509"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/clientfactory"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/model"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSecureTLS(t *testing.T) {
	server := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Simulate a simple 200 OK response
			w.WriteHeader(http.StatusOK)
		}))
	defer server.Close()

	m := model.SubModel{
		URL:         types.StringValue(server.URL),
		Username:    types.StringValue("user"),
		Password:    types.StringValue("secret"),
		AccessToken: types.StringValue("token"),
		Insecure:    types.BoolValue(false),
	}
	cf := clientfactory.New(m)
	c, err := cf.NewClient(context.Background())
	if err != nil {
		t.Fatal("Failed to create client", err)
	}
	u := c.UsersAPI.GetUser(context.Background(), 1)
	_, _, err = u.Execute()
	if err == nil {
		t.Fatal("Failed to raise error", err)
	}
	var certErr x509.UnknownAuthorityError
	if !errors.As(err, &certErr) {
		t.Fatalf("Expected UnknownAuthorityError, got: %v", err)
	}
}

func TestInsecureTLS(t *testing.T) {
	server := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Simulate a simple 200 OK response
			w.WriteHeader(http.StatusOK)
		}))
	defer server.Close()

	m := model.SubModel{
		URL:         types.StringValue(server.URL),
		Username:    types.StringValue("user"),
		Password:    types.StringValue("secret"),
		AccessToken: types.StringValue("token"),
		Insecure:    types.BoolValue(true),
	}
	cf := clientfactory.New(m)
	c, err := cf.NewClient(context.Background())
	if err != nil {
		t.Fatal("Failed to create client", err)
	}
	u := c.UsersAPI.GetUser(context.Background(), 1)
	_, _, err = u.Execute()
	if err != nil {
		t.Fatal("Unexpected error", err)
	}
}
