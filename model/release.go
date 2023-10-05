package model

type Release struct {
	Year   int
	Month  int
	Shows  []ReleaseItem
	Movies []ReleaseItem
	Games  []ReleaseItem
	Other  []ReleaseItem
}