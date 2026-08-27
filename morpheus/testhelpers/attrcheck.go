package testhelpers

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func getResourceAttribute(s *terraform.State, name, key string) (string, error) {
	// these bits are copied from terraform-plugin-testing/helper/resource/testing.go
	ms := s.RootModule()
	rs, ok := ms.Resources[name]
	if !ok {
		return "", fmt.Errorf("Not found: %s in %s", name, ms.Path)
	}

	is := rs.Primary
	if is == nil {
		return "", fmt.Errorf("No primary instance: %s in %s", name, ms.Path)
	}

	v, ok := is.Attributes[key]
	if !ok {
		return "", nil
	}

	return v, nil
}

// CheckResourceAttrEqual checks if two attributes contain equal values
//
// resource[A/B] specifies the name of the resource from which to extract the
// attribute
//
// attr[A/B] specifies the name of the attribute belonging to matching resource
func CheckResourceAttrEqual(resourceA, attrA, resourceB, attrB string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrAValue, err := getResourceAttribute(s, resourceA, attrA)
		if err != nil {
			return err
		}

		attrBValue, err := getResourceAttribute(s, resourceB, attrB)
		if err != nil {
			return err
		}

		if attrAValue == attrBValue {
			return nil
		}

		return fmt.Errorf(
			"attribute '%s.%s' value '%s' does not much attribute '%s.%s' value '%s'",
			resourceA, attrA, attrAValue,
			resourceB, attrB, attrBValue,
		)
	}
}
