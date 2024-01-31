package utils

import (
	"dharlequin/google-calendar-converter/model"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

const LOG_MESSAGE_DELIMITER = "--"
const LOG_RECORD_DELIMITER = "  "

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

func NormalizeString(value string) string {
	// replace "/," when using commas in Google Calendar event title
	return strings.Replace(value, "\\,", ",", -1)
}

func StringToNumber(sValue string) int {
	number, err := strconv.Atoi(sValue)
	HandleError(err)

	return number
}

func NumberToString(iValue int) string {
	return strconv.FormatInt(int64(iValue), 10)
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

func AddEmptyLine() {
	fmt.Println()
}

func PrintParserLog(message string, record string, level int) {
	var messagePrefix string
	var recordPrefix string

	for i := 0; i < level; i++ {
		messagePrefix += LOG_MESSAGE_DELIMITER
		recordPrefix += LOG_RECORD_DELIMITER
	}

	fmt.Println(messagePrefix, message)

	if record != "" {
		fmt.Println(recordPrefix, record)
	}
}
