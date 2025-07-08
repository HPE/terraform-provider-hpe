// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package errors

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const sslCertErrorMsg = `

If you understand the potential security risks of accepting an untrusted server
certificate, you can bypass this error by setting "insecure = true" in your
provider configuration. Use this option with caution.

provider "hpe" {
   morpheus {
     url = "https://..."
     .
     .
     .
     insecure = true <-- set to true to ignore SSL certificate errors
  }
}
`

func ErrMsg(err error, resp *http.Response) string {
	var msg string

	if err != nil {
		msg = err.Error()
		if strings.Contains(err.Error(), "failed to verify certificate") {
			msg = msg + sslCertErrorMsg
		}
	}

	if resp != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return msg
		}
		code := http.StatusText(resp.StatusCode)
		msg = fmt.Sprintf("%s (%s): %s", msg, code, string(bodyBytes))
	}

	return msg
}
