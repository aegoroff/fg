package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func newDay(c conf) *cobra.Command {
	cmd := &simple{
		use:   "day",
		a:     "d",
		descr: "Group files by day",
		g:     dayGrouper,
	}
	return newSimpleGroupingCmd(c, cmd)
}

func newMonth(c conf) *cobra.Command {
	cmd := &simple{
		use:   "month",
		a:     "m",
		descr: "Group files by month",
		g:     monthGrouper,
	}
	return newSimpleGroupingCmd(c, cmd)
}

func newYear(c conf) *cobra.Command {
	cmd := &simple{
		use:   "year",
		a:     "y",
		descr: "Group files by year",
		g:     yearGrouper,
	}
	return newSimpleGroupingCmd(c, cmd)
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
