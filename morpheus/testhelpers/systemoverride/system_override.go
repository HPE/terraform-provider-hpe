package systemoverride

import (
	"flag"
	"fmt"
	"strings"
	"testing"
)

var systemOverrides map[string]string

type testOverrideFlags []string

func (f *testOverrideFlags) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *testOverrideFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func ParseFlags() {
	if systemOverrides != nil {
		return
	}

	systemOverrides = make(map[string]string)

	var overrideSystem testOverrideFlags
	flag.Var(&overrideSystem, "morpheus.system-override", "override test system")
	flag.Parse()

	for _, override := range overrideSystem {
		m := strings.Split(override, "=")
		systemOverrides[m[0]] = m[1]
	}
}

func GetPreferred(t *testing.T, defaultSystem string) string {
	if system, ok := systemOverrides["all"]; ok {
		return system
	}

	if system, ok := systemOverrides[t.Name()]; ok {
		return system
	}

	return defaultSystem
}

type SystemTestParameters struct {
	Name   string
	Params parameters
}

type parameters map[string]string

func (p parameters) ToSlice() []string {
	var o []string
	for k, v := range p {
		o = append(o, k, v)
	}

	return o
}

func GetParameters(system string, testParameters ...SystemTestParameters) parameters {
	for _, systemParams := range testParameters {
		if strings.EqualFold(system, systemParams.Name) {
			return systemParams.Params
		}
	}

	panic("no parameters found for specified system")
}
