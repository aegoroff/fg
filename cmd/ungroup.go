package cmd

import (
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

// ungroupCmd represents the ungroup command
var ungroupCmd = &cobra.Command{
	Use:     "ungroup",
	Aliases: []string{"u"},
	Short:   "Ungroups file in a directory i.e. copies all files from subdirectories into parent one",
	RunE: func(cmd *cobra.Command, args []string) error {
		isClean, err := cmd.Flags().GetBool(removeEmptyParamName)
		if err != nil {
			return err
		}

		return ungroup(appFileSystem, isClean)
	},
}

func init() {
	rootCmd.AddCommand(ungroupCmd)
	ungroupCmd.Flags().BoolP(removeEmptyParamName, "c", false, "Remove empty subdirectories after ungrouping")
}

func ungroup(fs afero.Fs, isClean bool) error {
	base, err := fs.Open(basePath)
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
				subch <- filepath.Join(basePath, item.Name())
			}
		}
	}()

	filech := make(chan *fileItem, 16)

	// enumerate files in all subdirs
	go func() {
		defer close(filech)
		for sub := range subch {
			s, err := fs.Open(sub)
			if err != nil {
				continue
			}
			defer s.Close()

			items, err := s.Readdir(-1)
			if err != nil {
				continue
			}

			for _, file := range items {
				// skip directories
				if file.IsDir() {
					continue
				}

				// skip files if necessary
				if filterFile(file.Name(), include, exclude) {
					continue
				}

				filech <- &fileItem{path: sub, name: file.Name()}
			}
		}
	}()

	uniquePaths := make(map[string]interface{})
	oldSubDirs := make(map[string]interface{})

	// rename files
	for f := range filech {
		oldFilePath := filepath.Join(f.path, f.name)
		newFilePath := filepath.Join(basePath, f.name)

		if _, ok := oldSubDirs[f.path]; !ok {
			oldSubDirs[f.path] = nil
		}

		if _, ok := uniquePaths[newFilePath]; ok {
			d, f := filepath.Split(oldFilePath)
			sep := string(os.PathSeparator)
			baseDirParts := strings.Split(strings.Trim(basePath, sep), sep)
			dirParts := strings.Split(strings.Trim(d, sep), sep)
			newNameParts := append(dirParts[len(baseDirParts):], f)

			newFilePath = filepath.Join(basePath, strings.Join(newNameParts, "-"))
		}
		rename(fs, oldFilePath, newFilePath)
		uniquePaths[newFilePath] = nil
	}

	// cleanup old dirs
	if isClean {
		for k, _ := range oldSubDirs {
			if !isDirEmpty(fs, k) {
				continue
			}
			err = fs.Remove(k)
			if err != nil {
				log.Printf("%v", err)
			}
		}
	}

	return nil
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
