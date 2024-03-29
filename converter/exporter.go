package converter

import (
	"bufio"
	"dharlequin/google-calendar-converter/model"
	"dharlequin/google-calendar-converter/utils"
	"fmt"
	"os"
	"strings"
)

const FILE_TIME_LAYOUT = "02/01"

func ExportToFile(release *model.Release) {
	fileName := fmt.Sprintf("releases-%d-%d.txt", release.Month, release.Year)
	target, err := os.Create(fileName)
	utils.HandleError(err)

	defer target.Close()

	w := bufio.NewWriter(target)

	w.WriteString(getTitleName(release.Month))

	writeCategory(release.Movies, "ФИЛЬМЫ", w)
	writeCategory(release.Shows, "СЕРИАЛЫ", w)
	writeCategory(release.Games, "ИГРЫ", w)
	writeCategory(release.Other, "ДРУГОЕ", w)

	w.Flush()

	fmt.Println("Created file:", fileName)
}

func getTitleName(month int) string {
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
		sMonth = "декабря"
	}

	return strings.ToUpper(fmt.Sprintf("**релизы %s**\n\n", sMonth))
}

func writeCategory(items []model.ReleaseItem, category string, w *bufio.Writer) {
	if len(items) > 0 {
		w.WriteString(fmt.Sprintf("**%s:**\n\n", category))
		writeReleaseItems(items, w)
	}
}

func writeReleaseItems(items []model.ReleaseItem, w *bufio.Writer) {
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
