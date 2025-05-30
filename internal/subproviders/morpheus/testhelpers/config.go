package testhelpers

import (
	"os"
	"testing"
)

func BuildProviderBlock(t *testing.T) string {
	var providerBlock string
	// prefer access token if both creds AND token are set
	if os.Getenv("TF_ACC_MORPHEUS_ACCESS_TOKEN") != "" {
		providerBlock = RenderExample(t, "../../provider-access-token.tf.tmpl",
			"AccessToken", os.Getenv("TF_ACC_MORPHEUS_ACCESS_TOKEN"),
			"Url", os.Getenv("TF_ACC_MORPHEUS_URL"),
			"Insecure", os.Getenv("TF_ACC_MORPHEUS_INSECURE"),
		)
	} else {
		providerBlock = RenderExample(t, "../../provider-credentials.tf.tmpl",
			"Username", os.Getenv("TF_ACC_MORPHEUS_USERNAME"),
			"Password", os.Getenv("TF_ACC_MORPHEUS_PASSWORD"),
			"Url", os.Getenv("TF_ACC_MORPHEUS_URL"),
			"Insecure", os.Getenv("TF_ACC_MORPHEUS_INSECURE"),
		)
	}

	return providerBlock
}
