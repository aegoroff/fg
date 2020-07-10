package cmd

import "github.com/spf13/afero"

type conf interface {
	base() string
	include() string
	exclude() string
	fs() afero.Fs
}

type fgConf struct {
	bp   string
	incl string
	excl string
	afs  afero.Fs
}

func newFgConf(fs afero.Fs) *fgConf {
	return &fgConf{
		afs: fs,
	}
}

func (f *fgConf) base() string { return f.bp }

func (f *fgConf) include() string { return f.incl }

func (f *fgConf) exclude() string { return f.excl }

func (f *fgConf) fs() afero.Fs { return f.afs }
