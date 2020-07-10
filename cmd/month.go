package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func newMonth() *cobra.Command {
	return newCmd("month", "m", "Group files by month", monthFunc)
}

func monthFunc(_ *cobra.Command, _ []string) error {
	return group(appFileSystem, monthGrouper)
}

func monthGrouper(file os.FileInfo) []string {
	year, month, _ := file.ModTime().Date()
	return []string{fmt.Sprintf("%d-%02d", year, month)}
}
