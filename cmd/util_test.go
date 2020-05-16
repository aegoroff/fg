package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const utilTestFileName = "file.txt"

func TestMatchPathPattern_MatchResultIfErrorTrue_ReturnTrue(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	pattern := "*.txt"

	// Act
	result := matchPathPattern(pattern, utilTestFileName, true)

	// Assert
	ass.Truef(result, "File name %s should match pattern: %s but it not matched", utilTestFileName, pattern)
}

func TestMatchPathPattern_MatchResultIfErrorFalse_ReturnTrue(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	pattern := "*.txt"

	// Act
	result := matchPathPattern(pattern, utilTestFileName, false)

	// Assert
	ass.Truef(result, "File name %s should match pattern: %s but it not matched", utilTestFileName, pattern)
}

func TestMatchPathPattern_EmptyPatternResultIfErrorFalse_ReturnFalse(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	var pattern string

	// Act
	result := matchPathPattern(pattern, utilTestFileName, false)

	// Assert
	ass.Falsef(result, "File name %s should not match empty pattern but it matched", utilTestFileName)
}

func TestMatchPathPattern_EmptyPatternResultIfErrorTrue_ReturnTrue(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	var pattern string

	// Act
	result := matchPathPattern(pattern, utilTestFileName, true)

	// Assert
	ass.Truef(result, "File name %s should match empty pattern but it not matched", utilTestFileName)
}

func TestMatchPathPattern_ErrorPatternResultIfErrorFalse_ReturnFalse(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	pattern := "[e-"

	// Act
	result := matchPathPattern(pattern, utilTestFileName, false)

	// Assert
	ass.Falsef(result, "File name %s should not match empty pattern but it matched", utilTestFileName)
}

func TestMatchPathPattern_ErrorPatternResultIfErrorTrue_ReturnTrue(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	pattern := "[e-"

	// Act
	result := matchPathPattern(pattern, utilTestFileName, true)

	// Assert
	ass.Truef(result, "File name %s should match empty pattern but it not matched", utilTestFileName)
}
