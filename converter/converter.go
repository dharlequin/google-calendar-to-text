package converter

import (
	"bufio"
	"dharlequin/google-calendar-converter/utils"
	"fmt"
	"os"
	"strings"
	"time"
)

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

const TIME_LAYOUT = "20060102"

const RECORD_START = "BEGIN:VEVENT"
const RECORD_TIME = "DTSTART;VALUE=DATE:"
const RECORD_TITLE = "SUMMARY:"
const RECORD_CATEGORY = "DESCRIPTION:"
const RECORD_END = "END:VEVENT"

func ConvertFileToText(filePath string, year int, month int) {
	source, err := os.Open(filePath)
	utils.HandleError(err)

	defer source.Close()

	var shows []ReleaseItem
	var movies []ReleaseItem
	var games []ReleaseItem
	var other []ReleaseItem

	var item ReleaseItem

	// this is the key to look for during the next iteration
	key := RECORD_START

	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		var line = scanner.Text()
		fmt.Println(line)

		// new calendar record starts with this
		if key == RECORD_START && line == RECORD_START {
			fmt.Println("--New record found")
			item = ReleaseItem{}
			key = RECORD_TIME
			continue
		}

		// this is date and time of calendar record
		if key == RECORD_TIME && strings.HasPrefix(line, RECORD_TIME) {
			fmt.Println("--Found record date")
			var datePart = utils.ExtractStringValue(line, RECORD_TIME)
			var recordTime, err = time.Parse(TIME_LAYOUT, datePart)
			utils.HandleError(err)

			// this is the one we are looking for
			if recordTime.Year() == year && int(recordTime.Month()) == month {
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
		if key == RECORD_CATEGORY && strings.HasPrefix(line, RECORD_CATEGORY) {
			fmt.Println("--Found record category")
			var category = utils.ExtractStringValue(line, RECORD_CATEGORY)

			item.Category = category

			key = RECORD_TITLE
			continue
		}

		// this is the record title
		if key == RECORD_TITLE && strings.HasPrefix(line, RECORD_TITLE) {
			fmt.Println("--Found record summary")
			summary := utils.ExtractStringValue(line, RECORD_TITLE)

			sumParts := strings.Split(summary, " - ")
			item.Title = sumParts[0]

			if len(sumParts) > 1 {
				item.Comments = sumParts[1]
			}

			switch strings.ToLower(item.Category) {
			case SHOW:
				shows = append(shows, item)
			case MOVIE:
				movies = append(movies, item)
			case GAME:
				games = append(games, item)
			default:
				other = append(other, item)
			}

			key = RECORD_START
			continue
		}

		// end of file
		if line == "END:VCALENDAR" {
			fmt.Println("-- Reached end of file")
		}
	}

	// sorting does not seem to be required, as it is done by Google Calendar internally
	fmt.Println("Shows", shows)
	fmt.Println("Movies", movies)
	fmt.Println("Games", games)
	fmt.Println("Other", other)

	fileName := fmt.Sprintf("releases-%d-%d.txt", month, year)
	target, err := os.Create(fileName)
	utils.HandleError(err)
	defer target.Close()

	w := bufio.NewWriter(target)
	w.WriteString(getFileName(month))

	printCategory(movies, "ФИЛЬМЫ", w)
	printCategory(shows, "СЕРИАЛЫ", w)
	printCategory(games, "ИГРЫ", w)
	printCategory(other, "ДРУГОЕ", w)

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
			item = fmt.Sprintf("%s - **%s** - %s\n", i.Date.Format("02/01"), i.Title, i.Comments)
		} else {
			item = fmt.Sprintf("%s - **%s**\n", i.Date.Format("02/01"), i.Title)
		}

		w.WriteString(item)
	}

	w.WriteString("\n")
}
