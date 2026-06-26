// Client is the driver for interfacing with the Morpheus API
package morpheus

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type clientOptions struct {
	debug           bool
	insecure        bool
	skipLogin       bool
	errCallbackFunc func(err error) error
}

type ClientOption func(*clientOptions)

// WithDebug allows users to enable dumping
// of http requests and responses to stderr.
// Note: This is not recommended for production use - it
// may output sensitive data to logs.
func WithDebug(debug bool) ClientOption {
	return func(options *clientOptions) {
		options.debug = debug
	}
}

// Note: only non-nil callbackFunc return values
// will be used as returned error in Execute()
func WithErrCallbackFunc(callbackFunc func(err error) error) ClientOption {
	return func(options *clientOptions) {
		options.errCallbackFunc = callbackFunc
	}
}

// WithInsecure controls server TLS certificate verification.
// If set to true, it will accept any certificate presented by the server.
// This option should be used with caution.
func WithInsecure(insecure bool) ClientOption {
	return func(options *clientOptions) {
		options.insecure = insecure
	}
}

// The SkipLogin() option is used to explicitly skip the login step performed
// at the start of each call to client.Execute().
// We need to skip login if using e.g. a custom HTTP transport to handle auth.
func SkipLogin() ClientOption {
	return func(options *clientOptions) {
		options.skipLogin = true
	}
}

type Client struct {
	Url             string
	Username        string
	Password        string
	AccessToken     string // todo: make internal
	RefreshToken    string // todo: make internal
	AuthenticatedAt time.Time
	ExpiresIn       int64
	Scope           string
	UserAgent       string
	// Headers map[string]string
	// BaseURL   *url.URL
	HTTPClient *http.Client
	// LastLoginDate time
	// requests []*Request
	lastRequest     *Request
	lastResponse    *Response
	requestCount    int64
	successCount    int64
	errorCount      int64
	debug           bool
	insecure        bool
	skipLogin       bool
	errCallbackFunc func(err error) error
}

// func (client * Client) String() string {
//         return fmt.Sprintf("Client Url: %s Username: %s Logged In: %b", client.Url, client.Username, client.IsLoggedIn())
// }

func (client *Client) IsInsecure() bool {
	return client.insecure
}

func (client *Client) IsSkipLogin() bool {
	return client.skipLogin
}

func (client *Client) IsLoggedIn() bool {
	return client.AccessToken != ""
}

func (client *Client) RequestCount() int64 {
	return client.requestCount
}

func (client *Client) SuccessCount() int64 {
	return client.successCount
}

func (client *Client) ErrorCount() int64 {
	return client.errorCount
}

func (client *Client) incrementRequests(req *Request, resp *Response) {
	client.lastRequest = req
	client.lastResponse = resp
	client.requestCount++
	if resp.Success {
		client.successCount++
	} else {
		client.errorCount++
	}
}

func (client *Client) LastRequest() *Request {
	return client.lastRequest
}

func (client *Client) LastResponse() *Response {
	return client.lastResponse
}

// parseJsonToResult parses json into the given output (struct).
// The type of the ouput determines how it is parsed.
func parseJsonToResult(data []byte, output interface{}) error {
	var err error
	if data != nil {
		err = json.Unmarshal(data, &output)
	}
	return err
}

func NewClient(url string, options ...ClientOption) (client *Client) {
	userAgent := "hpe-morpheus-terraform-provider"

	opts := clientOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	return &Client{
		Url:             url,
		UserAgent:       userAgent,
		debug:           opts.debug,
		insecure:        opts.insecure,
		skipLogin:       opts.skipLogin,
		errCallbackFunc: opts.errCallbackFunc,
	}
}

func (client *Client) SetUsername(username string) *Client {
	// clear access token if switching users
	if client.Username != username {
		client.ClearAccessToken()
		// client.AccessToken = ""
	}
	client.Username = username
	return client
}

func (client *Client) SetPassword(password string) *Client {
	client.Password = password
	return client
}

func (client *Client) SetUsernameAndPassword(username string, password string) *Client {
	client.SetUsername(username)
	client.SetPassword(password)
	return client
}

func (client *Client) SetAccessToken(accessToken string, refreshToken string, expiresIn int64, scope string) *Client {
	client.AccessToken = accessToken
	client.RefreshToken = refreshToken
	client.ExpiresIn = expiresIn
	client.Scope = scope
	return client
}

func (client *Client) ClearAccessToken() *Client {
	client.AccessToken = ""
	client.RefreshToken = ""
	client.ExpiresIn = 0
	client.Scope = ""
	return client
}

func (client *Client) Execute(req *Request) (*Response, error) {
	// first, login if needed
	if !client.skipLogin {
		if !client.IsLoggedIn() && client.Username != "" {
			loginResp, loginErr := client.Login()
			if loginErr != nil {
				return loginResp, loginErr
			}
		}
	}

	// The response object to be returned
	resp := &Response{}

	// potential error to be returned
	var err error

	// construct the request
	httpMethod := req.Method
	if httpMethod == "" {
		// httpMethod = "GET"
		return nil, errors.New("invalid Request: Method is required eg. GET,POST,PUT,DELETE")
	}

	var url string = client.Url + req.Path

	// If we're not using a custom HTTP client, use the default one
	// with TLS settings based on the Morpheus client config
	if client.HTTPClient == nil {
		client.HTTPClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: client.insecure,
				},
			},
		}
	}

	b, err := json.Marshal(req.Body)
	if err != nil {
		return nil, err
	}

	httpReqBody := bytes.NewReader(b)

	httpReq, err := http.NewRequest(req.Method, url, httpReqBody)
	if err != nil {
		return nil, err
	}

	httpReq.Form = make(map[string][]string)

	// set query params
	q := httpReq.URL.Query()
	for k, v := range req.QueryParams {
		q.Add(k, v)
	}
	httpReq.URL.RawQuery = q.Encode()

	// set Headers
	// Set default headers: application/json
	for k, v := range req.Headers {
		httpReq.Header.Add(k, v)
	}

	// add Authorization Header with our access token
	// Not needed with custom transport-level token refresh, but
	// we can keep this in here to maintain similar behaviour to before.
	if !req.SkipAuthorization {
		// NOTE: This will break things if we don't set the token first.
		// if httpReq.Header.Get("Authorization") == "" {
		// 	httpReq.Header.Set("Authorization", "Bearer "+client.AccessToken)
		// }
	}

	// TODO: This could probably be made simpler by using
	// a custom ParseForm() method for the Morpheus request struct.
	// set body
	if httpMethod == "POST" || httpMethod == "PUT" || httpMethod == "PATCH" {
		// FormData means use application/x-www-form-urlencoded

		// ParseForm() will allocate a map for Form, PostForm.
		// It's a no-op otherwise on the httpReq Form data fields right now,
		// since we aren't parsing form data from the body; we need to convert it
		// from a Morpheus request struct.
		// We allocate the memory to avoid nil pointer dereferences in conversions.
		err = httpReq.ParseForm()
		if err != nil {
			return nil, err
		}

		if req.FormData != nil {
			// NOTE: Might need to tidy this up....
			for k, v := range req.FormData {
				if v != "" {
					// WARNING: PANIC HERE unless we explicitly allocate a map
					// for httpReq.Form, which we do in calling httpReq.ParseForm() earlier.
					httpReq.Form.Set(k, v)
				}
			}
			if httpReq.Header.Get("Content-Type") == "" {
				httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
		}

		// NOTE: Need to test that this works correctly with the net/http client
		if req.IsMultiPart {

			var buf bytes.Buffer
			w := multipart.NewWriter(&buf)
			for _, v := range req.MultiPartFiles {
				part, err := w.CreateFormFile(v.ParameterName, v.FileName)
				if err != nil {
					return nil, err
				}
				_, err = part.Write(v.FileContent)
				if err != nil {
					return nil, err
				}
			}

			err = w.Close()
			if err != nil {
				return nil, err
			}

			httpReq.Body = io.NopCloser(&buf)
			httpReq.ContentLength = int64(buf.Len())
			httpReq.Header.Set("Content-Type", w.FormDataContentType())
		}

		// NOTE: Need to test that this works correctly with the net/http client
		if req.IsStream {
			httpReq.Body = io.NopCloser(strings.NewReader(req.StreamBody))
		}

		if req.Body != nil {
			if httpReq.Header.Get("Content-Type") == "" {
				httpReq.Header.Set("Content-Type", "application/json")
			}
		}

		// Set default headers: application/json
		if httpReq.Header.Get("Content-Type") == "" {
			httpReq.Header.Set("Content-Type", "application/json")
		}
	}

	// Set default Accept header
	if httpReq.Header.Get("Accept") == "" {
		httpReq.Header.Set("Accept", "application/json")
	}

	// Dump the request if debug is enabled.
	// This is also what the OpenAPI Generator client does.
	if client.debug {
		dump, err := httputil.DumpRequestOut(httpReq, true)
		if err != nil {
			return resp, err
		}
		log.Printf("\nREQUEST:\n%s\n", string(dump))
	}

	httpResp, err := client.HTTPClient.Do(httpReq)
	if httpResp == nil || err != nil {
		if httpResp == nil {
			return resp, errors.New("nil HTTP response")
		}
		return resp, err
	}

	// Dump the response if debug is enabled.
	// This is also what the OpenAPI Generator client does.
	if client.debug {
		dump, err := httputil.DumpResponse(httpResp, true)
		if err != nil {
			return resp, err
		}
		log.Printf("\nRESPONSE:\n%s\n", string(dump))
	}

	receivedTime := time.Now()

	httpRespBody, err := io.ReadAll(httpResp.Body)
	httpResp.Body.Close()
	// we need the body available for a second read for returning
	// formatted errs in non-success cases.
	httpResp.Body = io.NopCloser(bytes.NewReader(httpRespBody))
	if err != nil {
		return nil, err
	}

	// Convert a net/http response into a Morpheus Response object
	resp = &Response{
		Success:      httpResp.StatusCode > 199 && httpResp.StatusCode < 300, // This is what Resty considers a success.
		StatusCode:   httpResp.StatusCode,
		Status:       httpResp.Status,
		ReceivedAt:   receivedTime,
		Size:         int64(len(httpRespBody)), // This is the same as what Resty does
		Body:         httpRespBody,
		HTTPResponse: httpResp,
	}

	if client.errCallbackFunc != nil {
		customErr := client.errCallbackFunc(err)
		// only use non-nil errors
		if customErr != nil {
			return resp, customErr
		}
	}

	// attempt to parse as json, populates JsonData
	var parsedResult interface{}
	jsonError := parseJsonToResult(resp.Body, &parsedResult)
	resp.JsonData = parsedResult
	resp.JsonParseError = jsonError

	// attempt to parse json into specified result type
	// arbitrary interface{} data is parsed and stored in here
	// The result type is specified in the request right now.
	resp.Result = req.Result
	if resp.Result != nil {
		jsonParseResultError := parseJsonToResult(resp.Body, &resp.Result)
		if jsonParseResultError != nil {
			// maybe actually treat this as a failure..
			log.Printf("Failed to parse JSON result for type %T. Parse Error: %v", resp.Result, jsonParseResultError)
			// log.Errorf("Parse Error: %v", jsonParseResultError)
			// err = jsonParseResultError
			// resp.Success = false
		}
	}

	// print for debugging
	// avoid printing request body for now, it may have secrets.
	// if req.Body != nil {
	// 	log.Printf(fmt.Sprintf("==> Request: %s %s JSON: %s", req.Method, url, req.Body))
	// } else if req.FormData != nil {
	// 	log.Printf(fmt.Sprintf("==> Request: %s %s BODY: %s", req.Method, url, req.FormData))
	// } else {
	// 	log.Printf(fmt.Sprintf("==> Request: %s %s", req.Method, url))
	// }

	// uncomment this for lots of output...
	// log.Printf("API Response: [%v] %d %s", resp.Success, resp.StatusCode, resp.Body)
	// if resp.Success {
	// 	log.Printf("API Response: %d %s", resp.StatusCode, resp.Body)
	// } else {
	// 	log.Printf(fmt.Sprintf("Bad API Response: %d %s", resp.StatusCode, resp.Body))
	// }
	// if err != nil {
	// 	log.Printf("API Error: %v", err)
	// }

	client.incrementRequests(req, resp)

	// so we can still access the Morpheus Resp body in non-success cases
	if !resp.Success {
		return resp, errWithBody(nil, httpResp)
	}

	// success
	return resp, nil
}

func (client *Client) Get(req *Request) (*Response, error) {
	req.Method = "GET"
	return client.Execute(req)
}

func (client *Client) Post(req *Request) (*Response, error) {
	req.Method = "POST"
	return client.Execute(req)
}

func (client *Client) Put(req *Request) (*Response, error) {
	req.Method = "PUT"
	return client.Execute(req)
}

func (client *Client) Delete(req *Request) (*Response, error) {
	req.Method = "DELETE"
	return client.Execute(req)
}

func (client *Client) Patch(req *Request) (*Response, error) {
	req.Method = "PATCH"
	return client.Execute(req)
}

func (client *Client) Head(req *Request) (*Response, error) {
	req.Method = "HEAD"
	return client.Execute(req)
}

func (client *Client) Options(req *Request) (*Response, error) {
	req.Method = "OPTIONS"
	return client.Execute(req)
}

// func (client * Client) List(req * Request) (*Response, error) {
// 	req.Method = "LIST"
// 	return client.Execute(req)
// }

type LoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
}

// NOTE: Login not currently working correctly.
func (client *Client) Login() (*Response, error) {
	// already logged in
	if client.IsLoggedIn() {
		// log.Printf("Login skipped. Already logged in as: %v", client.Username)
		return nil, nil
	} else {
		// c(fmt.Sprintf("Logging in as %s at %s", client.Username, client.Url))
		loginRequest := &Request{
			Method: "POST",
			Path:   "/oauth/token",
			QueryParams: map[string]string{
				"client_id":  "morph-api",
				"grant_type": "password",
				"scope":      "write",
				"username":   client.Username,
			},
			FormData: map[string]string{
				//"username": client.username,
				"password": client.Password,
			},
			Timeout:   10,
			SkipLogin: true,
		}
		resp, err := client.Execute(loginRequest)

		if resp.Success {
			var loginResult LoginResult
			jsonErr := json.Unmarshal(resp.Body, &loginResult)
			if jsonErr != nil {
				// logError(fmt.Sprintf("Error parsing JSON result for type %T [%v]", loginResult, jsonErr))
				return resp, jsonErr
			}
			// log.Printf("LOGIN RESPONSE: ", resp, err)
			// log.Printf("PARSED LOGIN RESULT: ", loginResult)

			if loginResult.AccessToken != "" {
				client.SetAccessToken(loginResult.AccessToken, loginResult.RefreshToken, loginResult.ExpiresIn, loginResult.Scope)
				// log.Printf("Logged in as %v @ %v", client.Username, client.Url)
				// log.Printf("Access Token: ", client.AccessToken)
			} else {
				err = errors.New("Login failed, unable to parse access token from login response")
				// logError(err)
			}
			// client.setLastLoginResult(loginResult)
			return resp, err
		} else {
			log.Printf("Login Failure: %v", resp)
			return resp, err
		}

	}

	// return resp, err
}

func (client *Client) Logout() (*Response, error) {
	client.ClearAccessToken()
	// client.AccessToken = ""
	// client.RefreshToken = ""
	// client.ExpiresIn = 0
	// client.Scope = ""
	client.Username = ""
	client.Password = ""
	// there is no serverside endpoint for this right now
	// create mock response
	resp := &Response{Success: true, Status: "200 OK", StatusCode: 200}
	return resp, nil
}
