package errfmt

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Operation represents a CRUD operation for consistent messaging.
type Operation string

const (
	OpCreate Operation = "creating"
	OpRead   Operation = "reading"
	OpUpdate Operation = "updating"
	OpDelete Operation = "deleting"
	OpImport Operation = "importing"
)

// DiagError adds a standardized error diagnostic for a failed API operation.
//
// Produces:
//
//	Summary: "Error {operation} {resourceType}"
//	Detail:  "{resourceType} \"{name}\" failed: {error details}"
func DiagError(
	diags *diag.Diagnostics,
	op Operation,
	resourceType string,
	resourceName string,
	err error,
	httpResp *http.Response,
) {
	summary := fmt.Sprintf("Error %s %s", op, resourceType)
	var detail string
	if resourceName != "" {
		detail = fmt.Sprintf("%s %q failed: %s", resourceType, resourceName, SafeErrMsg(err, httpResp))
	} else {
		detail = fmt.Sprintf("%s failed: %s", resourceType, SafeErrMsg(err, httpResp))
	}
	diags.AddError(summary, detail)
}

// DiagErrorf adds a standardized error diagnostic with a custom detail message.
func DiagErrorf(diags *diag.Diagnostics, op Operation, resourceType string, format string, args ...any) {
	summary := fmt.Sprintf("Error %s %s", op, resourceType)
	detail := fmt.Sprintf(format, args...)
	diags.AddError(summary, detail)
}

// DiagClientError adds the standard "failed to create client" error.
func DiagClientError(diags *diag.Diagnostics, err error) {
	diags.AddError(
		"Error creating API client",
		fmt.Sprintf("Failed to create Morpheus API client: %s", err.Error()),
	)
}

// SafeErrMsg is a nil-safe version of ErrMsg that handles nil error and nil response.
func SafeErrMsg(err error, resp *http.Response) string {
	if err == nil && resp == nil {
		return "unknown error"
	}
	if err == nil {
		return fmt.Sprintf("HTTP %d: %s", resp.StatusCode, safeReadBody(resp))
	}
	if resp == nil {
		msg := err.Error()
		if strings.Contains(msg, "failed to verify certificate") {
			msg = msg + sslCertErrorMsg
		}

		return msg
	}
	// Both present
	msg := err.Error()
	if strings.Contains(msg, "failed to verify certificate") {
		msg = msg + sslCertErrorMsg
	}
	body := safeReadBody(resp)
	code := http.StatusText(resp.StatusCode)
	if body != "" {
		return fmt.Sprintf("%s (%s): %s", msg, code, body)
	}

	return fmt.Sprintf("%s (%s)", msg, code)
}

// IsNotFound returns true if the HTTP response indicates a 404.
// Nil-safe — returns false for nil response.
func IsNotFound(httpResp *http.Response) bool {
	return httpResp != nil && httpResp.StatusCode == http.StatusNotFound
}

// IsSuccess returns true if the HTTP response has a 2xx status code.
// Nil-safe — returns false for nil response.
func IsSuccess(httpResp *http.Response) bool {
	return httpResp != nil && httpResp.StatusCode >= 200 && httpResp.StatusCode < 300
}

// CheckResponse validates an API response and returns a formatted error if it failed.
// Returns nil if err is nil and response is successful (2xx).
func CheckResponse(err error, httpResp *http.Response) error {
	if err != nil {
		return fmt.Errorf("%s", SafeErrMsg(err, httpResp))
	}
	if httpResp != nil && httpResp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", httpResp.StatusCode, safeReadBody(httpResp))
	}

	return nil
}

// safeReadBody reads the response body safely, handling nil and already-consumed bodies.
func safeReadBody(resp *http.Response) string {
	if resp == nil || resp.Body == nil {
		return ""
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil || len(bodyBytes) == 0 {
		return http.StatusText(resp.StatusCode)
	}

	return string(bodyBytes)
}
