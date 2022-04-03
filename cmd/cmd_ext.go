package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"path"
	"strings"
)

func newExt(c conf) *cobra.Command {
	cmd := &simple{
		use:   "ext",
		a:     "e",
		descr: "Group files by file extension",
		g:     extGrouper,
	}
	return newSimpleGroupingCmd(c, cmd)
}

func extGrouper(file os.FileInfo) []string {
	ext := path.Ext(file.Name())
	if len(ext) == 0 {
		return []string{"no extension"}
	}
	return []string{strings.TrimPrefix(ext, ".")}
}
