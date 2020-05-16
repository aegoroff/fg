package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// monthCmd represents the month command
var monthCmd = &cobra.Command{
	Use:     "month",
	Aliases: []string{"m"},
	Short:   "Group files by month",
	Run: func(cmd *cobra.Command, args []string) {
		fg(appFileSystem, monthGrouper)
	},
}

func init() {
	rootCmd.AddCommand(monthCmd)
}

func monthGrouper(file os.FileInfo) string {
	year, month, _ := file.ModTime().Date()
	return fmt.Sprintf("%d-%02d", year, month)
}