package main

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

const fgTestFileName = "file.txt"

func TestFilterFile_OnlyIncludeSetFileMatchPattern_ReturnFalse(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	include := "*.txt"
	var exclude string

	// Act
	result := filterFile(fgTestFileName, include, exclude)

	// Assert
	ass.Falsef(result, "File name %s should not be filtered by %s but it was", fgTestFileName, include)
}

func TestFilterFile_OnlyIncludeSetFileNotMatchPattern_ReturnTrue(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	include := "*.mov"
	var exclude string

	// Act
	result := filterFile(fgTestFileName, include, exclude)

	// Assert
	ass.Truef(result, "File name %s should be filtered by %s but it wasn't", fgTestFileName, include)
}

func TestFilterFile_OnlyExcludeSetFileMatchPattern_ReturnTrue(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	var include string
	exclude := "*.txt"

	// Act
	result := filterFile(fgTestFileName, include, exclude)

	// Assert
	ass.Truef(result, "File name %s should be filtered by %s but it wasn't", fgTestFileName, include)
}

func TestFilterFile_OnlyExcludeSetFileNotMatchPattern_ReturnFalse(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	var include string
	exclude := "*.mov"

	// Act
	result := filterFile(fgTestFileName, include, exclude)

	// Assert
	ass.Falsef(result, "File name %s should not filtered by %s but it was", fgTestFileName, include)
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

func TestFg_GroupByExt_FilesMoved(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	const path = "dir"
	const content = "src"
	const f1 = "/f1.txt"
	const f2 = "/f2.txt"
	const f3 = "/f3.xml"
	const f4 = "/f4.html"

	//ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(path, 0755)
	afero.WriteFile(memfs, path+f1, []byte(content), 0644)
	afero.WriteFile(memfs, path+f2, []byte(content), 0644)
	afero.WriteFile(memfs, path+f3, []byte(content), 0644)
	afero.WriteFile(memfs, path+f4, []byte(content), 0644)

	opt := Options{
		Path:    path,
		GroupBy: "ext",
	}

	// Act
	fg(opt, memfs)

	// Assert
	_, err := memfs.Stat(path + "/txt" + f1)
	ass.NoError(err)
	_, err = memfs.Stat(path + "/txt" + f2)
	ass.NoError(err)
	_, err = memfs.Stat(path + "/xml" + f3)
	ass.NoError(err)
	_, err = memfs.Stat(path + "/html" + f4)
	ass.NoError(err)
	_, err = memfs.Stat(path + f1)
	ass.Error(err)
	_, err = memfs.Stat(path + f2)
	ass.Error(err)
	_, err = memfs.Stat(path + f3)
	ass.Error(err)
	_, err = memfs.Stat(path + f4)
	ass.Error(err)
}

func TestFg_GroupByFirst3Letters_FilesMoved(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	const path = "dir"
	const content = "src"
	const f1 = "/file1.txt"
	const f2 = "/file2.txt"
	const f3 = "/dile.xml"
	const f4 = "/eile.html"

	//ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	memfs.MkdirAll(path, 0755)
	afero.WriteFile(memfs, path+f1, []byte(content), 0644)
	afero.WriteFile(memfs, path+f2, []byte(content), 0644)
	afero.WriteFile(memfs, path+f3, []byte(content), 0644)
	afero.WriteFile(memfs, path+f4, []byte(content), 0644)

	opt := Options{
		Path:    path,
		GroupBy: "l3",
	}

	// Act
	fg(opt, memfs)

	// Assert
	_, err := memfs.Stat(path + "/fil" + f1)
	ass.NoError(err)
	_, err = memfs.Stat(path + "/fil" + f2)
	ass.NoError(err)
	_, err = memfs.Stat(path + "/dil" + f3)
	ass.NoError(err)
	_, err = memfs.Stat(path + "/eil" + f4)
	ass.NoError(err)
	_, err = memfs.Stat(path + f1)
	ass.Error(err)
	_, err = memfs.Stat(path + f2)
	ass.Error(err)
	_, err = memfs.Stat(path + f3)
	ass.Error(err)
	_, err = memfs.Stat(path + f4)
	ass.Error(err)
}
