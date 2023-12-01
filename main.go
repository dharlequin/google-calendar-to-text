package main

import (
	"dharlequin/google-calendar-converter/converter"
	"dharlequin/google-calendar-converter/utils"
	"fmt"
)

var Version = "v1.1.2"

func main() {
	fmt.Println("Version:\t", Version)

	var filePath string
	var year int
	var month int

	filePath = utils.GetFilePath("Enter path to Google Calendar file (.ics):")
	year = utils.GetYearFromUser("Specify year (as 2023) to look for, leave empty to use current - %s:")
	month = utils.GetMonthFromUser("Specify month (from 1 to 12) to look for, leave empty to use current - %s:")

	release := converter.GetParsedReleases(filePath, year, month)

	converter.ExportToFile(&release)
}
