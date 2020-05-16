package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// extCmd represents the ext command
var extCmd = &cobra.Command{
	Use:     "ext",
	Aliases: []string{"e"},
	Short:   "Group files by file extension",
	RunE: func(cmd *cobra.Command, args []string) error {
		return group(appFileSystem, extGrouper)
	},
}

func init() {
	rootCmd.AddCommand(extCmd)
}

func extGrouper(file os.FileInfo) []string {
	name := file.Name()
	parts := strings.Split(name, ".")
	if len(parts) < 2 {
		return []string{"no extension"}
	}
	return []string{parts[len(parts)-1]}
}
