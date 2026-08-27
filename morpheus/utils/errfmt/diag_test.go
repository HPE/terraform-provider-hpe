package errfmt

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestSafeErrMsg_BothNil(t *testing.T) {
	msg := SafeErrMsg(nil, nil)
	if msg != "unknown error" {
		t.Fatalf("expected 'unknown error', got %q", msg)
	}
}

func TestSafeErrMsg_NilResponse(t *testing.T) {
	msg := SafeErrMsg(errors.New("connection refused"), nil)
	if msg != "connection refused" {
		t.Fatalf("expected 'connection refused', got %q", msg)
	}
}

func TestSafeErrMsg_NilError(t *testing.T) {
	resp := &http.Response{StatusCode: 404}
	msg := SafeErrMsg(nil, resp)
	if !strings.Contains(msg, "HTTP 404") {
		t.Fatalf("expected 'HTTP 404' in message, got %q", msg)
	}
}

func TestSafeErrMsg_SSLError(t *testing.T) {
	msg := SafeErrMsg(errors.New("failed to verify certificate"), nil)
	if !strings.Contains(msg, "insecure = true") {
		t.Fatalf("expected SSL hint in message, got %q", msg)
	}
}

func TestIsNotFound_True(t *testing.T) {
	resp := &http.Response{StatusCode: 404}
	if !IsNotFound(resp) {
		t.Fatal("expected true for 404")
	}
}

func TestIsNotFound_False(t *testing.T) {
	resp := &http.Response{StatusCode: 200}
	if IsNotFound(resp) {
		t.Fatal("expected false for 200")
	}
}

func TestIsNotFound_Nil(t *testing.T) {
	if IsNotFound(nil) {
		t.Fatal("expected false for nil")
	}
}

func TestIsSuccess_True(t *testing.T) {
	for _, code := range []int{200, 201, 204} {
		resp := &http.Response{StatusCode: code}
		if !IsSuccess(resp) {
			t.Fatalf("expected true for %d", code)
		}
	}
}

func TestIsSuccess_False(t *testing.T) {
	for _, code := range []int{400, 404, 500} {
		resp := &http.Response{StatusCode: code}
		if IsSuccess(resp) {
			t.Fatalf("expected false for %d", code)
		}
	}
}

func TestIsSuccess_Nil(t *testing.T) {
	if IsSuccess(nil) {
		t.Fatal("expected false for nil")
	}
}

func TestCheckResponse_Success(t *testing.T) {
	resp := &http.Response{StatusCode: 200}
	if err := CheckResponse(nil, resp); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCheckResponse_GoError(t *testing.T) {
	err := CheckResponse(errors.New("timeout"), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("expected 'timeout' in error, got %q", err.Error())
	}
}

func TestCheckResponse_HTTPError(t *testing.T) {
	resp := &http.Response{StatusCode: 500}
	err := CheckResponse(nil, resp)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "HTTP 500") {
		t.Fatalf("expected 'HTTP 500' in error, got %q", err.Error())
	}
}

func TestDiagError_Format(t *testing.T) {
	var diags diag.Diagnostics
	DiagError(&diags, OpCreate, "environment", "production", errors.New("conflict"), nil)

	if !diags.HasError() {
		t.Fatal("expected diagnostics to have error")
	}
	err := diags.Errors()[0]
	if err.Summary() != "Error creating environment" {
		t.Fatalf("unexpected summary: %q", err.Summary())
	}
	if !strings.Contains(err.Detail(), `environment "production" failed`) {
		t.Fatalf("unexpected detail: %q", err.Detail())
	}
}

func TestDiagError_NoName(t *testing.T) {
	var diags diag.Diagnostics
	DiagError(&diags, OpRead, "network_pool", "", errors.New("timeout"), nil)

	err := diags.Errors()[0]
	if !strings.Contains(err.Detail(), "network_pool failed") {
		t.Fatalf("unexpected detail: %q", err.Detail())
	}
}

func TestDiagErrorf(t *testing.T) {
	var diags diag.Diagnostics
	DiagErrorf(&diags, OpDelete, "instance", "failed to parse ID %q: %v", "abc", errors.New("not a number"))

	err := diags.Errors()[0]
	if err.Summary() != "Error deleting instance" {
		t.Fatalf("unexpected summary: %q", err.Summary())
	}
	if !strings.Contains(err.Detail(), "failed to parse ID") {
		t.Fatalf("unexpected detail: %q", err.Detail())
	}
}

func TestDiagClientError(t *testing.T) {
	var diags diag.Diagnostics
	DiagClientError(&diags, errors.New("invalid token"))

	err := diags.Errors()[0]
	if err.Summary() != "Error creating API client" {
		t.Fatalf("unexpected summary: %q", err.Summary())
	}
	if !strings.Contains(err.Detail(), "invalid token") {
		t.Fatalf("unexpected detail: %q", err.Detail())
	}
}
