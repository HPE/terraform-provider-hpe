// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package notify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	version "github.com/hashicorp/go-version"
)

const (
	RegistryUrl = "https://registry.terraform.io/v1/providers?provider=hpe"

	envHPEIgnoreVersionCheck = "HPE_IGNORE_VERSION_CHECK"

	outdatedProviderErrFmt = `
	A new version of the HPE Morpheus Provider is available (installed version: %s).

	Please upgrade to the latest version (%s) by following the instructions at:

	https://registry.terraform.io/providers/HPE/hpe/latest

	Or by updating the version constraint in your provider configuration block to include
	the latest version and run 'terraform init'.

	To suppress this message, upgrade to the latest version or set the environment variable
	HPE_IGNORE_VERSION_CHECK.
`
)

type registryResponse200 struct {
	Meta      registryMeta       `json:"meta"`
	Providers []registryProvider `json:"providers"`
}

type registryMeta struct {
	Limit         int `json:"limit"`
	CurrentOffset int `json:"current_offset"`
}

type registryProvider struct {
	Id          string `json:"id"`
	Owner       string `json:"owner"`
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Version     string `json:"version"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
	Source      string `json:"source"`
	PublishedAt string `json:"published_at"`
	Downloads   int    `json:"downloads"`
	Tier        string `json:"tier"`
	LogoUrl     string `json:"logo_url"`
}

// Notifier is enabled by default, unless we set IGNORE_VERSION_CHECK
func IsEnabled() bool {
	_, isSet := os.LookupEnv(envHPEIgnoreVersionCheck)
	return !isSet
}

// Dials CloudFlare DNS server to quickly test internet connectivity.
func TryDial() error {
	// We can't use IPV4 ICMP without elevated permissions, so we'll do it over TCP.
	if _, err := net.DialTimeout("tcp", "1.1.1.1:443", 250*time.Millisecond); err != nil {
		return err
	}

	return nil
}

func GetProviderVersion(url string) (*version.Version, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Don't bother enabling HTTP tracing for now on this.
	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var registryResp registryResponse200
	err = json.Unmarshal(b, &registryResp)
	if err != nil {
		return nil, err
	}

	if len(registryResp.Providers) != 1 {
		return nil, errors.New("length of providers must be 1")
	}

	provider := registryResp.Providers[0]

	ver, err := version.NewVersion(provider.Version)
	if err != nil {
		return nil, err
	}

	return ver, nil
}

func CompareProviderVersion(local, remote *version.Version) error {
	if local == nil {
		return errors.New("local must not be nil")
	}

	if remote == nil {
		return errors.New("remote must not be nil")
	}

	if local.LessThan(remote) {
		return fmt.Errorf(outdatedProviderErrFmt, local.String(), remote.String())
	}

	return nil
}
