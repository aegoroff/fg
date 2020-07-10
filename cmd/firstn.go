package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
)

const numberParamName = "number"

func newFirstn() *cobra.Command {
	cmd := newCmd("firstn", "fn", "Group files by first N letters of a name. By default 3", firstnFunc)

	cmd.Flags().IntP(numberParamName, "n", 3, "The number of first letters that used to group files")
	return cmd
}

func firstnFunc(cmd *cobra.Command, _ []string) error {
	num, err := cmd.Flags().GetInt(numberParamName)
	if err != nil {
		return err
	}
	if num <= 0 {
		return errors.New("number must be positive")
	}

	return group(appFileSystem, func(info os.FileInfo) []string {
		return firstGrouper(num, info)
	})
}

func firstGrouper(num int, file os.FileInfo) []string {
	name := file.Name()
	if len(name) < num {
		return []string{name}
	}
	return []string{name[0:num]}
}
