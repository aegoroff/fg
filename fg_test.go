package main

import (
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
