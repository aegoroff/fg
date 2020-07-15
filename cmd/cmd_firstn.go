package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
)

const numberParamName = "number"

func newFirstn(c conf) *cobra.Command {
	short := "Group files by first N letters of a name. By default 3"
	cmd := newCmd("firstn", "fn", short, func(cmd *cobra.Command, _ []string) error {
		num, err := cmd.Flags().GetInt(numberParamName)
		if err != nil {
			return err
		}
		if num <= 0 {
			return errors.New("number must be positive")
		}

		flt := NewFilter(c.include(), c.exclude())
		g := newGrouper(c.fs(), c.root(), func(info os.FileInfo) []string {
			return firstGrouper(num, info)
		})

		return g.group(flt)
	})

	cmd.Flags().IntP(numberParamName, "n", 3, "The number of first letters that used to group files")
	return cmd
}

func firstGrouper(num int, file os.FileInfo) []string {
	name := file.Name()
	if len(name) < num {
		return []string{name}
	}
	return []string{name[0:num]}
}
