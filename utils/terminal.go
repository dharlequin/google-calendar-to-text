package utils

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func GetFilePath(instruction string) string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(instruction)
	scanner.Scan()
	path := scanner.Text()

	ValidateFilePath(path)
	return path
}

func GetYearFromUser(instruction string) int {
	sYear := getNumberFromUser(instruction)

	if sYear == "" {
		return time.Now().Year()
	}

	ValidateYear(sYear)

	return ConvertToNumber(sYear)
}

func GetMonthFromUser(instruction string) int {
	sMonth := getNumberFromUser(instruction)

	if sMonth == "" {
		return int(time.Now().Month())
	}

	month := ConvertToNumber(sMonth)
	ValidateMonth(month)

	return month
}

func getNumberFromUser(instruction string) string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(instruction)
	scanner.Scan()
	return scanner.Text()
}
