// This provides common functions for all morpheus.Client tests.
package morpheus_test

import (
	"crypto/x509"
	_ "encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	_ "strconv"
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

var (
	// DEBUG = (os.Getenv("MORPHEUS_DEBUG") == "true")
	testUrl         = os.Getenv("MORPHEUS_TEST_URL")
	testUsername    = os.Getenv("MORPHEUS_TEST_USERNAME")
	testPassword    = os.Getenv("MORPHEUS_TEST_PASSWORD")
	testAccessToken = os.Getenv("MORPHEUS_TEST_TOKEN")
	// if os.Getenv("MORPHEUS_TEST_ACCESS_TOKEN") != "" {
	// 	testAccessToken = os.Getenv("MORPHEUS_TEST_ACCESS_TOKEN")
	// }
	testRefreshToken = os.Getenv("MORPHEUS_TEST_REFRESH_TOKEN")
	// testUrl                  string
	// testUsername             string
	// testPassword             string
	sharedTestClient *morpheus.Client
)

// TestMain hooks up setup/teardown methods
// func TestMain(m *testing.M) {
//     setup()
//     code := m.Run()
//     teardown()
//     os.Exit(code)
// }

// func setup() {
// 	//t.Logf(fmt.Sprintf("Setup test suite..."))
// 	log.Printf(fmt.Sprintf("Setup test suite..."))
// }

// func teardown() {
// 	//log.Printf("Teardown test suite...")
// 	if sharedTestClient != nil {
// 		// t.Logf(fmt.Sprintf("Test Client Summary | Success: %d, Error: %d", sharedTestClient.RequestCount(), sharedTestClient.ErrorCount()))
// 		log.Printf(fmt.Sprintf("Test Client Summary | Success: %d, Error: %d", sharedTestClient.RequestCount(), sharedTestClient.ErrorCount())
// 	}
// }

func getNewClient(t *testing.T) *morpheus.Client {
	if testUrl == "" {
		t.Errorf("MORPHEUS_TEST_URL must be set to to run tests")
	}
	t.Logf("Initializing new client for %v", testUrl)
	client := morpheus.NewClient(testUrl)
	return client
}

func TestSecureTLS(t *testing.T) {
	server := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Simulate a simple 200 OK response
			w.WriteHeader(http.StatusOK)
		}))
	defer server.Close()

	client := morpheus.NewClient(server.URL)

	testRequest := &morpheus.Request{
		Method: "GET",
		Path:   "/test",
	}

	_, err := client.Execute(testRequest)

	var certErr x509.UnknownAuthorityError
	if !errors.As(err, &certErr) {
		t.Fatalf("Expected UnknownAuthorityError, got: %v", err)
	}
}

func TestInsecureTLS(t *testing.T) {
	server := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Simulate a simple 200 OK response
			w.WriteHeader(http.StatusOK)
		}))
	defer server.Close()

	Server := morpheus.NewClient(server.URL, morpheus.Insecure())

	testRequest := &morpheus.Request{
		Method: "GET",
		Path:   "/test",
	}

	resp, err := Server.Execute(testRequest)
	assertResponse(t, resp, err)
}

// getTestClient returns a Client that is shared between all tests.
// It is configured via the following environment variables.
// MORPHEUS_TEST_URL - Morpheus Appliance URL
// MORPHEUS_TEST_ACCESS_TOKEN - Morpheus API access token
// MORPHEUS_TEST_USERNAME - Morpheus username
// MORPHEUS_TEST_PASSWORD - Morpheus password
func getTestClient(t *testing.T) *morpheus.Client {
	if sharedTestClient == nil {
		if testUrl == "" {
			panic("MORPHEUS_TEST_URL must be set to to run tests.")
		}
		client := morpheus.NewClient(testUrl)
		if testAccessToken != "" {
			client.SetAccessToken(testAccessToken, testRefreshToken, 0, "write")
			// validate api access token by hitting /api/whoami
			resp, err := client.Whoami()
			// assertResponse(t, resp, err)
			if resp.Success == true {

				whoamiResult := resp.Result.(*morpheus.WhoamiResult)
				currentUsername := whoamiResult.User.Username
				// So this just sets the testUsername
				// it might be better to error if the name does not match...
				if testUsername != currentUsername {
					// need to stop the test right away!
					panic(fmt.Sprintf("MORPHEUS_TEST_USERNAME does not match that of MORPHEUS_TEST_TOKEN. Expected [%v], got [%v]", testUsername, currentUsername))
					// t.Fatalf(fmt.Sprintf("MORPHEUS_TEST_USERNAME does not match that of MORPHEUS_TEST_TOKEN. Expected [%v], got [%v]", testUsername, currentUsername))
				}
				testUsername = currentUsername

			} else {
				assertResponse(t, resp, err)
				// panic("MORPHEUS_TEST_TOKEN could not be validated.")
			}
		} else {
			// authenticate with username and password
			if testUsername == "" || testPassword == "" {
				panic("MORPHEUS_TEST_TOKEN or MORPHEUS_TEST_USERNAME and MORPHEUS_TEST_PASSWORD must be set to to run tests.")
			}
			log.Printf("Initializing test client for %v @ %v", testUsername, testUrl)
			client.SetUsernameAndPassword(testUsername, testPassword)
			resp, err := client.Login()
			assertResponse(t, resp, err)
		}
		sharedTestClient = client
	}
	return sharedTestClient
}

// func parseInt64(s string) (int64, error) {
// 	n, err := strconv.ParseInt(s, 10, 64)
// 	return n, err
// }

func TestPing(t *testing.T) {
	client := getNewClient(t)
	testRequest := &morpheus.Request{
		Method: "GET",
		Path:   "/api/setup/check",
		QueryParams: map[string]string{
			"foo":    "bar",
			"gotest": "true",
		},
	}
	resp, err := client.Execute(testRequest)
	assertResponse(t, resp, err)
}

func TestLogin(t *testing.T) {
	client := getNewClient(t)
	client.SetUsernameAndPassword(testUsername, testPassword)
	resp, err := client.Login()
	assertResponse(t, resp, err)
}

func TestLoginWithToken(t *testing.T) {
	client := getNewClient(t)
	client.SetAccessToken(testAccessToken, testRefreshToken, 0, "write")
	resp, err := client.Whoami()
	assertResponse(t, resp, err)
}

func TestLogout(t *testing.T) {
	client := getNewClient(t)
	client.SetUsernameAndPassword(testUsername, testPassword)
	resp, err := client.Login()
	assertResponse(t, resp, err)
	if err == nil {
		_, err = client.Logout()
		assertError(t, err)
		// assertResponse(t, resp, err)
	}
}

func TestLoginRepeated(t *testing.T) {
	client := getNewClient(t)
	client.SetUsernameAndPassword(testUsername, testPassword)
	resp, err := client.Login()
	assertResponse(t, resp, err)
	if err == nil {
		_, _ = client.Login()
		_, _ = client.Login()
		_, _ = client.Login()
		_, _ = client.Login()
		// assertResponse(t, resp, err)
		assertError(t, err)
	}
}
