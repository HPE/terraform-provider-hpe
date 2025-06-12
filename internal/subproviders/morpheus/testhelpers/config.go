package testhelpers

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	ErrOneOrMoreEnvVarsNotSet = "One or more environment variables were not set: "

	//nolint:gosec
	envVarURL         = "TF_VAR_testacc_morpheus_url"
	envVarUsername    = "TF_VAR_testacc_morpheus_username"
	envVarPassword    = "TF_VAR_testacc_morpheus_password"
	envVarAccessToken = "TF_VAR_testacc_morpheus_access_token"
	envVarInsecure    = "TF_VAR_testacc_morpheus_insecure"

	providerConfigCredentials = `
variable "testacc_morpheus_url" {}
variable "testacc_morpheus_insecure" {
  default = false
}
variable "testacc_morpheus_username" {}
variable "testacc_morpheus_password" {}

provider "hpe" {
  morpheus {
    url          = var.testacc_morpheus_url
    insecure     = var.testacc_morpheus_insecure
    username     = var.testacc_morpheus_username
    password     = var.testacc_morpheus_password
  }
}
`
	providerConfigAccessToken = `
variable "testacc_morpheus_url" {}
variable "testacc_morpheus_insecure" {
  default = false
}
variable "testacc_morpheus_access_token" {}

provider "hpe" {
  morpheus {
    url          = var.testacc_morpheus_url
    insecure     = var.testacc_morpheus_insecure
    access_token = var.testacc_morpheus_access_token
  }
}
`
)

func buildNotSetMessage(setMap map[string]bool) error {
	var n int
	var notSet string
	// need to do it like this so it's in order; we can't range over the map as it'll be unordered
	if isSet, ok := setMap[envVarURL]; ok && !isSet {
		n++
		notSet += fmt.Sprintf("%s, ", envVarURL)
	}
	if isSet, ok := setMap[envVarUsername]; ok && !isSet {
		n++
		notSet += fmt.Sprintf("%s, ", envVarUsername)
	}
	if isSet, ok := setMap[envVarPassword]; ok && !isSet {
		n++
		notSet += fmt.Sprintf("%s, ", envVarPassword)
	}
	if isSet, ok := setMap[envVarAccessToken]; ok && !isSet {
		n++
		notSet += fmt.Sprintf("%s, ", envVarAccessToken)
	}

	if n > 0 {
		return errors.New(strings.TrimRight(ErrOneOrMoreEnvVarsNotSet+notSet, ", "))
	}

	return nil
}

func BuildProviderBlock() (string, error) {

	setVars := make(map[string]bool)

	// true if set, false if not
	setVars[envVarURL] = os.Getenv(envVarURL) != ""
	setVars[envVarUsername] = os.Getenv(envVarUsername) != ""
	setVars[envVarPassword] = os.Getenv(envVarPassword) != ""
	setVars[envVarAccessToken] = os.Getenv(envVarAccessToken) != ""
	// we don't need to check if insecure is set - it's completely optional

	// prefer access token if both creds AND token are set
	if setVars[envVarAccessToken] {
		delete(setVars, envVarUsername)
		delete(setVars, envVarPassword)

		err := buildNotSetMessage(setVars)
		if err != nil {
			return "", err
		}

		return providerConfigAccessToken, nil
	} else if setVars[envVarUsername] {
		delete(setVars, envVarAccessToken)

		err := buildNotSetMessage(setVars)
		if err != nil {
			return "", err
		}

		return providerConfigCredentials, nil

	}

	return "", buildNotSetMessage(setVars)
}
