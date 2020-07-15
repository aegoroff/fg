package cmd

import (
	"log"
	"path/filepath"
)

// Filter defines file filtering interface
type Filter interface {
	// Skip filters file specified if necessary
	Skip(file string) bool
}

// NewFilter creates new filter
func NewFilter(include string, exclude string) Filter {
	return &filter{
		incl: newIncluder(include),
		excl: newExcluder(exclude),
	}
}

// Skip filters file specified if necessary
func (f *filter) Skip(file string) bool {
	return !f.incl.match(file) || f.excl.match(file)
}

type filter struct {
	incl matcher
	excl matcher
}

type matcher interface {
	match(file string) bool
}

func newExcluder(pattern string) matcher {
	return &matching{pattern: pattern, resultIfError: false}
}

func newIncluder(pattern string) matcher {
	return &matching{pattern: pattern, resultIfError: true}
}

type matching struct {
	pattern       string
	resultIfError bool
}

// Returns resultIfError in case of empty pattern or pattern parsing error
func (m *matching) match(file string) bool {
	result, err := filepath.Match(m.pattern, file)
	if err != nil {
		log.Printf("%v", err)
		result = m.resultIfError
	} else if len(m.pattern) == 0 {
		result = m.resultIfError
	}
	return result
}
