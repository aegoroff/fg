package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version defines program version
var Version = "0.1.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"ver"},
	Short:   "Print the version number of fg",
	Long:    `All software has versions. This is fg's`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Fprintf(appWriter, "fg v%s\n", Version)
		return err
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}