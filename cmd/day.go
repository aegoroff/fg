package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func newDay() *cobra.Command {
	return newSimpleGroupingCmd("day", "d", "Group files by day", dayGrouper)
}

func dayGrouper(file os.FileInfo) []string {
	year, month, day := file.ModTime().Date()
	return []string{fmt.Sprintf("%d-%02d-%02d", year, month, day)}
}
