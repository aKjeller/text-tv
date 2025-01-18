package svttext

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const URL_FMT = "https://www.svt.se/text-tv/api/%s"

type Page struct {
	Text  string
	Image string
}

func GetPage(id string) (Page, error) {
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

	return Page{Text: response.Data.SubPages[0].AltText, Image: response.Data.SubPages[0].GifAsBase64}, nil
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
