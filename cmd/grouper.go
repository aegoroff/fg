package cmd

import (
	"github.com/spf13/afero"
	"log"
	"os"
	"path/filepath"
)

type grouping func(os.FileInfo) []string

type grouper struct {
	*renamer
	basePath string
	grp      grouping
}

type renamer struct {
	fs afero.Fs
}

func newRenamer(fs afero.Fs) *renamer {
	return &renamer{fs: fs}
}

func newGrouper(fs afero.Fs, basePath string, grouping grouping) *grouper {
	return &grouper{
		renamer:  newRenamer(fs),
		basePath: basePath,
		grp:      grouping,
	}
}

func (g *grouper) group(flt *filter) error {
	f, err := g.fs.Open(basePath)
	if err != nil {
		return err
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		// skip directories
		if file.IsDir() {
			continue
		}

		// skip files if necessary
		if flt.filterFile(file.Name()) {
			continue
		}

		// Only files grouped
		g.groupFile(file)
	}
	return nil
}

func (g *grouper) groupFile(file os.FileInfo) {
	// Group key will be subdirectory (of base dir) name
	subdirs := g.grp(file)

	parts := []string{g.basePath}
	parts = append(parts, subdirs...)

	targetDirPath := filepath.Join(parts...)

	// Directory may not exist
	if _, err := g.fs.Stat(targetDirPath); os.IsNotExist(err) {
		if err := g.fs.Mkdir(targetDirPath, os.ModeDir); err != nil {
			log.Printf("%v", err)
			return
		}
	}

	oldFilePath := filepath.Join(g.basePath, file.Name())
	newFilePath := filepath.Join(targetDirPath, file.Name())

	g.rename(oldFilePath, newFilePath)
}

func (r *renamer) rename(oldFilePath string, newFilePath string) {
	if err := r.fs.Rename(oldFilePath, newFilePath); err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("File %s moved to %s", oldFilePath, newFilePath)
	}
}
