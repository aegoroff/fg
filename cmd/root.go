package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var appFileSystem afero.Fs
var appWriter io.Writer

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
	appFileSystem = afero.NewOsFs()
}

// Execute starts package running
func Execute(args ...string) error {
	rootCmd := newRoot()

	if args != nil && len(args) > 0 {
		rootCmd.SetArgs(args)
	}

	conf := newFgConf(appFileSystem)

	rootCmd.PersistentFlags().StringVarP(&conf.bp, "path", "p", "", "REQUIRED. Directory path whose files will be grouped by folders.")
	rootCmd.PersistentFlags().StringVarP(&conf.incl, "include", "i", "", "Only files whose names match the pattern specified by the option are grouped.")
	rootCmd.PersistentFlags().StringVarP(&conf.excl, "exclude", "e", "", "Exclude files whose names match pattern specified by the option from grouping.")

	rootCmd.AddCommand(newDay(conf))
	rootCmd.AddCommand(newMonth(conf))
	rootCmd.AddCommand(newYear(conf))
	rootCmd.AddCommand(newExt(conf))
	rootCmd.AddCommand(newFirstn(conf))
	rootCmd.AddCommand(newUngroup(conf))
	rootCmd.AddCommand(newVersion())

	return rootCmd.Execute()
}
