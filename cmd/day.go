package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// dayCmd represents the day command
var dayCmd = &cobra.Command{
	Use:     "day",
	Aliases: []string{"d"},
	Short:   "Group files by day",
	Run: func(cmd *cobra.Command, args []string) {
		fg(appFileSystem, dayGrouper)
	},
}

func init() {
	rootCmd.AddCommand(dayCmd)
}

func dayGrouper(file os.FileInfo) []string {
	year, month, day := file.ModTime().Date()
	return []string{fmt.Sprintf("%d-%02d-%02d", year, month, day)}
}
