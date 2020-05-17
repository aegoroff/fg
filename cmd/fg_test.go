package cmd

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const fgTestFileName = "file.txt"

func TestFilterFile_OnlyOneOptionSetFileNotMatchPattern_ReturnFalse(t *testing.T) {
	var tests = []struct {
		include string
		exclude string
	}{
		{"*.txt", ""},
		{"", "*.mov"},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)

		// Act
		result := filterFile(fgTestFileName, test.include, test.exclude)

		// Assert
		ass.Falsef(result, "File name %s should not be filtered but it was", fgTestFileName)
	}
}

func TestFilterFile_OnlyOneOptionSetAndItFiltersFile_ReturnTrue(t *testing.T) {
	var tests = []struct {
		include string
		exclude string
	}{
		{"*.mov", ""},
		{"", "*.txt"},
	}
	for _, test := range tests {
		// Arrange
		ass := assert.New(t)

		// Act
		result := filterFile(fgTestFileName, test.include, test.exclude)

		// Assert
		ass.Truef(result, "File name %s should be filtered but it wasn't", fgTestFileName)
	}
}

func TestFilterFile_BothPatternsEmpty_ReturnFalse(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	var include string
	var exclude string

	// Act
	result := filterFile(fgTestFileName, include, exclude)

	// Assert
	ass.Falsef(result, "File name %s should not filtered because patterns not set but it was", fgTestFileName)
}

func TestFilterFile_BothPatternsSetMatchOnlyInclude_ReturnFalse(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	include := "*.txt"
	exclude := "*.mov"

	// Act
	result := filterFile(fgTestFileName, include, exclude)

	// Assert
	ass.Falsef(result, "File name %s should not filtered by %s but it was", fgTestFileName, exclude)
}

func TestFilterFile_BothPatternsSetMatchBoth_ReturnTrue(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	include := "*.txt"
	exclude := "*.txt"

	// Act
	result := filterFile(fgTestFileName, include, exclude)

	// Assert
	ass.Truef(result, "File name %s should be filtered by %s but it wasn't", fgTestFileName, exclude)
}

func TestFilterFile_BothPatternsSetMatchOnlyExclude_ReturnTrue(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	include := "*.mov"
	exclude := "*.txt"

	// Act
	result := filterFile(fgTestFileName, include, exclude)

	// Assert
	ass.Truef(result, "File name %s should be filtered by %s but it wasn't", fgTestFileName, exclude)
}

func TestFilterFile_BothPatternsSetMatchNoneOfThem_ReturnTrue(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	include := "*.mov"
	exclude := "*.mov"

	// Act
	result := filterFile(fgTestFileName, include, exclude)

	// Assert
	ass.Truef(result, "File name %s should be filtered by include: %s but it wasn't", fgTestFileName, include)
}

func TestFg_GroupingTests_FilesMoved(t *testing.T) {
	y, m, d := time.Now().Date()

	today := "/" + fmt.Sprintf("%d-%02d-%02d", y, m, d)
	month := "/" + fmt.Sprintf("%d-%02d", y, m)
	year := "/" + fmt.Sprintf("%d", y)

	var tests = []struct {
		option string
		dir    string
		file1  string
		file2  string
		file3  string
		file4  string
		sub1   string
		sub2   string
		sub3   string
		sub4   string
	}{
		{"d", "dir", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", today, today, today, today},
		{"day", "dir", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", today, today, today, today},
		{"m", "dir", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", month, month, month, month},
		{"month", "dir", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", month, month, month, month},
		{"y", "dir", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", year, year, year, year},
		{"year", "dir", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", year, year, year, year},
		{"ext", "dir", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", "/txt", "/txt", "/xml", "/html"},
		{"ext", "dir", "/f1", "/f2.txt", "/f3.xml", "/f4.html", "/no extension", "/txt", "/xml", "/html"},
		{"letter", "dir", "/file1.txt", "/file2.txt", "/dile.xml", "/eile.html", "/fil", "/fil", "/dil", "/eil"},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		const content = "src"

		memfs := afero.NewMemMapFs()
		memfs.MkdirAll(test.dir, 0755)
		afero.WriteFile(memfs, test.dir+test.file1, []byte(content), 0644)
		afero.WriteFile(memfs, test.dir+test.file2, []byte(content), 0644)
		afero.WriteFile(memfs, test.dir+test.file3, []byte(content), 0644)
		afero.WriteFile(memfs, test.dir+test.file4, []byte(content), 0644)
		appFileSystem = memfs

		// Act
		rootCmd.SetArgs([]string{test.option, "-p", test.dir})
		rootCmd.Execute()

		// Assert
		_, err := memfs.Stat(test.dir + test.sub1 + test.file1)
		ass.NoError(err)
		_, err = memfs.Stat(test.dir + test.sub2 + test.file2)
		ass.NoError(err)
		_, err = memfs.Stat(test.dir + test.sub3 + test.file3)
		ass.NoError(err)
		_, err = memfs.Stat(test.dir + test.sub4 + test.file4)
		ass.NoError(err)
		_, err = memfs.Stat(test.dir + test.file1)
		ass.Error(err)
		_, err = memfs.Stat(test.dir + test.file2)
		ass.Error(err)
		_, err = memfs.Stat(test.dir + test.file3)
		ass.Error(err)
		_, err = memfs.Stat(test.dir + test.file4)
		ass.Error(err)

		files := getFileNamesInDir(memfs, test.dir)
		ass.Equal(0, len(files), "The number of files in target dont match")
	}
}

func TestFg_UngroupingTests_FilesMoved(t *testing.T) {
	var tests = []struct {
		option string
		dir    string
		file1  string
		file2  string
		file3  string
		sub1   string
		sub2   string
		sub3   string
	}{
		{"u", "dir", "/f1.txt", "/f2.txt", "/f3.xml", "/txt", "/txt", "/xml"},
		{"u", "dir", "/f1.txt", "/f1.txt", "/f3.xml", "/txt", "/txt1", "/xml"},
		{"u", "dir", "/f1.txt", "/f1.txt", "/f1.txt", "/txt", "/txt1", "/xml"},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		const content = "src"

		memfs := afero.NewMemMapFs()
		memfs.MkdirAll(test.dir, 0755)
		afero.WriteFile(memfs, test.dir+test.sub1+test.file1, []byte(content), 0644)
		afero.WriteFile(memfs, test.dir+test.sub2+test.file2, []byte(content), 0644)
		afero.WriteFile(memfs, test.dir+test.sub3+test.file3, []byte(content), 0644)
		appFileSystem = memfs

		// Act
		rootCmd.SetArgs([]string{test.option, "-p", test.dir})
		rootCmd.Execute()

		// Assert
		_, err := memfs.Stat(test.dir + test.sub1 + test.file1)
		ass.Error(err)
		_, err = memfs.Stat(test.dir + test.sub2 + test.file2)
		ass.Error(err)
		_, err = memfs.Stat(test.dir + test.sub3 + test.file3)
		ass.Error(err)
		_, err = memfs.Stat(test.dir + test.file1)
		ass.NoError(err)
		_, err = memfs.Stat(test.dir + test.file2)
		ass.NoError(err)
		_, err = memfs.Stat(test.dir + test.file3)
		ass.NoError(err)

		_, err = memfs.Stat(test.dir + test.sub1)
		ass.NoError(err)

		_, err = memfs.Stat(test.dir + test.sub2)
		ass.NoError(err)

		_, err = memfs.Stat(test.dir + test.sub3)
		ass.NoError(err)

		files := getFileNamesInDir(memfs, test.dir)
		ass.Equal(3, len(files), "The number of files in target dont match")
	}
}

func TestFg_UngroupingTestAndClean_FilesMovedOldDirsRemoved(t *testing.T) {
	var tests = []struct {
		option string
		dir    string
		file1  string
		file2  string
		sub1   string
		sub2   string
	}{
		{"u", "dir", "/f1.txt", "/f2.txt", "/txt", "/txt"},
		{"u", "dir", "/f1.txt", "/f1.txt", "/txt", "/txt1"},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		const content = "src"

		memfs := afero.NewMemMapFs()
		memfs.MkdirAll(test.dir, 0755)
		afero.WriteFile(memfs, test.dir+test.sub1+test.file1, []byte(content), 0644)
		afero.WriteFile(memfs, test.dir+test.sub2+test.file2, []byte(content), 0644)
		appFileSystem = memfs

		// Act
		rootCmd.SetArgs([]string{test.option, "-p", test.dir, "-c"})
		rootCmd.Execute()

		// Assert
		_, err := memfs.Stat(test.dir + test.sub1 + test.file1)
		ass.Error(err)
		_, err = memfs.Stat(test.dir + test.sub2 + test.file2)
		ass.Error(err)
		_, err = memfs.Stat(test.dir + test.file1)
		ass.NoError(err)
		_, err = memfs.Stat(test.dir + test.file2)
		ass.NoError(err)

		_, err = memfs.Stat(test.dir + test.sub1)
		ass.Error(err)

		_, err = memfs.Stat(test.dir + test.sub2)
		ass.Error(err)

		files := getFileNamesInDir(memfs, test.dir)
		ass.Equal(2, len(files), "The number of files in target dont match")
	}
}

func TestFg_UngroupingTestWithFiltering_CountMovedFilesAsSpecified(t *testing.T) {
	var tests = []struct {
		dir     string
		file1   string
		file2   string
		include string
		movedCount int
	}{
		{"dir", "/f1.txt", "/f1.xml", "*.txt", 1},
		{"dir", "/f1.txt", "/f1.xml", "*.exe", 0},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		const content = "src"

		sub := "/sub"
		memfs := afero.NewMemMapFs()
		memfs.MkdirAll(test.dir, 0755)
		afero.WriteFile(memfs, test.dir+sub+test.file1, []byte(content), 0644)
		afero.WriteFile(memfs, test.dir+sub+test.file2, []byte(content), 0644)
		appFileSystem = memfs

		// Act
		rootCmd.SetArgs([]string{"u", "-p", test.dir, "-i", test.include})
		rootCmd.Execute()

		// Assert
		files := getFileNamesInDir(memfs, test.dir)
		ass.Equal(test.movedCount, len(files), "The number of files in target dont match")
	}
}

func getFileNamesInDir(fs afero.Fs, path string) []string {
	base, _ := fs.Open(path)
	defer base.Close()
	items, _ := base.Readdir(-1)
	var files []string
	for _, file := range items {
		if file.IsDir() {
			continue
		}
		files = append(files, file.Name())
	}
	return files
}
