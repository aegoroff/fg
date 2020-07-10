package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func newExt() *cobra.Command {
	return newCmd("ext", "e", "Group files by file extension", extFunc)
}

func extFunc(_ *cobra.Command, _ []string) error {
	return group(appFileSystem, extGrouper)
}

func extGrouper(file os.FileInfo) []string {
	name := file.Name()
	parts := strings.Split(name, ".")
	if len(parts) < 2 {
		return []string{"no extension"}
	}
	return []string{parts[len(parts)-1]}
}
