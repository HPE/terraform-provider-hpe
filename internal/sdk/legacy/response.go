package morpheus

import (
	"fmt"
	"net/http"
	//"errors"
	"time"
	//"encoding/json"
	// "github.com/gormorpheus/morpheusapi/client"
)

type Response struct {
	HTTPResponse *http.Response
	Success      bool
	StatusCode   int
	Status       string
	Body         []byte
	Error        error
	ReceivedAt   time.Time
	Size         int64

	// This holds the parsed JSON for convenience
	JsonData interface{}

	// This holds any error encountering JsonData
	JsonParseError error

	Result interface{}

	// the request is stored in here for referencing the http request.
	request *Request

	// Errors []error
	// Took int // ms
}

func (resp *Response) String() string {
	return fmt.Sprintf("Response HTTP: %v Success: %v  Size: %dB Body: %s",
		resp.Status, resp.Success, resp.Size, resp.Body)
}

func (resp *Response) SetRequest(req *Request) *Response {
	resp.request = req
	return resp
}

func (resp *Response) GetRequest() (req *Request) {
	// if resp.request == nil {
	// 	return nil
	// }
	return resp.request
}

// todo: use something like this so we don't have to typecast everywhere.
// API response interface
type APIResponse interface {

	// whether the request was successful (200 OK) or not.
	Success() bool

	// HTTP status code eg. 200
	StatusCode() int

	// HTTP status message .eg "OK"
	Status() string

	// response body byte array
	Body() []byte

	// an error ocured
	Error() error

	// timestamp of the request
	ReceivedAt() time.Time

	// number of bytes in the response
	Size() int64

	// This holds the parsed JSON for convenience
	JsonData() interface{}

	// This holds any error encountering JsonData
	JsonParseError() error

	// the parsed result, in the specified type.
	Result() interface{}

	// the request is stored in here for referencing the http request.
	Request() *Request

	// errors that happened during the request
	Errors() []error
}
