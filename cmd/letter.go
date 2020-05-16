package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

const numberParamName = "number"

// letterCmd represents the letter command
var letterCmd = &cobra.Command{
	Use:     "letter",
	Aliases: []string{},
	Short:   "Group files by first N letters of a name. By default 3",
	RunE: func(cmd *cobra.Command, args []string) error {

		num, err := cmd.Flags().GetInt(numberParamName)
		if err != nil {
			return err
		}

		fg(appFileSystem, func(info os.FileInfo) []string {
			return firstGrouper(num, info)
		})
		return nil
	},
}

func init() {
	rootCmd.AddCommand(letterCmd)
	letterCmd.Flags().IntP(numberParamName, "n", 3, "The number of first letters that used to group files")
}

func firstGrouper(num int, file os.FileInfo) []string {
	name := file.Name()
	if len(name) < num {
		return []string{name}
	}
	return []string{name[0:num]}
}
