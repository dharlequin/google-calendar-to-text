package converter

import (
	"bufio"
	"dharlequin/google-calendar-converter/model"
	"dharlequin/google-calendar-converter/utils"
	"fmt"
	"os"
	"strings"
	"time"
)

const SHOW = "show"
const MOVIE = "movie"
const GAME = "game"
const OTHER = "other"

const CALENDAR_TIME_LAYOUT = "20060102"

const RECORD_START = "BEGIN:VEVENT"
const RECORD_DATE = "DTSTART;VALUE=DATE:"
const RECORD_TITLE = "SUMMARY:"
const RECORD_CATEGORY = "DESCRIPTION:"
const RECORD_END = "END:VEVENT"

const EVENT_TITLE_DELIMITER = " - "

func GetParsedReleases(filePath string, year int, month int) model.Release {
	var release = model.Release{
		Year:  year,
		Month: month,
	}

	parseCalendar(filePath, &release)

	utils.AddEmptyLine()
	fmt.Println("Shows:", release.Shows)
	utils.AddEmptyLine()
	fmt.Println("Movies:", release.Movies)
	utils.AddEmptyLine()
	fmt.Println("Games:", release.Games)
	utils.AddEmptyLine()
	fmt.Println("Other:", release.Other)
	utils.AddEmptyLine()

	utils.SortCategories(&release)

	return release
}

func parseCalendar(filePath string, release *model.Release) {
	source, err := os.Open(filePath)
	utils.HandleError(err)

	defer source.Close()

	var item model.ReleaseItem

	// this is the key to look for during the next iteration
	key := RECORD_START

	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		var line = scanner.Text()

		if !strings.HasPrefix(line, key) {
			continue
		}

		// new calendar record starts with this
		if line == RECORD_START {
			utils.AddEmptyLine()
			utils.PrintParserLog("New record found:", line, 1)

			item = model.ReleaseItem{}
			key = RECORD_DATE
			continue
		}

		// this is date and time of calendar record
		if strings.HasPrefix(line, RECORD_DATE) {
			utils.PrintParserLog("Found record date:", line, 2)

			var datePart = utils.ExtractStringValue(line, RECORD_DATE)
			var recordTime, err = time.Parse(CALENDAR_TIME_LAYOUT, datePart)
			utils.HandleError(err)

			// this is the one we are looking for
			if recordTime.Year() == release.Year && int(recordTime.Month()) == release.Month {
				utils.PrintParserLog("Within target parameters:", "", 3)

				item.Date = recordTime

				key = RECORD_CATEGORY
				continue
			} else {
				utils.PrintParserLog("Outside of target parameters", "", 3)

				key = RECORD_START
				continue
			}
		}

		// this is the record category
		if strings.HasPrefix(line, RECORD_CATEGORY) {
			utils.PrintParserLog("Found record category:", line, 2)

			var category = utils.ExtractStringValue(line, RECORD_CATEGORY)

			item.Category = category

			key = RECORD_TITLE
			continue
		}

		// this is the record title
		if strings.HasPrefix(line, RECORD_TITLE) {
			utils.PrintParserLog("Found record summary:", line, 2)

			summary := utils.ExtractStringValue(line, RECORD_TITLE)
			summary = utils.NormalizeString(summary)

			sumParts := strings.Split(summary, EVENT_TITLE_DELIMITER)
			item.Title = sumParts[0]

			if len(sumParts) > 1 {
				for i := 1; i < len(sumParts)-1; i++ {
					item.Title += EVENT_TITLE_DELIMITER + sumParts[i]
				}

				item.Comments = sumParts[len(sumParts)-1]
			}

			switch strings.ToLower(item.Category) {
			case SHOW:
				release.Shows = append(release.Shows, item)
			case MOVIE:
				release.Movies = append(release.Movies, item)
			case GAME:
				release.Games = append(release.Games, item)
			default:
				release.Other = append(release.Other, item)
			}

			key = RECORD_START
			continue
		}
	}
}
