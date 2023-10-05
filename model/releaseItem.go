package model

import "time"

type ReleaseItem struct {
	Title    string
	Comments string
	Category string
	Date     time.Time
}
