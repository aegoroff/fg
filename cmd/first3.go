package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"strings"
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
	name := file.Name()
	parts := strings.Split(name, ".")
	if len(parts) < 2 {
		return "no extension"
	}
	return parts[len(parts)-1]
}
