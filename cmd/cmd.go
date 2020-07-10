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

func newSimpleGroupingCmd(c conf, use string, alias string, short string, g grouping) *cobra.Command {
	return newCmd(use, alias, short, func(_ *cobra.Command, _ []string) error {
		flt := newFilter(c.include(), c.exclude())
		grp := newGrouper(c.fs(), c.base(), g)

		return grp.group(flt)
	})
}
