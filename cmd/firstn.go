package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
)

const numberParamName = "number"

// firstnCmd represents the firstn command
var firstnCmd = &cobra.Command{
	Use:     "firstn",
	Aliases: []string{"fn"},
	Short:   "Group files by first N letters of a name. By default 3",
	RunE: func(cmd *cobra.Command, args []string) error {

		num, err := cmd.Flags().GetInt(numberParamName)
		if err != nil {
			return err
		}
		if num <= 0 {
			return errors.New("Number must be positive")
		}

		return group(appFileSystem, func(info os.FileInfo) []string {
			return firstGrouper(num, info)
		})
	},
}

func init() {
	rootCmd.AddCommand(firstnCmd)
	firstnCmd.Flags().IntP(numberParamName, "n", 3, "The number of first letters that used to group files")
}

func firstGrouper(num int, file os.FileInfo) []string {
	name := file.Name()
	if len(name) < num {
		return []string{name}
	}
	return []string{name[0:num]}
}
