package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func newDay() *cobra.Command {
	return newSimpleGroupingCmd("day", "d", "Group files by day", dayGrouper)
}

func newMonth() *cobra.Command {
	return newSimpleGroupingCmd("month", "m", "Group files by month", monthGrouper)
}

func newYear() *cobra.Command {
	return newSimpleGroupingCmd("year", "y", "Group files by year", yearGrouper)
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
