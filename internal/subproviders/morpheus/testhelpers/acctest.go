package testhelpers

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Helper function to pull out a value from state.
// This can be used to extract a computed value from a resource or data source
func ExtractValue(s *terraform.State, resourceAddress, attributeKey string) (string, error) {
	rs, ok := s.RootModule().Resources[resourceAddress]
	if !ok {
		return "", fmt.Errorf("resource not found: %s", resourceAddress)
	}

	val, ok := rs.Primary.Attributes[attributeKey]
	if !ok {
		return "", fmt.Errorf("attribute not found: %s", val)
	}

	return val, nil
}
