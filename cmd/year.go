package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// yearCmd represents the year command
var yearCmd = &cobra.Command{
	Use:     "year",
	Aliases: []string{"y"},
	Short:   "Group files by year",
	Run: func(cmd *cobra.Command, args []string) {
		fg(appFileSystem, yearGrouper)
	},
}

func init() {
	rootCmd.AddCommand(yearCmd)
}

func yearGrouper(file os.FileInfo) string {
	year, _, _ := file.ModTime().Date()
	return fmt.Sprintf("%d", year)
}
