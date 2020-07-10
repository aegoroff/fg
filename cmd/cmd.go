package cmd

import "github.com/spf13/cobra"

type simpler interface {
	name() string
	alias() string
	short() string
	grouper() grouping
}

type simple struct {
	use   string
	a     string
	descr string
	g     grouping
}

func (s *simple) name() string { return s.use }

func (s *simple) alias() string { return s.a }

func (s *simple) short() string { return s.descr }

func (s *simple) grouper() grouping { return s.g }

type cmdFunc func(cmd *cobra.Command, args []string) error

func newCmd(use string, alias string, short string, f cmdFunc) *cobra.Command {
	return &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   short,
		RunE:    f,
	}
}

func newSimpleGroupingCmd(c conf, s simpler) *cobra.Command {
	return newCmd(s.name(), s.alias(), s.short(), func(_ *cobra.Command, _ []string) error {
		flt := newFilter(c.include(), c.exclude())
		grp := newGrouper(c.fs(), c.root(), s.grouper())

		return grp.group(flt)
	})
}
