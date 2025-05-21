// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package errors

import (
	"fmt"
	"io"
	"net/http"
)

func ErrMsg(err error, resp *http.Response) string {
	var msg string

	if err != nil {
		msg = err.Error()
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
