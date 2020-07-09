package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func newMonth() *cobra.Command {
	return &cobra.Command{
		Use:     "month",
		Aliases: []string{"m"},
		Short:   "Group files by month",
		RunE: func(cmd *cobra.Command, args []string) error {
			return group(appFileSystem, monthGrouper)
		},
	}
}

func monthGrouper(file os.FileInfo) []string {
	year, month, _ := file.ModTime().Date()
	return []string{fmt.Sprintf("%d-%02d", year, month)}
}
