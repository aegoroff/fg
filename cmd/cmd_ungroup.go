package cmd

import (
	"github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type fileItem struct {
	path string
	name string
}

const removeEmptyParamName = "clean"

func newUngroup(c conf) *cobra.Command {
	short := "Ungroups file in a directory i.e. copies all files from subdirectories into parent one"
	cmd := newCmd("ungroup", "u", short, func(cmd *cobra.Command, _ []string) error {
		isClean, err := cmd.Flags().GetBool(removeEmptyParamName)
		if err != nil {
			return err
		}

		return ungroup(c, isClean)
	})

	cmd.Flags().BoolP(removeEmptyParamName, "c", false, "Remove empty subdirectories after ungrouping")
	return cmd
}

func ungroup(c conf, isClean bool) error {
	base, err := c.fs().Open(c.base())
	if err != nil {
		return err
	}

	items, err := base.Readdir(-1)
	if err != nil {
		base.Close()
		return err
	}

	subch := make(chan string, 16)

	// Enumerate all subdirs
	go func() {
		defer close(subch)
		// Close base path after reading all subdirs
		defer base.Close()
		for _, item := range items {
			if item.IsDir() {
				subch <- filepath.Join(c.base(), item.Name())
			}
		}
	}()

	filech := make(chan *fileItem, 16)

	flt := newFilter(c.include(), c.exclude())

	// enumerate files in all subdirs
	go func() {
		defer close(filech)
		for sub := range subch {
			s, err := c.fs().Open(sub)
			if err != nil {
				continue
			}

			items, err := s.Readdir(-1)
			if err != nil {
				s.Close()
				continue
			}

			for _, file := range items {
				// skip directories
				if file.IsDir() {
					continue
				}

				// skip files if necessary
				if flt.filterFile(file.Name()) {
					continue
				}

				filech <- &fileItem{path: sub, name: file.Name()}
			}
			s.Close()
		}
	}()

	uniquePaths := make(collections.StringHashSet)
	oldSubDirs := make(collections.StringHashSet)

	r := newRenamer(c.fs())
	// rename files
	for f := range filech {
		oldFilePath := filepath.Join(f.path, f.name)
		newFilePath := filepath.Join(c.base(), f.name)

		oldSubDirs.Add(f.path)

		if uniquePaths.Contains(newFilePath) {
			newFilePath = createNewPath(c.base(), oldFilePath)
		}
		r.rename(oldFilePath, newFilePath)
		uniquePaths.Add(newFilePath)
	}

	// cleanup old dirs
	if isClean {
		removeDirectories(c.fs(), oldSubDirs.Items())
	}

	return nil
}

func createNewPath(basePath, oldFilePath string) string {
	d, f := filepath.Split(oldFilePath)
	sep := string(os.PathSeparator)
	baseDirParts := strings.Split(strings.Trim(basePath, sep), sep)
	dirParts := strings.Split(strings.Trim(d, sep), sep)
	newNameParts := append(dirParts[len(baseDirParts):], f)

	return filepath.Join(basePath, strings.Join(newNameParts, "-"))
}

func removeDirectories(fs afero.Fs, oldSubDirs []string) {
	for _, k := range oldSubDirs {
		if !isDirEmpty(fs, k) {
			continue
		}
		err := fs.Remove(k)
		if err != nil {
			log.Printf("%v", err)
		}
	}
}

func isDirEmpty(fs afero.Fs, path string) bool {
	base, err := fs.Open(path)
	if err != nil {
		return false
	}

	defer base.Close()
	items, err := base.Readdir(-1)
	if err != nil {
		return false
	}
	for _, file := range items {
		if file.IsDir() {
			continue
		}
		return false
	}
	return true
}
