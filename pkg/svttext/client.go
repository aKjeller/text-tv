package svttext

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const URL_FMT = "https://www.svt.se/text-tv/api/%s"

type Client struct {
	useCache bool
	cache    map[string]Page
	mu       sync.RWMutex
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		cache: make(map[string]Page),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type Option func(*Client)

func WithCaching() Option {
	return func(c *Client) {
		c.useCache = true
	}
}

func (c *Client) GetPage(id string) (Page, error) {
	if c.useCache {
		c.mu.RLock()
		page, ok := c.cache[id]
		c.mu.RUnlock()
		if ok {
			return page, nil
		}
	}

	// TODO: avoid duplicate fetches of same id
	page, err := getPageFromTextTv(id)
	if err != nil {
		return Page{}, fmt.Errorf("failed to get page from text-tv: %w", err)
	}

	if c.useCache {
		c.mu.Lock()
		c.cache[id] = page
		c.mu.Unlock()
	}

	return page, nil
}

func (c *Client) CacheSize() int {
	return len(c.cache)
}

func getPageFromTextTv(id string) (Page, error) {
	url := fmt.Sprintf(URL_FMT, id)
	resp, err := http.Get(url)
	if err != nil {
		return Page{}, fmt.Errorf("error during GET request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Page{}, fmt.Errorf("error reading response body: %w", err)
	}

	var response response
	if err := json.Unmarshal(body, &response); err != nil {
		return Page{}, fmt.Errorf("error parsing JSON response: %w", err)
	}

	page, err := newPage(response.Data)
	if err != nil {
		return page, fmt.Errorf("failed to create new page, %w", err)
	}
	return page, nil
}

type response struct {
	Status string `json:"status"`
	Data   data   `json:"data"`
}

type data struct {
	PageNumber string    `json:"pageNumber"`
	PrevPage   string    `json:"prevPage"`
	NextPage   string    `json:"nextPage"`
	SubPages   []subPage `json:"subPages"`
	Meta       meta      `json:"meta"`
}

type subPage struct {
	SubPageNumber string `json:"subPageNumber"`
	GifAsBase64   string `json:"gifAsBase64"`
	ImageMap      string `json:"imageMap"`
	AltText       string `json:"altText"`
}

type meta struct {
	Updated string `json:"updated"`
}
