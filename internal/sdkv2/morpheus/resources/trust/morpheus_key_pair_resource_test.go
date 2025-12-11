// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package trust_test

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
)

// nolint: lll
const (
	pubKey  = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQCy2vldNrzaqe2vvimpNDQgDHJhCpcUNhkJ0ulpJD+uJMqyn/RpcVzsq3XO+sSvweykJ1bXmUAV/2p/btPnUXsWT2gTaDbXtPYgRpmc8jswpNdl0XEdH3UEpb/ABFU55/LfMY0fvMswbvzzLqK6zxPHmQMqrt+5p0xrzLuqGlaw7w== joe@work`
	privKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAlwAAAAdzc2gtcn
NhAAAAAwEAAQAAAIEAstr5XTa82qntr74pqTQ0IAxyYQqXFDYZCdLpaSQ/riTKsp/0aXFc
7Kt1zvrEr8HspCdW15lAFf9qf27T51F7Fk9oE2g217T2IEaZnPI7MKTXZdFxHR91BKW/wA
RVOefy3zGNH7zLMG788y6ius8Tx5kDKq7fuadMa8y7qhpWsO8AAAIQ/2cKZv9nCmYAAAAH
c3NoLXJzYQAAAIEAstr5XTa82qntr74pqTQ0IAxyYQqXFDYZCdLpaSQ/riTKsp/0aXFc7K
t1zvrEr8HspCdW15lAFf9qf27T51F7Fk9oE2g217T2IEaZnPI7MKTXZdFxHR91BKW/wARV
Oefy3zGNH7zLMG788y6ius8Tx5kDKq7fuadMa8y7qhpWsO8AAAADAQABAAAAgFZCA1ewSX
6Py6Ehbkg7dBQszJD+oYRO3t59CLL7l3auKc/iEuczlCRUQPn0uR0mwrEcg+Zw85ZoW31f
/vSluF1ns6cTcdNj62RV911O35mZofaf2KRjW2C/dpwaS5O/yhKR2uwnI1DX6du3olxueI
hJ4KaD/aZSXGc65XlI/GWhAAAAQB36SILFiFEQ1aTMzD/dgPGK+vBxxz9V7hvJMZWdLfUf
6kRWgx8fn2FgqaoxcXD451MIQdfrrCFnFvUnblZoscgAAABBAOzH+bFyMe+KTNisD5/3zK
lQCDmBxDs7msWhAGJsdMzR/lAyiXd0FaP1FQETbSl62N9sE3mWzFGwNJURJ0yS6lEAAABB
AMFfXwtstG7rrIp7Iv1TKAlVdF7pmlGsoLGonPnAKKllp/+PlwHmdyhY36oOXgZ93bRkDX
3e+bGyAa+ELsrV1z8AAAAUanF1aWdsZXlAbGFwdG9wLXdvcmsBAgMEBQYH
----END OPENSSH PRIVATE KEY-----`
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	return tf5to6server.UpgradeServer(context.Background(), sdkv2morpheus.Provider().GRPCProvider)
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func RenderMorpheusKeyPairConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       acctest.RandomWithPrefix(t.Name()),
		"PublicKey":  pubKey,
		"PrivateKey": privKey,
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"morpheus_key_pair_resource.tf.tmpl",
		"Name", defaults["Name"],
		"PublicKey", defaults["PublicKey"],
		"PrivateKey", defaults["PrivateKey"],
	)
}

func TestAccMorpheusKeyPairExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderMorpheusKeyPairConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_key_pair.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_key_pair.example",
			"public_key",
			pubKey,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Plan
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			// Apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
			},
		},
	})
}
