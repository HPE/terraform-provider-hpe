// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package trust

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_key_pair/resource.tf morpheus_key_pair_resource.tf.tmpl Name example-key-pair PublicKey "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC..." PrivateKey "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...\n-----END RSA PRIVATE KEY-----"

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
-----END OPENSSH PRIVATE KEY-----`
)

func RenderKeyPairConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       "example-key-pair",
		"PublicKey":  pubKey,
		"PrivateKey": privKey,
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_key_pair_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
