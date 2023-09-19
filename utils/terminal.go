package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetFilePath(instruction string) string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(instruction)
	scanner.Scan()
	path := scanner.Text()

	validateFilePath(path)
	return path
}

func validateFilePath(name string) {
	if name == "" {
		log.Fatalln("Value cannot be empty")
	}
	if !strings.HasSuffix(name, ".ics") {
		log.Fatalln("This is probably not a Google Calendar file (.ics)")
	}
}

func GetYearFromUser(instruction string) int {
	sYear := getNumberFromUser(instruction)

	if sYear == "" {
		return time.Now().Year()
	}

	validateYear(sYear)

	return convertToNumber(sYear)
}

func GetMonthFromUser(instruction string) int {
	sMonth := getNumberFromUser(instruction)

	if sMonth == "" {
		return int(time.Now().Month())
	}

	validateMonth(sMonth)

	return convertToNumber(sMonth)
}

func getNumberFromUser(instruction string) string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(instruction)
	scanner.Scan()
	return scanner.Text()
}

func validateYear(year string) {
	if len([]rune(year)) < 4 {
		log.Fatalln("Invalid year parameter, must contain 4 digits")
	}
}

func validateMonth(month string) {
	if len([]rune(month)) > 2 {
		log.Fatalln("Invalid month parameter, must contain 1 or 2 digits")
	}
	if month == "0" {
		log.Fatalln("Invalid month parameter, must be between 1 and 12")
	}
}

func convertToNumber(sValue string) int {
	number, err := strconv.Atoi(sValue)
	HandleError(err)

	return number
}
