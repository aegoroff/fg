package cmd

import (
	"github.com/spf13/afero"
	"log"
	"os"
	"path/filepath"
)

type grouper struct {
	*filter
	fs       afero.Fs
	basePath string
}

type filter struct {
	include string
	exclude string
}

func newFilter(include string, exclude string) *filter {
	return &filter{
		include: include,
		exclude: exclude,
	}
}

func newGrouper(fs afero.Fs, basePath string, flt *filter) *grouper {
	return &grouper{
		filter:   flt,
		fs:       fs,
		basePath: basePath,
	}
}

func (g *grouper) group(grouper Grouping) error {
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
		if g.filterFile(file.Name()) {
			continue
		}

		// Only files grouped
		g.groupFile(file, grouper)
	}
	return nil
}

func (f *filter) filterFile(file string) bool {
	isInclude := matchPathPattern(f.include, file, true)
	isExclude := matchPathPattern(f.exclude, file, false)

	return !isInclude || isExclude
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

// Returns resultIfError in case of empty pattern or pattern parsing error
func matchPathPattern(pattern string, file string, resultIfError bool) bool {
	result, err := filepath.Match(pattern, file)
	if err != nil {
		log.Printf("%v", err)
		result = resultIfError
	} else if len(pattern) == 0 {
		result = resultIfError
	}
	return result
}
