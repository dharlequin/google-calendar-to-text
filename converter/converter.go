package converter

import (
	"bufio"
	"dharlequin/google-calendar-converter/utils"
	"fmt"
	"os"
	"strings"
	"time"
)

type Release struct {
	Year   int
	Month  int
	Shows  []ReleaseItem
	Movies []ReleaseItem
	Games  []ReleaseItem
	Other  []ReleaseItem
}

type ReleaseItem struct {
	Title    string
	Comments string
	Category string
	Date     time.Time
}

const SHOW = "show"
const MOVIE = "movie"
const GAME = "game"
const OTHER = "other"

const CALENDAR_TIME_LAYOUT = "20060102"
const FILE_TIME_LAYOUT = "02/01"

const RECORD_START = "BEGIN:VEVENT"
const RECORD_DATE = "DTSTART;VALUE=DATE:"
const RECORD_TITLE = "SUMMARY:"
const RECORD_CATEGORY = "DESCRIPTION:"
const RECORD_END = "END:VEVENT"

func ConvertFileToText(filePath string, year int, month int) {
	var release = Release{
		Year:  year,
		Month: month,
	}

	parseCalendar(filePath, &release)

	// sorting does not seem to be required, as it is done by Google Calendar when adding events
	fmt.Println("Shows", release.Shows)
	fmt.Println("Movies", release.Movies)
	fmt.Println("Games", release.Games)
	fmt.Println("Other", release.Other)

	exportToFile(&release)
}

func parseCalendar(filePath string, release *Release) {
	source, err := os.Open(filePath)
	utils.HandleError(err)

	defer source.Close()

	var item ReleaseItem

	// this is the key to look for during the next iteration
	key := RECORD_START

	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		var line = scanner.Text()
		fmt.Println(line)

		if !strings.HasPrefix(line, key) {
			continue
		}

		// new calendar record starts with this
		if line == RECORD_START {
			fmt.Println("--New record found")
			item = ReleaseItem{}
			key = RECORD_DATE
			continue
		}

		// this is date and time of calendar record
		if strings.HasPrefix(line, RECORD_DATE) {
			fmt.Println("--Found record date")
			var datePart = utils.ExtractStringValue(line, RECORD_DATE)
			var recordTime, err = time.Parse(CALENDAR_TIME_LAYOUT, datePart)
			utils.HandleError(err)

			// this is the one we are looking for
			if recordTime.Year() == release.Year && int(recordTime.Month()) == release.Month {
				fmt.Println("---Within target parameters")
				item.Date = recordTime

				key = RECORD_CATEGORY
				continue
			} else {
				key = RECORD_START
				continue
			}
		}

		// this is the record category
		if strings.HasPrefix(line, RECORD_CATEGORY) {
			fmt.Println("--Found record category")
			var category = utils.ExtractStringValue(line, RECORD_CATEGORY)

			item.Category = category

			key = RECORD_TITLE
			continue
		}

		// this is the record title
		if strings.HasPrefix(line, RECORD_TITLE) {
			fmt.Println("--Found record summary")
			summary := utils.ExtractStringValue(line, RECORD_TITLE)

			sumParts := strings.Split(summary, " - ")
			item.Title = sumParts[0]

			if len(sumParts) > 1 {
				item.Comments = sumParts[1]
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

func exportToFile(release *Release) {
	fileName := fmt.Sprintf("releases-%d-%d.txt", release.Month, release.Year)
	target, err := os.Create(fileName)
	utils.HandleError(err)

	defer target.Close()

	w := bufio.NewWriter(target)
	w.WriteString(getFileName(release.Month))

	printCategory(release.Movies, "ФИЛЬМЫ", w)
	printCategory(release.Shows, "СЕРИАЛЫ", w)
	printCategory(release.Games, "ИГРЫ", w)
	printCategory(release.Other, "ДРУГОЕ", w)

	w.Flush()
}

func getFileName(month int) string {
	var sMonth string
	switch month {
	case 1:
		sMonth = "января"
	case 2:
		sMonth = "февраля"
	case 3:
		sMonth = "марта"
	case 4:
		sMonth = "апреля"
	case 5:
		sMonth = "мая"
	case 6:
		sMonth = "июня"
	case 7:
		sMonth = "июля"
	case 8:
		sMonth = "августа"
	case 9:
		sMonth = "сентября"
	case 10:
		sMonth = "октября"
	case 11:
		sMonth = "ноября"
	case 12:
		sMonth = "января"
	}

	return strings.ToUpper(fmt.Sprintf("**релизы %s**\n\n", sMonth))
}

func printCategory(items []ReleaseItem, category string, w *bufio.Writer) {
	if len(items) > 0 {
		w.WriteString(fmt.Sprintf("**%s:**\n\n", category))
		printReleaseItems(items, w)
	}
}

func printReleaseItems(items []ReleaseItem, w *bufio.Writer) {
	for _, i := range items {
		var item string
		if i.Comments != "" {
			item = fmt.Sprintf("%s - **%s** - %s\n", i.Date.Format(FILE_TIME_LAYOUT), i.Title, i.Comments)
		} else {
			item = fmt.Sprintf("%s - **%s**\n", i.Date.Format(FILE_TIME_LAYOUT), i.Title)
		}

		w.WriteString(item)
	}

	w.WriteString("\n")
}
