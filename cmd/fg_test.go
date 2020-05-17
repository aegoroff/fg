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
		file1  string
		file2  string
		file3  string
		file4  string
		sub1   string
		sub2   string
		sub3   string
		sub4   string
	}{
		{"d", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", today, today, today, today},
		{"day", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", today, today, today, today},
		{"m", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", month, month, month, month},
		{"month", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", month, month, month, month},
		{"y", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", year, year, year, year},
		{"year", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", year, year, year, year},
		{"ext", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", "/txt", "/txt", "/xml", "/html"},
		{"e", "/f1.txt", "/f2.txt", "/f3.xml", "/f4.html", "/txt", "/txt", "/xml", "/html"},
		{"ext", "/f1", "/f2.txt", "/f3.xml", "/f4.html", "/no extension", "/txt", "/xml", "/html"},
		{"firstn", "/file1.txt", "/file2.txt", "/dile.xml", "/eile.html", "/fil", "/fil", "/dil", "/eil"},
		{"fn", "/file1.txt", "/file2.txt", "/dile.xml", "/eile.html", "/fil", "/fil", "/dil", "/eil"},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		const content = "src"
		dir := "dir"

		memfs := afero.NewMemMapFs()
		memfs.MkdirAll(dir, 0755)
		afero.WriteFile(memfs, dir+test.file1, []byte(content), 0644)
		afero.WriteFile(memfs, dir+test.file2, []byte(content), 0644)
		afero.WriteFile(memfs, dir+test.file3, []byte(content), 0644)
		afero.WriteFile(memfs, dir+test.file4, []byte(content), 0644)
		appFileSystem = memfs

		// Act
		rootCmd.SetArgs([]string{test.option, "-p", dir, "-i", ""})
		rootCmd.Execute()

		// Assert
		_, err := memfs.Stat(dir + test.sub1 + test.file1)
		ass.NoError(err)
		_, err = memfs.Stat(dir + test.sub2 + test.file2)
		ass.NoError(err)
		_, err = memfs.Stat(dir + test.sub3 + test.file3)
		ass.NoError(err)
		_, err = memfs.Stat(dir + test.sub4 + test.file4)
		ass.NoError(err)
		_, err = memfs.Stat(dir + test.file1)
		ass.Error(err)
		_, err = memfs.Stat(dir + test.file2)
		ass.Error(err)
		_, err = memfs.Stat(dir + test.file3)
		ass.Error(err)
		_, err = memfs.Stat(dir + test.file4)
		ass.Error(err)

		files := getFileNamesInDir(memfs, dir)
		ass.Equal(0, len(files), "The number of files in target dont match")
	}
}

func TestFg_UngroupingTests_FilesMoved(t *testing.T) {
	var tests = []struct {
		option string
		file1  string
		file2  string
		file3  string
		sub1   string
		sub2   string
		sub3   string
	}{
		{"u", "/f1.txt", "/f2.txt", "/f3.xml", "/txt", "/txt", "/xml"},
		{"ungroup", "/f1.txt", "/f2.txt", "/f3.xml", "/txt", "/txt", "/xml"},
		{"u", "/f1.txt", "/f1.txt", "/f3.xml", "/txt", "/txt1", "/xml"},
		{"u", "/f1.txt", "/f1.txt", "/f1.txt", "/txt", "/txt1", "/xml"},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		const content = "src"
		dir := "dir"

		memfs := afero.NewMemMapFs()
		memfs.MkdirAll(dir, 0755)
		afero.WriteFile(memfs, dir+test.sub1+test.file1, []byte(content), 0644)
		afero.WriteFile(memfs, dir+test.sub2+test.file2, []byte(content), 0644)
		afero.WriteFile(memfs, dir+test.sub3+test.file3, []byte(content), 0644)
		appFileSystem = memfs

		// Act
		rootCmd.SetArgs([]string{test.option, "-p", dir, "-i", ""})
		rootCmd.Execute()

		// Assert
		_, err := memfs.Stat(dir + test.sub1 + test.file1)
		ass.Error(err)
		_, err = memfs.Stat(dir + test.sub2 + test.file2)
		ass.Error(err)
		_, err = memfs.Stat(dir + test.sub3 + test.file3)
		ass.Error(err)
		_, err = memfs.Stat(dir + test.file1)
		ass.NoError(err)
		_, err = memfs.Stat(dir + test.file2)
		ass.NoError(err)
		_, err = memfs.Stat(dir + test.file3)
		ass.NoError(err)

		_, err = memfs.Stat(dir + test.sub1)
		ass.NoError(err)

		_, err = memfs.Stat(dir + test.sub2)
		ass.NoError(err)

		_, err = memfs.Stat(dir + test.sub3)
		ass.NoError(err)

		files := getFileNamesInDir(memfs, dir)
		ass.Equal(3, len(files), "The number of files in target dont match")
	}
}

func TestFg_UngroupingTestAndClean_FilesMovedOldDirsRemoved(t *testing.T) {
	var tests = []struct {
		file1 string
		file2 string
		sub1  string
		sub2  string
	}{
		{"/f1.txt", "/f2.txt", "/txt", "/txt"},
		{"/f1.txt", "/f1.txt", "/txt", "/txt1"},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		const content = "src"
		dir := "dir"

		memfs := afero.NewMemMapFs()
		memfs.MkdirAll(dir, 0755)
		afero.WriteFile(memfs, dir+test.sub1+test.file1, []byte(content), 0644)
		afero.WriteFile(memfs, dir+test.sub2+test.file2, []byte(content), 0644)
		appFileSystem = memfs

		// Act
		rootCmd.SetArgs([]string{"u", "-p", dir, "-c", "-i", ""})
		rootCmd.Execute()

		// Assert
		_, err := memfs.Stat(dir + test.sub1 + test.file1)
		ass.Error(err)
		_, err = memfs.Stat(dir + test.sub2 + test.file2)
		ass.Error(err)
		_, err = memfs.Stat(dir + test.file1)
		ass.NoError(err)
		_, err = memfs.Stat(dir + test.file2)
		ass.NoError(err)

		_, err = memfs.Stat(dir + test.sub1)
		ass.Error(err)

		_, err = memfs.Stat(dir + test.sub2)
		ass.Error(err)

		files := getFileNamesInDir(memfs, dir)
		ass.Equal(2, len(files), "The number of files in target dont match")
	}
}

func TestFg_UngroupingTestWithFiltering_CountMovedFilesAsSpecified(t *testing.T) {
	var tests = []struct {
		file1      string
		file2      string
		include    string
		movedCount int
	}{
		{"/f1.txt", "/f1.xml", "*.txt", 1},
		{"/f1.txt", "/f1.xml", "*.exe", 0},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		const content = "src"
		dir := "dir"

		sub := "/sub"
		memfs := afero.NewMemMapFs()
		memfs.MkdirAll(dir, 0755)
		afero.WriteFile(memfs, dir+sub+test.file1, []byte(content), 0644)
		afero.WriteFile(memfs, dir+sub+test.file2, []byte(content), 0644)
		appFileSystem = memfs

		// Act
		rootCmd.SetArgs([]string{"u", "-p", dir, "-i", test.include})
		rootCmd.Execute()

		// Assert
		files := getFileNamesInDir(memfs, dir)
		ass.Equal(test.movedCount, len(files), "The number of files in target dont match")
	}
}

func TestFg_UngroupingTestWithFilteringAndCleaning_CountMovedFilesAsSpecifiedNotEmptySubdirsExist(t *testing.T) {
	var tests = []struct {
		file1      string
		file2      string
		include    string
		movedCount int
	}{
		{"/f1.txt", "/f1.xml", "*.txt", 1},
		{"/f1.txt", "/f2.txt", "f1.txt", 1},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		const content = "src"
		dir := "dir"

		sub := "/sub"
		memfs := afero.NewMemMapFs()
		memfs.MkdirAll(dir, 0755)
		afero.WriteFile(memfs, dir+sub+test.file1, []byte(content), 0644)
		afero.WriteFile(memfs, dir+sub+test.file2, []byte(content), 0644)
		appFileSystem = memfs

		// Act
		rootCmd.SetArgs([]string{"u", "-p", dir, "-i", test.include, "-c"})
		rootCmd.Execute()

		// Assert
		files := getFileNamesInDir(memfs, dir)
		ass.Equal(test.movedCount, len(files), "The number of files in target dont match")

		_, err := memfs.Stat(dir + sub)
		ass.NoError(err)
	}
}

func TestFg_UngroupingTestReadOnlyTarget_CountMovedFilesAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	const content = "src"
	dir := "dir"
	sub := "/sub"
	file1 := "/f1.txt"
	file2 := "/f1.xml"

	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(dir, 0755)
	afero.WriteFile(memfs, dir+sub+file1, []byte(content), 0644)
	afero.WriteFile(memfs, dir+sub+file2, []byte(content), 0644)

	appFileSystem = afero.NewReadOnlyFs(memfs)

	// Act
	rootCmd.SetArgs([]string{"u", "-p", dir, "-c", "-i", ""})
	rootCmd.Execute()

	// Assert
	files := getFileNamesInDir(memfs, dir)
	ass.Equal(0, len(files), "The number of files in target dont match")

	_, err := memfs.Stat(dir + sub)
	ass.NoError(err)
}

func TestFg_UngroupingTestSubdirWithoutFiles_CountMovedFilesAsSpecifiedSubdirWithoutFilesNotRemoved(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	const content = "src"
	dir := "dir1"
	sub := "/dub"
	sub1 := "/dub1"
	file1 := "/f3.txt"
	file2 := "/f3.xml"

	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(dir, 0755)
	memfs.MkdirAll(dir+sub1, 0755)
	afero.WriteFile(memfs, dir+sub+file1, []byte(content), 0644)
	afero.WriteFile(memfs, dir+sub+file2, []byte(content), 0644)

	appFileSystem = memfs

	// Act
	rootCmd.SetArgs([]string{"u", "-p", dir, "-c", "-i", ""})
	rootCmd.Execute()

	// Assert
	files := getFileNamesInDir(memfs, dir)
	ass.Equal(2, len(files), "The number of files in target dont match")

	_, err := memfs.Stat(dir + sub1)
	ass.NoError(err)

	_, err = memfs.Stat(dir + sub)
	ass.Error(err)
}

func TestFg_GroupingTestReadOnlyTarget_CountNotMovedFilesAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	const content = "src"
	dir := "dir"
	sub := "/txt"
	file1 := "/f1.txt"
	file2 := "/f2.txt"

	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(dir, 0755)
	afero.WriteFile(memfs, dir+file1, []byte(content), 0644)
	afero.WriteFile(memfs, dir+file2, []byte(content), 0644)

	appFileSystem = afero.NewReadOnlyFs(memfs)

	// Act
	rootCmd.SetArgs([]string{"e", "-p", dir, "-i", ""})
	rootCmd.Execute()

	// Assert
	files := getFileNamesInDir(memfs, dir)
	ass.Equal(2, len(files), "The number of files in target dont match")

	_, err := memfs.Stat(dir + sub)
	ass.Error(err)
}

func TestFg_GroupingTestFirstNFileNameShort_CountMovedFilesAsSpecifiedTargetPathAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	const content = "src"
	dir := "dir"
	file1 := "/f1.t"
	file2 := "/f2.txt"

	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(dir, 0755)
	afero.WriteFile(memfs, dir+file1, []byte(content), 0644)
	afero.WriteFile(memfs, dir+file2, []byte(content), 0644)

	appFileSystem = memfs

	// Act
	rootCmd.SetArgs([]string{"fn", "-p", dir, "-i", "", "-n", "5"})
	rootCmd.Execute()

	// Assert
	files := getFileNamesInDir(memfs, dir)
	ass.Equal(0, len(files), "The number of files in target dont match")

	files = getFileNamesInDir(memfs, dir + "/f1.t/")
	ass.Equal(1, len(files), "The number of files in target dont match")

	files = getFileNamesInDir(memfs, dir + "/f2.tx/")
	ass.Equal(1, len(files), "The number of files in target dont match")
}

func TestFg_GroupingTestFirstNFileInvalidNum_FilesNotMoved(t *testing.T) {
	var tests = []struct {
		num      string
	}{
		{"-1" },
		{"0" },
		{"xxx" },
	}
	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		const content = "src"
		dir := "dir"
		file1 := "/f1.txt"
		file2 := "/f2.txt"

		memfs := afero.NewMemMapFs()
		memfs.MkdirAll(dir, 0755)
		afero.WriteFile(memfs, dir+file1, []byte(content), 0644)
		afero.WriteFile(memfs, dir+file2, []byte(content), 0644)

		appFileSystem = memfs

		// Act
		rootCmd.SetArgs([]string{"fn", "-p", dir, "-i", "", "-n", test.num})
		rootCmd.Execute()

		// Assert
		files := getFileNamesInDir(memfs, dir)
		ass.Equal(2, len(files), "The number of files in target dont match")
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
