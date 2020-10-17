package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version defines program version
var Version = "0.2.0"

func newVersion() *cobra.Command {
	cmd := newCmd("version", "ver", "Print the version number of fgr", versionFunc)
	cmd.Long = `All software has versions. This is fgr's`

	return cmd
}

func versionFunc(_ *cobra.Command, _ []string) error {
	_, err := fmt.Fprintf(appWriter, "fgr v%s\n", Version)
	return err
}
