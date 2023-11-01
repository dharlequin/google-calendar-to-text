package main

import (
	"dharlequin/google-calendar-converter/converter"
	"dharlequin/google-calendar-converter/utils"
)

func main() {
	var filePath string
	var year int
	var month int

	filePath = utils.GetFilePath("Enter path to Google Calendar file (.ics)")
	year = utils.GetYearFromUser("Specify year (as 2023) to look for, leave empty to user current")
	month = utils.GetMonthFromUser("Specify month (as 1 or 10) to look for, leave empty to use current")

	release := converter.GetParsedReleases(filePath, year, month)

	converter.ExportToFile(&release)
}
