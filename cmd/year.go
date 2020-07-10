package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func newYear() *cobra.Command {
	return newSimpleGroupingCmd("year", "y", "Group files by year", yearGrouper)
}

func yearGrouper(file os.FileInfo) []string {
	year, _, _ := file.ModTime().Date()
	return []string{fmt.Sprintf("%d", year)}
}
