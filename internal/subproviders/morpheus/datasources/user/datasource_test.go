// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/clientfactory"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/environment"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/user/consts"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

func TestMain(m *testing.M) {
	code := m.Run()

	// Stop and save all VCR recorders
	for testName, recorder := range globalRecorders {
		if err := recorder.Stop(); err != nil {
			fmt.Printf("❌ Error stopping VCR recorder for %s: %v\n", testName, err)
		} else {
			fmt.Printf("💾 VCR Cassette saved for test: %s\n", testName)
		}
	}

	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusUserDataSourceFindByUsername(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	// Get username for VCR replay - use cassette username in replay mode, generate new in record mode
	username := getTestUsernameForMode(t, "TestAccMorpheusUserDataSourceFindByUsername")

	// Set the username for VCR replay
	setTestUsername("TestAccMorpheusUserDataSourceFindByUsername", username)

	providerConfig := testhelpers.ProviderBlock()

	userResourceConfig := `
resource "hpe_morpheus_user" "test_user" {
	username = "` + username + `"
	role_ids = [1]
	email    = "foo@testacc.com"
	password_wo = "Test123!!"
}
`

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-username.tf.tmpl", "Username", username)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test",
			"username",
			username,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: createTestProviderFactories("TestAccMorpheusUserDataSourceFindByUsername"),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + userResourceConfig,
			},
			{
				Config: providerConfig + userResourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

var globalRecorders map[string]*recorder.Recorder
var testUsernameMap map[string]string // Maps test name to current username

func init() {
	globalRecorders = make(map[string]*recorder.Recorder)
	testUsernameMap = make(map[string]string)
}

// setTestUsername stores the current test's username for use in hooks
func setTestUsername(testName, username string) {
	testUsernameMap[testName] = username
	fmt.Printf("🏷️  Test username set for %s: %s\n", testName, username)
}

// getTestUsername retrieves the current test's username
func getTestUsername(testName string) string {
	if username, exists := testUsernameMap[testName]; exists {
		return username
	}
	return ""
}

// getTestUsernameForMode returns appropriate username based on test mode (record vs replay)
func getTestUsernameForMode(t *testing.T, testName string) string {
	if testing.Short() {
		// In replay-only mode, try to extract username from cassette
		cassettePath := fmt.Sprintf("testdata/%s.yaml", testName)
		fmt.Printf("🔍 Looking for cassette at: %s\n", cassettePath)
		if username := extractUsernameFromCassette(cassettePath, testName); username != "" {
			fmt.Printf("✅ Extracted username from cassette: %s\n", username)
			return username
		}
		fmt.Printf("⚠️  Could not extract username from cassette, using fallback\n")
		// Fallback to deterministic username if cassette doesn't exist or can't be read
		return fmt.Sprintf("%s-replay", testName)
	}
	// In record mode, generate a new random username
	return acctest.RandomWithPrefix(testName)
}

// extractUsernameFromCassette reads the cassette file and extracts the username from the first POST request
func extractUsernameFromCassette(cassettePath, testName string) string {
	// Try to read the cassette file
	if _, err := os.Stat(cassettePath); os.IsNotExist(err) {
		fmt.Printf("❌ Cassette file does not exist: %s\n", cassettePath)
		return ""
	}

	content, err := os.ReadFile(cassettePath)
	if err != nil {
		fmt.Printf("❌ Error reading cassette file: %v\n", err)
		return ""
	}

	// Look for username in the request body using regex that handles JSON escaping
	// The pattern looks for \"username\":\"USERNAME\" where backslashes escape the quotes
	re := regexp.MustCompile(`\\"username\\":\\"([^"\\]+)\\"`)
	if matches := re.FindStringSubmatch(string(content)); len(matches) > 1 {
		username := matches[1]
		fmt.Printf("✅ Found username in cassette: %s\n", username)
		return username
	}
	fmt.Printf("❌ No username found in cassette content\n")
	return ""
}

// customUsernameMatcher is a custom matcher for VCR that normalizes usernames
// in request bodies to handle dynamic test names with random suffixes
func customUsernameMatcher(incomingReq *http.Request, recordedReq cassette.Request) bool {
	fmt.Printf("🔍 Custom matcher called - Method: %s, URL: %s\n", incomingReq.Method, incomingReq.URL.String())

	// For GET requests with username query parameters, normalize the URL before comparing
	if incomingReq.Method == "GET" && strings.Contains(incomingReq.URL.String(), "username=") {
		normalizedIncomingURL := normalizeUsernameInURL(incomingReq.URL.String())
		normalizedRecordedURL := normalizeUsernameInURL(recordedReq.URL)

		if normalizedRecordedURL != normalizedIncomingURL {
			return false
		}

		// Check method match
		if recordedReq.Method != incomingReq.Method {
			return false
		}

		// Check headers match (excluding dynamic ones like Authorization)
		return headersMatch(recordedReq.Headers, incomingReq.Header)
	}

	// For other requests, use the original logic
	// Check if URL and method match first
	if recordedReq.URL != incomingReq.URL.String() || recordedReq.Method != incomingReq.Method {
		fmt.Printf("❌ URL or method mismatch: recorded URL=%s, incoming URL=%s, recorded method=%s, incoming method=%s\n",
			recordedReq.URL, incomingReq.URL.String(), recordedReq.Method, incomingReq.Method)
		return false
	}

	// Check headers match (excluding dynamic ones like Authorization)
	if !headersMatch(recordedReq.Headers, incomingReq.Header) {
		fmt.Printf("❌ Headers mismatch\n")
		return false
	}

	// For POST requests with JSON bodies containing usernames, normalize them
	if incomingReq.Method == "POST" && recordedReq.Body != "" && incomingReq.Body != nil {
		// Read the incoming request body
		var bodyBytes []byte
		if incomingReq.Body != nil {
			var err error
			bodyBytes, err = io.ReadAll(incomingReq.Body)
			if err != nil {
				return false
			}
			// Restore the body for subsequent use
			incomingReq.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		result := normalizedBodiesMatch(recordedReq.Body, string(bodyBytes))
		fmt.Printf("🔍 Body comparison result: %v\n", result)
		if !result {
			fmt.Printf("❌ Body mismatch - recorded: %s\n", recordedReq.Body)
			fmt.Printf("❌ Body mismatch - incoming: %s\n", string(bodyBytes))
		}
		return result
	}

	// For GET requests and other requests without body, standard matching is sufficient
	// since URL, method, and headers already match
	return true
}

// normalizeUsernameInURL replaces dynamic username suffixes in URL query parameters
func normalizeUsernameInURL(url string) string {
	// Find username query parameter and normalize it
	re := regexp.MustCompile(`(\?|&)username=([^&]+)`)
	return re.ReplaceAllStringFunc(url, func(match string) string {
		parts := strings.SplitN(match, "=", 2)
		if len(parts) == 2 {
			username := parts[1]
			// Replace the random suffix with a placeholder
			usernameRe := regexp.MustCompile(`^([^-]+)-\d+$`)
			if matches := usernameRe.FindStringSubmatch(username); len(matches) > 1 {
				return parts[0] + "=" + matches[1] + "-PLACEHOLDER"
			}
		}
		return match
	})
}

// headersMatch compares headers, ignoring dynamic ones like Authorization
func headersMatch(recorded map[string][]string, incoming http.Header) bool {
	// Compare only static headers - Authorization is deliberately excluded
	// since it's removed by the BeforeSaveHook for security
	staticHeaders := []string{"Accept", "Content-Type", "User-Agent"}

	for _, headerName := range staticHeaders {
		recordedValues := recorded[headerName]
		incomingValues := incoming[headerName]

		if len(recordedValues) != len(incomingValues) {
			return false
		}

		for i, val := range recordedValues {
			if i >= len(incomingValues) || val != incomingValues[i] {
				return false
			}
		}
	}

	return true
}

// normalizedBodiesMatch compares JSON request bodies with username normalization
func normalizedBodiesMatch(recordedBody, incomingBody string) bool {
	// Parse both bodies as JSON
	var recordedData, incomingData map[string]interface{}

	if err := json.Unmarshal([]byte(recordedBody), &recordedData); err != nil {
		// If not JSON, fall back to exact match
		return recordedBody == incomingBody
	}

	if err := json.Unmarshal([]byte(incomingBody), &incomingData); err != nil {
		return false
	}

	// Normalize usernames and passwords in both request bodies to handle sanitization
	normalizeDataForMatching(recordedData)
	normalizeDataForMatching(incomingData)

	// Convert back to JSON strings for comparison
	recordedNormalized, err1 := json.Marshal(recordedData)
	incomingNormalized, err2 := json.Marshal(incomingData)

	if err1 != nil || err2 != nil {
		return false
	}

	return string(recordedNormalized) == string(incomingNormalized)
}

// normalizeUsernameInData replaces dynamic username suffixes with a fixed placeholder
func normalizeUsernameInData(data map[string]interface{}) {
	if user, ok := data["user"].(map[string]interface{}); ok {
		// List of username fields that need normalization
		usernameFields := []string{"username", "linuxUsername", "windowsUsername"}

		for _, field := range usernameFields {
			if username, ok := user[field].(string); ok && username != "" {
				// Replace the random suffix with a placeholder
				// Pattern: TestName-RandomNumber -> TestName-PLACEHOLDER
				re := regexp.MustCompile(`^([^-]+)-\d+$`)
				if matches := re.FindStringSubmatch(username); len(matches) > 1 {
					user[field] = matches[1] + "-PLACEHOLDER"
				}
			}
		}
	}
}

// normalizeDataForMatching normalizes both usernames and passwords for VCR matching
func normalizeDataForMatching(data map[string]interface{}) {
	if user, ok := data["user"].(map[string]interface{}); ok {
		// Normalize username fields
		usernameFields := []string{"username", "linuxUsername", "windowsUsername"}
		for _, field := range usernameFields {
			if username, ok := user[field].(string); ok && username != "" {
				// Replace the random suffix with a placeholder
				// Pattern: TestName-RandomNumber -> TestName-PLACEHOLDER
				re := regexp.MustCompile(`^([^-]+)-\d+$`)
				if matches := re.FindStringSubmatch(username); len(matches) > 1 {
					user[field] = matches[1] + "-PLACEHOLDER"
				}
			}
		}

		// Normalize password fields to match sanitized cassettes
		passwordFields := []string{"password", "linuxPassword", "windowsPassword"}
		for _, field := range passwordFields {
			if password, ok := user[field].(string); ok && password != "" {
				// Replace any actual password with the redacted placeholder
				user[field] = "[REDACTED]"
			}
		}
	}
}

// beforeResponseReplayHook replaces recorded usernames with current test usernames
func beforeResponseReplayHook(testName string) func(*cassette.Interaction) error {
	return func(i *cassette.Interaction) error {
		currentUsername := getTestUsername(testName)
		if currentUsername == "" {
			// If no username is set, we can't do replacement
			return nil
		}

		// List of username fields that need replacement in responses
		usernameFields := []string{"username", "linuxUsername", "windowsUsername"}

		for _, field := range usernameFields {
			// Find the recorded username pattern in the response body for this field
			fieldPattern := fmt.Sprintf(`"%s":"([^"]+)"`, field)
			re := regexp.MustCompile(fieldPattern)
			matches := re.FindAllStringSubmatch(i.Response.Body, -1)

			for _, match := range matches {
				if len(match) > 1 {
					recordedUsername := match[1]

					// Check if this looks like a test-generated username with the pattern TestName-RandomNumber
					basePattern := regexp.MustCompile(`^([^-]+)-\d+$`)
					if baseMatches := basePattern.FindStringSubmatch(recordedUsername); len(baseMatches) > 1 {
						recordedBaseName := baseMatches[1]

						// Check if the recorded base name matches our current test name pattern
						if strings.Contains(testName, recordedBaseName) || strings.Contains(currentUsername, recordedBaseName) {
							// Replace with the current test's username
							i.Response.Body = strings.ReplaceAll(i.Response.Body, recordedUsername, currentUsername)

							fmt.Printf("🔄 VCR %s replaced: %s -> %s in response for test %s\n",
								field, recordedUsername, currentUsername, testName)
						}
					}
				}
			}
		}

		// Also handle displayName which might be "FirstName LastName" format
		displayNamePattern := `"displayName":"([^"]+)"`
		re := regexp.MustCompile(displayNamePattern)
		matches := re.FindAllStringSubmatch(i.Response.Body, -1)

		for _, match := range matches {
			if len(match) > 1 {
				recordedDisplayName := match[1]
				// Check if displayName contains the recorded username pattern
				basePattern := regexp.MustCompile(`([^-]+)-\d+`)
				if baseMatches := basePattern.FindStringSubmatch(recordedDisplayName); len(baseMatches) > 1 {
					recordedBaseName := baseMatches[1]

					if strings.Contains(testName, recordedBaseName) || strings.Contains(currentUsername, recordedBaseName) {
						// Replace the username part in displayName
						newDisplayName := basePattern.ReplaceAllString(recordedDisplayName, currentUsername)
						i.Response.Body = strings.ReplaceAll(i.Response.Body, recordedDisplayName, newDisplayName)

						fmt.Printf("🔄 VCR displayName replaced: %s -> %s in response for test %s\n",
							recordedDisplayName, newDisplayName, testName)
					}
				}
			}
		}

		return nil
	}
}

// debugTransport wraps the VCR transport to log all requests
type debugTransport struct {
	Transport http.RoundTripper
}

func (d *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Printf("🌐 HTTP Request: %s %s\n", req.Method, req.URL.String())
	resp, err := d.Transport.RoundTrip(req)
	if err != nil {
		fmt.Printf("❌ HTTP Error: %v\n", err)
		return resp, err
	}
	fmt.Printf("✅ HTTP Response: %d %s\n", resp.StatusCode, resp.Status)
	return resp, err
}

func getOrCreateRecorder(testName string) (*recorder.Recorder, error) {
	if r, exists := globalRecorders[testName]; exists {
		fmt.Printf("🔄 Reusing existing VCR Recorder for test: %s\n", testName)
		return r, nil
	}

	// Determine mode based on testing.Short()
	var mode recorder.Mode
	var modeStr string
	if testing.Short() {
		// Short mode: replay only (fast tests)
		mode = recorder.ModeReplayOnly
		modeStr = "replay-only"
	} else {
		// Full mode: record or replay with new episodes (can record new interactions)
		mode = recorder.ModeReplayWithNewEpisodes
		modeStr = "record-with-new-episodes"
	}

	cassetteName := fmt.Sprintf("testdata/%s", testName)
	r, err := recorder.New(cassetteName,
		recorder.WithMode(mode),
		// Add custom matcher to handle dynamic usernames in request bodies
		recorder.WithMatcher(customUsernameMatcher),
		// Enable replayable interactions to allow same interaction to be replayed multiple times
		recorder.WithReplayableInteractions(true),
		// Add hook to replace usernames in responses during replay
		recorder.WithHook(beforeResponseReplayHook(testName), recorder.BeforeResponseReplayHook),
		// Add hook to remove sensitive headers before saving to cassette
		recorder.WithHook(func(i *cassette.Interaction) error {
			// Remove Authorization header to prevent sensitive tokens from being saved
			if i.Request.Headers != nil {
				delete(i.Request.Headers, "Authorization")
				fmt.Printf("🔒 VCR Authorization header removed for %s %s\n", i.Request.Method, i.Request.URL)
			}

			// Sanitize sensitive data in request bodies
			if i.Request.Body != "" {
				// Remove passwords from request bodies
				passwordRegex := regexp.MustCompile(`"password":"[^"]*"`)
				i.Request.Body = passwordRegex.ReplaceAllString(i.Request.Body, `"password":"[REDACTED]"`)

				linuxPasswordRegex := regexp.MustCompile(`"linuxPassword":"[^"]*"`)
				i.Request.Body = linuxPasswordRegex.ReplaceAllString(i.Request.Body, `"linuxPassword":"[REDACTED]"`)

				windowsPasswordRegex := regexp.MustCompile(`"windowsPassword":"[^"]*"`)
				i.Request.Body = windowsPasswordRegex.ReplaceAllString(i.Request.Body, `"windowsPassword":"[REDACTED]"`)
			}

			// Sanitize sensitive data in response bodies
			if i.Response.Body != "" {
				// Replace account names with generic names
				accountNameRegex := regexp.MustCompile(`"name":"[^"]*QA[^"]*"`)
				i.Response.Body = accountNameRegex.ReplaceAllString(i.Response.Body, `"name":"Test Account"`)
			}

			// Remove sensitive response headers
			if i.Response.Headers != nil {
				delete(i.Response.Headers, "Set-Cookie")
				delete(i.Response.Headers, "XSRF-TOKEN")

				// Sanitize CSP nonces
				if csp, exists := i.Response.Headers["Content-Security-Policy"]; exists {
					for j, cspValue := range csp {
						nonceRegex := regexp.MustCompile(`'nonce-[^']*'`)
						i.Response.Headers["Content-Security-Policy"][j] = nonceRegex.ReplaceAllString(cspValue, "'nonce-REDACTED'")
					}
				}
			}

			return nil
		}, recorder.BeforeSaveHook),
		// Keep the existing after capture hook for logging
		recorder.WithHook(func(i *cassette.Interaction) error {
			fmt.Printf("🎥 VCR AfterCapture [%s]: %s %s -> %d %s\n",
				testName, i.Request.Method, i.Request.URL, i.Response.Code, i.Response.Status)
			return nil
		}, recorder.AfterCaptureHook))
	if err != nil {
		return nil, err
	}

	globalRecorders[testName] = r

	fmt.Printf("🎬 VCR Recorder created for test: %s (cassette: %s, mode: %s, new cassette: %v)\n",
		testName, cassetteName, modeStr, r.IsNewCassette())

	return r, nil
}

func newProviderWithVCR(testName string) (tfprotov6.ProviderServer, error) {
	r, err := getOrCreateRecorder(testName)
	if err != nil {
		return nil, err
	}

	f := func(m model.SubModel) *clientfactory.ClientFactory {
		// example of passing in custom http client
		client := r.GetDefaultClient()

		// Wrap the client transport to log ALL HTTP requests
		originalTransport := client.Transport
		client.Transport = &debugTransport{
			Transport: originalTransport,
		}

		fmt.Printf("🔌 VCR HTTP Client created for test %s: %T (transport: %T)\n", testName, client, client.Transport)

		return clientfactory.New(
			m,
			clientfactory.WithFactoryHTTPClient(client),
		)
	}
	providerInstance := provider.New("test", morpheus.New(morpheus.WithClientFactory(f)))()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

func createTestProviderFactories(testName string) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"hpe": func() (tfprotov6.ProviderServer, error) {
			return newProviderWithVCR(testName)
		},
	}
}

func TestAccMorpheusUserDataSourceFindById(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	username := getTestUsernameForMode(t, "TestAccMorpheusUserDataSourceFindById")

	// Set the username for VCR replay
	setTestUsername("TestAccMorpheusUserDataSourceFindById", username)

	providerConfig := testhelpers.ProviderBlock()

	userResourceConfig := `
resource "hpe_morpheus_user" "test_user" {
	username = "` + username + `"
	role_ids = [1]
	email    = "foo@testacc.com"
	password_wo = "Test123!!"
}
`

	dataSourceConfig := `
    data "hpe_morpheus_user" "test" {
        id = hpe_morpheus_user.test_user.id
    }
    `

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test",
			"username",
			username,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: createTestProviderFactories("TestAccMorpheusUserDataSourceFindById"),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + userResourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusUserDataSourceNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
      data "hpe_morpheus_user" "test" {
        username = "______"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_user.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorNoUserFound

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: createTestProviderFactories("TestAccMorpheusUserDataSourceNotFound"),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusUserDataSourceNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_user" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_user.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorNoValidUserTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: createTestProviderFactories("TestAccMorpheusUserDataSourceNoSearchAttrs"),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusUserDataSourceBothSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
      data "hpe_morpheus_user" "test" {
        id = "1"
        username = "testuser"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_user.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := environment.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: createTestProviderFactories("TestAccMorpheusUserDataSourceBothSearchAttrs"),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

// Test to verify that all of the attributes from a created user can be read
func TestAccMorpheusUserDataSourceVerifyAttributes(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	username := getTestUsernameForMode(t, "TestAccMorpheusUserDataSourceVerifyAttributes")

	// Set the username for VCR replay
	setTestUsername("TestAccMorpheusUserDataSourceVerifyAttributes", username)

	email := "foo@testacc.com"
	firstName := "TestFirst"
	lastName := "TestLast"

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := `
resource "hpe_morpheus_user" "test_all" {
  username     = "` + username + `"
  email        = "` + email + `"
  first_name   = "` + firstName + `"
  last_name    = "` + lastName + `"
  role_ids     = [1]
  password_wo  = "Test123!!"
  receive_notifications = true
  linux_username = "` + username + `"
  windows_username = "` + username + `"
}
`

	dataSourceConfig := `
data "hpe_morpheus_user" "test_all" {
  username = "` + username + `"
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test_all",
			"username",
			username,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test_all",
			"email",
			email,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test_all",
			"first_name",
			firstName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test_all",
			"last_name",
			lastName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test_all",
			"receive_notifications",
			"true",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test_all",
			"linux_username",
			username,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test_all",
			"windows_username",
			username,
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"id",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"roles.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"tenant.id",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"tenant.name",
		),
		/* Can't test this yet as we cannot assign a default persona via terraform yet
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"default_persona.id",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"default_persona.code",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"default_persona.name",
		),
		*/
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.blueprints.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.catalog_item_types.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.features.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.instance_types.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.personas.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.report_types.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.groups.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.workflows.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.tasks.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.vdi_pools.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.clouds.#",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: createTestProviderFactories("TestAccMorpheusUserDataSourceVerifyAttributes"),
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: false,
				Config:             providerConfig + resourceConfig,
			},
			{
				ExpectNonEmptyPlan: false,
				Config:             providerConfig + resourceConfig + dataSourceConfig,
				Check:              checkFn,
			},
		},
	})
}
