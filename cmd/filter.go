package cmd

import (
	"log"
	"path/filepath"
)

type filter struct {
	include string
	exclude string
}

func newFilter(include string, exclude string) *filter {
	return &filter{
		include: include,
		exclude: exclude,
	}
}

func (f *filter) filterFile(file string) bool {
	isInclude := matchPathPattern(f.include, file, true)
	isExclude := matchPathPattern(f.exclude, file, false)

	return !isInclude || isExclude
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
