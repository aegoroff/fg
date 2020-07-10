package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func newMonth() *cobra.Command {
	return newSimpleGroupingCmd("month", "m", "Group files by month", monthGrouper)
}

func monthGrouper(file os.FileInfo) []string {
	year, month, _ := file.ModTime().Date()
	return []string{fmt.Sprintf("%d-%02d", year, month)}
}
