package utils

import (
	"dharlequin/google-calendar-converter/model"
	"log"
	"sort"
	"strconv"
	"strings"
)

func ValidateFilePath(name string) {
	if name == "" {
		log.Fatalln("Value cannot be empty")
	}
	if !strings.HasSuffix(name, ".ics") {
		log.Fatalln("This is probably not a Google Calendar file (.ics)")
	}
}

func ExtractStringValue(line string, target string) string {
	return strings.ReplaceAll(line, target, "")
}

func ConvertToNumber(sValue string) int {
	number, err := strconv.Atoi(sValue)
	HandleError(err)

	return number
}

func ValidateMonth(month int) {
	if month < 1 && month > 12 {
		log.Fatalln("Invalid month parameter, must be between 1 and 12")
	}
}

func ValidateYear(year string) {
	if len([]rune(year)) != 4 {
		log.Fatalln("Invalid year parameter, must contain 4 digits")
	}
}

func SortCategories(release *model.Release) {
	sortCategory(release.Movies)
	sortCategory(release.Shows)
	sortCategory(release.Games)
	sortCategory(release.Other)
}

// sorts items in a release category by date preserving titles order within one date
func sortCategory(items []model.ReleaseItem) {
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Title < items[j].Title
	})

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Date.Before(items[j].Date)
	})
}
