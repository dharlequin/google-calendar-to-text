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
	AddEmptyLine()
	path := scanner.Text()

	ValidateFilePath(path)
	return path
}

func GetYearFromUser(instruction string) int {
	currentYear := time.Now().Year()

	sYear := getNumberFromUser(fmt.Sprintf(instruction, NumberToString(currentYear)))

	if sYear == "" {
		return currentYear
	}

	ValidateYear(sYear)

	return StringToNumber(sYear)
}

func GetMonthFromUser(instruction string) int {
	currentMonth := time.Now().Month()

	sMonth := getNumberFromUser(fmt.Sprintf(instruction, currentMonth.String()))

	if sMonth == "" {
		return int(currentMonth)
	}

	month := StringToNumber(sMonth)
	ValidateMonth(month)

	return month
}

func getNumberFromUser(instruction string) string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(instruction)
	// scanner.Scan() is used here to allow for empty inputs, which is not possible with fmt.Scan()
	scanner.Scan()
	AddEmptyLine()
	return scanner.Text()
}
