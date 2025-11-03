// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package httptrace

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ http.RoundTripper = TraceRoundTripper{}

func IsEnabled() bool {
	_, enabled := os.LookupEnv("MORPHEUS_API_HTTPTRACE")

	return enabled
}

func New(
	transport http.RoundTripper,
) http.RoundTripper {
	return TraceRoundTripper{
		Transport: transport,
	}
}

type TraceRoundTripper struct {
	Transport http.RoundTripper
}

func (t TraceRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBytes, err := httputil.DumpRequestOut(req, true)
	if err == nil {
		tflog.Info(req.Context(),
			"\n\n->TX\n\n"+
				string(reqBytes)+
				"--\n")
	} else {
		msg := fmt.Sprintf("Error tracing request: %v", err)
		tflog.Error(req.Context(), msg)
	}

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	respBytes, err := httputil.DumpResponse(resp, true)
	if err == nil {
		tflog.Info(req.Context(),
			"\n\n<-RX\n\n"+
				string(respBytes)+
				"\n--\n")
	} else {
		msg := fmt.Sprintf("Error tracing response: %v", err)
		tflog.Error(req.Context(), msg)
	}

	return resp, nil
}
