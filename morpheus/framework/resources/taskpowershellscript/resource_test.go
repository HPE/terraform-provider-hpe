// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskpowershellscript_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/taskpowershellscript"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

var importIgnore = []string{
	"remote_target_password_wo",
	"remote_target_password_wo_version",
	"write_only_revision",
}

// TestAccMorpheusTaskPowershellScriptResourceExampleOk exercises the local-script
// example: create, an empty plan after apply, and import.
func TestAccMorpheusTaskPowershellScriptResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	addr := "hpe_morpheus_task_powershell_script.tfexample_powershell_local"

	cfg, err := taskpowershellscript.RenderTaskPowershellScriptConfig(t, map[string]string{
		"Name": name,
		"Code": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "name", name),
					resource.TestCheckResourceAttr(addr, "code", name),
					resource.TestCheckResourceAttr(addr, "labels.0", "demo"),
					resource.TestCheckResourceAttr(addr, "labels.1", "terraform"),
					resource.TestCheckResourceAttr(addr, "source_type", "local"),
					resource.TestCheckResourceAttr(addr, "elevated_shell", "true"),
					resource.TestCheckResourceAttr(addr, "retryable", "true"),
					resource.TestCheckResourceAttr(addr, "retry_count", "1"),
					resource.TestCheckResourceAttr(addr, "retry_delay_seconds", "10"),
					resource.TestCheckResourceAttr(addr, "allow_custom_config", "true"),
					resource.TestCheckResourceAttrSet(addr, "write_only_revision"),
				),
			},
			{Config: providerConfig + cfg, PlanOnly: true},
			{
				ResourceName:            addr,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importIgnore,
			},
		},
	})
}

// TestAccMorpheusTaskPowershellScriptResourceGitExampleOk exercises the
// repository-script example.
func TestAccMorpheusTaskPowershellScriptResourceGitExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	repoID, ok := os.LookupEnv("TF_VAR_testacc_morpheus_git_repository_id")
	if !ok {
		t.Skip("TF_VAR_testacc_morpheus_git_repository_id not set; no code-repository fixture available for the repository-source test")
	}

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	addr := "hpe_morpheus_task_powershell_script.tfexample_powershell_git"

	cfg, err := taskpowershellscript.RenderTaskPowershellScriptGitConfig(t, map[string]string{
		"Name":         name,
		"Code":         name,
		"RepositoryId": repoID,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "name", name),
					resource.TestCheckResourceAttr(addr, "source_type", "repository"),
					resource.TestCheckResourceAttr(addr, "result_type", "json"),
					resource.TestCheckResourceAttr(addr, "script_path", "example.ps"),
					resource.TestCheckResourceAttr(addr, "version_ref", "master"),
					resource.TestCheckResourceAttr(addr, "repository_id", repoID),
				),
			},
			{Config: providerConfig + cfg, PlanOnly: true},
			{
				ResourceName:      addr,
				ImportState:       true,
				ImportStateVerify: true,
				// repository_id is config-only: the task API never returns the
				// git repository id (file.repository is always null on read),
				// so it cannot be recovered on import.
				ImportStateVerifyIgnore: append(append([]string{}, importIgnore...), "repository_id"),
			},
		},
	})
}

// TestAccMorpheusTaskPowershellScriptResourceUrlExampleOk exercises the url-script
// example.
func TestAccMorpheusTaskPowershellScriptResourceUrlExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	addr := "hpe_morpheus_task_powershell_script.tfexample_powershell_url"

	cfg, err := taskpowershellscript.RenderTaskPowershellScriptUrlConfig(t, map[string]string{
		"Name": name,
		"Code": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "name", name),
					resource.TestCheckResourceAttr(addr, "source_type", "url"),
					resource.TestCheckResourceAttr(addr, "result_type", "json"),
					resource.TestCheckResourceAttr(addr, "script_path", "https://example.com/example.ps"),
				),
			},
			{Config: providerConfig + cfg, PlanOnly: true},
			{
				ResourceName:            addr,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importIgnore,
			},
		},
	})
}
