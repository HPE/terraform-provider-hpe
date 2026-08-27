package testhelpers

import (
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m interface {
	Run() int
},
) int {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-sweep") {
			resource.TestMain(m)

			return 0
		}
	}

	return m.Run()
}
