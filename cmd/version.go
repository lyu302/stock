package cmd

import (
	"fmt"
	"os"
)

var version string

func DisplayVersion(module string) {
	if version != "" {
		fmt.Sprintf("Stock %s %s\n", module, version)
	} else {
		fmt.Sprintf("Stock %s %s\n", module, os.Getenv("RELEASE_VERSION"))
	}

	os.Exit(0)
}

func GetVersion() string {
	return version
}
