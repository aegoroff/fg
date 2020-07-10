package cmd

import (
	"github.com/spf13/afero"
	"log"
	"os"
	"path/filepath"
)

type grouper struct {
	fs       afero.Fs
	basePath string
}

func newGrouper(fs afero.Fs, basePath string) *grouper {
	return &grouper{
		fs:       fs,
		basePath: basePath,
	}
}

func (g *grouper) group(grouper Grouping, flt *filter) error {
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
		g.groupFile(file, grouper)
	}
	return nil
}

func (g *grouper) groupFile(file os.FileInfo, grouper Grouping) {
	// Group key will be subdirectory (of base dir) name
	subdirs := grouper(file)

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

func (g *grouper) rename(oldFilePath string, newFilePath string) {
	if err := g.fs.Rename(oldFilePath, newFilePath); err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("File %s moved to %s", oldFilePath, newFilePath)
	}
}
