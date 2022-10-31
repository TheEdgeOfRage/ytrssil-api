package feedparser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/lib/log"
	"github.com/paulrosania/go-charset/charset"
)

var (
	ErrInvalidChannelID = errors.New("invalid channel ID")
	ErrParseFailed      = errors.New("failed to parse feed")
)

var urlFormat = "https://www.youtube.com/feeds/videos.xml?channel_id=%s"

type Parser interface {
	Parse(channelID string) (*Channel, error)
	ParseThreadSafe(channelID string, channelChan chan *Channel, errChan chan error, mu *sync.Mutex, wg *sync.WaitGroup)
}

type parser struct {
	log log.Logger
}

func NewParser(l log.Logger) *parser {
	return &parser{
		log: l,
	}
}

func (p *parser) fetch(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		p.log.Log("level", "ERROR", "function", "feedparser.fetch", "call", "http.NewRequest", "error", err)
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		p.log.Log("level", "ERROR", "function", "feedparser.fetch", "call", "http.Do", "error", err)
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
func (p *parser) Parse(channelID string) (*Channel, error) {
	url := fmt.Sprintf(urlFormat, channelID)
	reader, err := p.fetch(url)
	if err != nil {
		return nil, err
	}

	defer reader.Close()
	xmlDecoder := xml.NewDecoder(reader)
	xmlDecoder.CharsetReader = charset.NewReader

	var channel Channel
	if err := xmlDecoder.Decode(&channel); err != nil {
		p.log.Log("level", "ERROR", "function", "feedparser.Parse", "call", "xml.Decode", "error", err)
		return nil, fmt.Errorf("%w: %s", ErrParseFailed, err.Error())
	}
	channel.ID = channelID
	return &channel, nil
}

// ParseThreadSafe calls Parse, but additionally accepts an out parameter to store the result,
// as well as a mutex and wait group to run multiple fetches in parallel
func (p *parser) ParseThreadSafe(
	channelID string, channelChan chan *Channel, errChan chan error, mu *sync.Mutex, wg *sync.WaitGroup,
) {
	channel, err := p.Parse(channelID)

	mu.Lock()
	channelChan <- channel
	errChan <- err
	mu.Unlock()
	wg.Done()
}
