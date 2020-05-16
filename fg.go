// This tool groups all files in the dir specified into several child directories. Grouping uses file modification time.
package main

import (
	"github.com/spf13/afero"
	"github.com/voxelbrain/goptions"
	"log"
	"os"
	"path/filepath"
)

type Options struct {
	Help    goptions.Help `goptions:"-h, --help, description='Show this help'"`
	Path    string        `goptions:"-p, --path, obligatory, description='Name to the directory whose files will be grouped by folders.'"`
	GroupBy string        `goptions:"-g, --groupby, description='Grouping mode. Only: day or d, month or m, year or y, ext (file extension), l3 (first 3 letters of a name) supported. If not set day used'"`
	Include string        `goptions:"-i, --include, description='Only files whose names match the pattern specified by the option are grouped.'"`
	Exclude string        `goptions:"-e, --exclude, description='Exclude files whose names match pattern specified by the option from grouping.'"`
}

func main() {
	options := Options{}

	goptions.ParseAndFail(&options)
	fs := afero.NewOsFs()
	if _, err := fs.Stat(options.Path); os.IsNotExist(err) {
		log.Fatalf("Directory '%s' does not exist. Details:\n  %v", options.Path, err)
	}

	fg(options, fs)
}

func fg(options Options, fs afero.Fs) {
	f, err := fs.Open(options.Path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		// skip directories
		if file.IsDir() {
			continue
		}

		// skip files if necessary
		if filterFile(file.Name(), options.Include, options.Exclude) {
			continue
		}

		// Only files grouped
		groupFile(file, options.Path, options.GroupBy, fs)
	}
}

func filterFile(file string, include string, exclude string) bool {
	isInclude := matchPathPattern(include, file, true)
	isExclude := matchPathPattern(exclude, file, false)

	return !isInclude || isExclude
}

func groupFile(file os.FileInfo, baseDirPath string, groupBy string, fs afero.Fs) {
	// Group key will be subdirectory (of base dir) name
	grpKey := getGroupKeyFromFileObject(file, groupBy)

	targetDirPath := filepath.Join(baseDirPath, grpKey)

	// Directory may not exist
	if _, err := fs.Stat(targetDirPath); os.IsNotExist(err) {
		if err := fs.Mkdir(targetDirPath, os.ModeDir); err != nil {
			log.Printf("%v", err)
			return
		}
	}

	sourcePath := filepath.Join(baseDirPath, file.Name())
	targetPath := filepath.Join(targetDirPath, file.Name())

	if err := fs.Rename(sourcePath, targetPath); err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("File %s moved to %s", sourcePath, targetPath)
	}
}
