// Various helper code
package main

import (
	"log"
	"path/filepath"
)

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
