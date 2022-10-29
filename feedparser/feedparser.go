package feedparser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/lib/log"
	"github.com/paulrosania/go-charset/charset"
)

var (
	ErrInvalidChannelID = errors.New("invalid channel ID")
	ErrParseFailed      = errors.New("failed to parse feed")
)

var urlFormat = "https://www.youtube.com/feeds/videos.xml?channel_id=%s"

// Video struct for each video in the feed
type Video struct {
	ID        string `xml:"id"`
	Title     string `xml:"title"`
	Published Date   `xml:"published"`
}

// Channel struct for RSS
type Channel struct {
	Name   string  `xml:"title"`
	Videos []Video `xml:"entry"`
}

func read(l log.Logger, url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		l.Log("level", "ERROR", "function", "feedparser.read", "call", "http.NewRequest", "error", err)
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Log("level", "ERROR", "function", "feedparser.read", "call", "http.Do", "error", err)
		return nil, err
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrInvalidChannelID
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get feed with status %d", response.StatusCode)
	}

	return response.Body, nil
}

// Parse parses a YouTube channel XML feed from a channel ID
func Parse(l log.Logger, channelID string) (*Channel, error) {
	url := fmt.Sprintf(urlFormat, channelID)
	reader, err := read(l, url)
	if err != nil {
		return nil, err
	}

	defer reader.Close()
	xmlDecoder := xml.NewDecoder(reader)
	xmlDecoder.CharsetReader = charset.NewReader

	var channel Channel
	if err := xmlDecoder.Decode(&channel); err != nil {
		l.Log("level", "ERROR", "function", "feedparser.read", "call", "xml.Decode", "error", err)
		return nil, fmt.Errorf("%w: %s", ErrParseFailed, err.Error())
	}
	return &channel, nil
}
