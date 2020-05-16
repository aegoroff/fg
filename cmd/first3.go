package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// first3Cmd represents the l3 command
var first3Cmd = &cobra.Command{
	Use:     "l3",
	Aliases: []string{},
	Short:   "Group files by first 3 letters of a name",
	Run: func(cmd *cobra.Command, args []string) {
		fg(appFileSystem, first3Grouper)
	},
}

func init() {
	rootCmd.AddCommand(first3Cmd)
}

func first3Grouper(file os.FileInfo) string {
	sz := 3
	name := file.Name()
	if len(name) < sz {
		return name
	}
	return name[0:sz]
}
