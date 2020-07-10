package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func newDay(c conf) *cobra.Command {
	return newSimpleGroupingCmd(c, "day", "d", "Group files by day", dayGrouper)
}

func newMonth(c conf) *cobra.Command {
	return newSimpleGroupingCmd(c, "month", "m", "Group files by month", monthGrouper)
}

func newYear(c conf) *cobra.Command {
	return newSimpleGroupingCmd(c, "year", "y", "Group files by year", yearGrouper)
}

func dayGrouper(file os.FileInfo) []string {
	year, month, day := file.ModTime().Date()
	return []string{fmt.Sprintf("%d-%02d-%02d", year, month, day)}
}

func monthGrouper(file os.FileInfo) []string {
	year, month, _ := file.ModTime().Date()
	return []string{fmt.Sprintf("%d-%02d", year, month)}
}

func yearGrouper(file os.FileInfo) []string {
	year, _, _ := file.ModTime().Date()
	return []string{fmt.Sprintf("%d", year)}
}
