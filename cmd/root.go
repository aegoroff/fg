package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// Grouping function
type Grouping func(os.FileInfo) []string

var appFileSystem = afero.NewOsFs()
var appWriter io.Writer

const pathParamName = "path"

var basePath string
var include string
var exclude string

func newRoot() *cobra.Command {
	return &cobra.Command{
		Use:   "fg",
		Short: "Grouping files tool",
		Long: ` fg is a small commandline app that allows you to easily group
 all files in the dir specified into several child subdirectories.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}

func init() {
	cobra.MousetrapHelpText = ""
	appWriter = os.Stdout
}

// Execute starts package running
func Execute(args ...string) error {
	rootCmd := newRoot()

	if args != nil && len(args) > 0 {
		rootCmd.SetArgs(args)
	}

	rootCmd.PersistentFlags().StringVarP(&basePath, pathParamName, "p", "", "REQUIRED. Directory path whose files will be grouped by folders.")
	rootCmd.PersistentFlags().StringVarP(&include, "include", "i", "", "Only files whose names match the pattern specified by the option are grouped.")
	rootCmd.PersistentFlags().StringVarP(&exclude, "exclude", "e", "", "Exclude files whose names match pattern specified by the option from grouping.")

	rootCmd.AddCommand(newDay())
	rootCmd.AddCommand(newExt())
	rootCmd.AddCommand(newFirstn())
	rootCmd.AddCommand(newMonth())
	rootCmd.AddCommand(newUngroup())
	rootCmd.AddCommand(newVersion())
	rootCmd.AddCommand(newYear())

	return rootCmd.Execute()
}
