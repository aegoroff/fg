package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func newDay() *cobra.Command {
	return &cobra.Command{
		Use:     "day",
		Aliases: []string{"d"},
		Short:   "Group files by day",
		RunE: func(cmd *cobra.Command, args []string) error {
			return group(appFileSystem, dayGrouper)
		},
	}
}

func dayGrouper(file os.FileInfo) []string {
	year, month, day := file.ModTime().Date()
	return []string{fmt.Sprintf("%d-%02d-%02d", year, month, day)}
}
