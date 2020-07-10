package cmd

import "github.com/spf13/cobra"

type cmdFunc func(cmd *cobra.Command, args []string) error

func newCmd(use string, alias string, short string, f cmdFunc) *cobra.Command {
	return &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   short,
		RunE:    f,
	}
}
