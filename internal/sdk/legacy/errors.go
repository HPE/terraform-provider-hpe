package morpheus

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func errWithBody(err error, resp *http.Response) error {
	var msg string

	if err != nil {
		msg = err.Error()
		return errors.New(msg)
	}

	if resp != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		msg = fmt.Sprintf("%d (%s): %s", resp.StatusCode, http.StatusText(resp.StatusCode), string(bodyBytes))
		return errors.New(msg)
	}

	return errors.New("Unknown error")
}
