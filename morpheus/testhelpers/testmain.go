package testhelpers

import (
	"os"
	"slices"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m interface {
	Run() int
},
) int {
	if slices.Contains(os.Args, "-sweep") {
		resource.TestMain(m)

		return 0
	}

	return m.Run()
}
