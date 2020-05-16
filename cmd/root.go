package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

type Grouping func(os.FileInfo) []string

var appFileSystem = afero.NewOsFs()

const pathParamName = "path"

var sourcesPath string
var include string
var exclude string

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "fg",
	Short: "Grouping files tool",
	Long: `fg is a small commandline app that allows you to easily group
all files in the dir specified into several child subdirectories.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	cobra.MousetrapHelpText = ""
	rootCmd.PersistentFlags().StringVarP(&sourcesPath, pathParamName, "p", "", "REQUIRED. Directory path whose files will be grouped by folders.")
	rootCmd.PersistentFlags().StringVarP(&include, "include", "i", "", "Only files whose names match the pattern specified by the option are grouped.")
	rootCmd.PersistentFlags().StringVarP(&exclude, "exclude", "e", "", "Exclude files whose names match pattern specified by the option from grouping.")
	rootCmd.MarkPersistentFlagRequired(pathParamName)
}

// Execute starts package running
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func group(fs afero.Fs, grouper Grouping) error {
	f, err := fs.Open(sourcesPath)
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
		groupFile(file, sourcesPath, fs, grouper)
	}
	return nil
}

func filterFile(file string, include string, exclude string) bool {
	isInclude := matchPathPattern(include, file, true)
	isExclude := matchPathPattern(exclude, file, false)

	return !isInclude || isExclude
}

func groupFile(file os.FileInfo, baseDirPath string, fs afero.Fs, grouper Grouping)  {
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

	sourcePath := filepath.Join(baseDirPath, file.Name())
	targetPath := filepath.Join(targetDirPath, file.Name())

	if err := fs.Rename(sourcePath, targetPath); err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("File %s moved to %s", sourcePath, targetPath)
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
