package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func newExt(c conf) *cobra.Command {
	return newSimpleGroupingCmd(c, "ext", "e", "Group files by file extension", extGrouper)
}

func extGrouper(file os.FileInfo) []string {
	name := file.Name()
	parts := strings.Split(name, ".")
	if len(parts) < 2 {
		return []string{"no extension"}
	}
	return []string{parts[len(parts)-1]}
}
