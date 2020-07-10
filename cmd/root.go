package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Grouping function
type Grouping func(os.FileInfo) []string

var appFileSystem = afero.NewOsFs()
var appWriter io.Writer

const pathParamName = "path"

var basePath string
var include string
var exclude string

func newRoot() *cobra.Command {
	return &cobra.Command{
		Use:   "fg",
		Short: "Grouping files tool",
		Long: ` fg is a small commandline app that allows you to easily group
 all files in the dir specified into several child subdirectories.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}

func init() {
	cobra.MousetrapHelpText = ""
	appWriter = os.Stdout
}

// Execute starts package running
func Execute(args ...string) error {
	rootCmd := newRoot()

	if args != nil && len(args) > 0 {
		rootCmd.SetArgs(args)
	}

	rootCmd.PersistentFlags().StringVarP(&basePath, pathParamName, "p", "", "REQUIRED. Directory path whose files will be grouped by folders.")
	rootCmd.PersistentFlags().StringVarP(&include, "include", "i", "", "Only files whose names match the pattern specified by the option are grouped.")
	rootCmd.PersistentFlags().StringVarP(&exclude, "exclude", "e", "", "Exclude files whose names match pattern specified by the option from grouping.")

	rootCmd.AddCommand(newDay())
	rootCmd.AddCommand(newExt())
	rootCmd.AddCommand(newFirstn())
	rootCmd.AddCommand(newMonth())
	rootCmd.AddCommand(newUngroup())
	rootCmd.AddCommand(newVersion())
	rootCmd.AddCommand(newYear())
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func group(fs afero.Fs, grouper Grouping) error {
	f, err := fs.Open(basePath)
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
		if filterFile(file.Name(), include, exclude) {
			continue
		}

		// Only files grouped
		groupFile(file, basePath, fs, grouper)
	}
	return nil
}

func filterFile(file string, include string, exclude string) bool {
	isInclude := matchPathPattern(include, file, true)
	isExclude := matchPathPattern(exclude, file, false)

	return !isInclude || isExclude
}

func groupFile(file os.FileInfo, baseDirPath string, fs afero.Fs, grouper Grouping) {
	// Group key will be subdirectory (of base dir) name
	subdirs := grouper(file)

	parts := []string{baseDirPath}
	parts = append(parts, subdirs...)

	targetDirPath := filepath.Join(parts...)

	// Directory may not exist
	if _, err := fs.Stat(targetDirPath); os.IsNotExist(err) {
		if err := fs.Mkdir(targetDirPath, os.ModeDir); err != nil {
			log.Printf("%v", err)
			return
		}
	}

	oldFilePath := filepath.Join(baseDirPath, file.Name())
	newFilePath := filepath.Join(targetDirPath, file.Name())

	rename(fs, oldFilePath, newFilePath)
}

func rename(fs afero.Fs, oldFilePath string, newFilePath string) {
	if err := fs.Rename(oldFilePath, newFilePath); err != nil {
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
