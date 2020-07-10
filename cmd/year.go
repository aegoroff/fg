package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func newYear() *cobra.Command {
	return newCmd("year", "y", "Group files by year", yearFunc)
}

func yearFunc(_ *cobra.Command, _ []string) error {
	return group(appFileSystem, yearGrouper)
}

func yearGrouper(file os.FileInfo) []string {
	year, _, _ := file.ModTime().Date()
	return []string{fmt.Sprintf("%d", year)}
}
