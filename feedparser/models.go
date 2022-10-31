package feedparser

// Video struct for each video in the feed
type Video struct {
	ID        string `xml:"id"`
	Title     string `xml:"title"`
	Published Date   `xml:"published"`
}

// Channel struct for RSS
type Channel struct {
	ID     string
	Name   string  `xml:"title"`
	Videos []Video `xml:"entry"`
}
