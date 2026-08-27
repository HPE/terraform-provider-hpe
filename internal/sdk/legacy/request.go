package morpheus

import (
	"fmt"
	// "github.com/go-resty/resty/v2"
	// "github.com/gormorpheus/morpheusapi/client"
	// "github.com/gormorpheus/morpheusapi/response"
)

type Request struct {
	Method      string
	Path        string
	QueryParams map[string]string
	Headers     map[string]string
	Body        map[string]interface{}
	// FormData interface{}
	// FormData map[string]interface{}
	FormData map[string]string

	// Client Client
	SkipLogin         bool // used for anonymous api calls, otherwise Login() is always called to get token
	SkipAuthorization bool // do not automatically add header for Authorization: Bearer AccessToken
	Timeout           int  // todo:  dictate request timeout

	Result interface{}

	IsMultiPart bool
	isFormData  bool
	IsStream    bool
	StreamBody  string
	// setContentLength    bool
	// isSaveResponse      bool
	// notParseResponse    bool
	// jsonEscapeHTML      bool
	// trace               bool
	// outputFile          string
	// fallbackContentType string
	// ctx                 context.Context
	// pathParams          map[string]string
	// values              map[string]interface{}
	// client              *Client
	// bodyBuf             *bytes.Buffer
	// clientTrace         *clientTrace
	MultiPartFiles []*FilePayload
	// multipartFields     []*MultipartField
}

type FilePayload struct {
	ParameterName string
	FileName      string
	FileContent   []byte
}

func (req *Request) String() string {
	return fmt.Sprintf("Request Method: %v Path: %v QueryParams: %v Body: %s",
		req.Method, req.Path, req.QueryParams, req.Body)
}
