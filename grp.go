// This file contains all grouping calculation code
package main

import (
	"fmt"
	"os"
)

const (
	Day        = "day"
	DayShort   = "d"
	Month      = "month"
	MonthShort = "m"
	Year       = "year"
	YearShort  = "y"
)

type Grouping func(os.FileInfo) string

var groupings = map[string]Grouping{
	Day:        getGroupingKeyDayAnsi,
	DayShort:   getGroupingKeyDayAnsi,
	Month:      getGroupingKeyMonthAnsi,
	MonthShort: getGroupingKeyMonthAnsi,
	Year:       getGroupingKeyYear,
	YearShort:  getGroupingKeyYear,
}

// Gets grouping key from file object specified
func getGroupKeyFromFileObject(file os.FileInfo, groupBy string) string {
	if action, ok := groupings[groupBy]; ok {
		return action(file)
	} else {
		return groupings[Day](file)
	}
}

func getGroupingKeyDayAnsi(file os.FileInfo) string {
	year, month, day := file.ModTime().Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func getGroupingKeyMonthAnsi(file os.FileInfo) string {
	year, month, _ := file.ModTime().Date()
	return fmt.Sprintf("%d-%02d", year, month)
}

func getGroupingKeyYear(file os.FileInfo) string {
	year, _, _ := file.ModTime().Date()
	return fmt.Sprintf("%d", year)
}
