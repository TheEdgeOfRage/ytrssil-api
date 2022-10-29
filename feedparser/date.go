package feedparser

import (
	"time"
)

// Date type
type Date string

// Parse (Date function) and returns Time, error
func (d Date) Parse() (time.Time, error) {
	return time.Parse(time.RFC3339, string(d)) // ISO8601
}
