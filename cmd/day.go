package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func newDay() *cobra.Command {
	return newCmd("day", "d", "Group files by day", dayFunc)
}

func dayFunc(_ *cobra.Command, _ []string) error {
	return group(appFileSystem, dayGrouper)
}

func dayGrouper(file os.FileInfo) []string {
	year, month, day := file.ModTime().Date()
	return []string{fmt.Sprintf("%d-%02d-%02d", year, month, day)}
}
