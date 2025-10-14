// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package user_test

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
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

	username := acctest.RandomWithPrefix(t.Name())
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

func init() {
	globalRecorders = make(map[string]*recorder.Recorder)
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

	username := acctest.RandomWithPrefix(t.Name())

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

	username := acctest.RandomWithPrefix(t.Name())
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
