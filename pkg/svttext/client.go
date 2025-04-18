package svttext

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const URL_FMT = "https://www.svt.se/text-tv/api/%s"

type Client struct {
	useCache bool
	cache    map[string]Page
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
	if page, ok := c.cache[id]; ok && c.useCache {
		return page, nil
	}

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

	if c.useCache {
		c.cache[id] = page
	}

	return page, nil
}

func (c *Client) CacheSize() int {
	return len(c.cache)
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
