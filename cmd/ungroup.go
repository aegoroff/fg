package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

// ungroupCmd represents the ungroup command
var ungroupCmd = &cobra.Command{
	Use:     "ungroup",
	Aliases: []string{"u"},
	Short:   "Ungroups file in a directory i.e. copies all files from subdirectories into parent one",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ungroup(appFileSystem)
	},
}

func init() {
	rootCmd.AddCommand(ungroupCmd)
}

func ungroup(fs afero.Fs) error {
	f, err := fs.Open(sourcesPath)
	if err != nil {
		return err
	}
	defer f.Close()

	items, err := f.Readdir(-1)
	if err != nil {
		return err
	}

	subdirs := []string{}
	for _, file := range items {
		// skip directories
		if file.IsDir() {
			subdirs = append(subdirs, file.Name())
		}
	}

	for _, d := range subdirs {
		sub := filepath.Join(sourcesPath, d)
		s, err := fs.Open(sub)
		if err != nil {
			continue
		}

		items, err := s.Readdir(-1)
		if err != nil {
			continue
		}

		for _, file := range items {
			// skip directories
			if file.IsDir() {
				continue
			}

			oldFilePath := filepath.Join(sourcesPath, d, file.Name())
			newFilePath := filepath.Join(sourcesPath, file.Name())

			if err := fs.Rename(oldFilePath, newFilePath); err != nil {
				log.Printf("%v", err)
			} else {
				log.Printf("File %s moved to %s", oldFilePath, newFilePath)
			}
		}

		s.Close()
	}

	return nil
}
