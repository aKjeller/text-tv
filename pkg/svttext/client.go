package svttext

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const URL_FMT = "https://www.svt.se/text-tv/api/%s"

type Page struct {
	PageNumber string    `json:"pageNumber"`
	PrevPage   string    `json:"prevPage"`
	NextPage   string    `json:"nextPage"`
	SubPages   []SubPage `json:"subPages"`
}

type SubPage struct {
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

	var subPages []SubPage
	for _, page := range response.Data.SubPages {
		subPages = append(subPages, SubPage{Text: page.AltText, Image: page.GifAsBase64})
	}

	return Page{
		PageNumber: response.Data.PageNumber,
		PrevPage:   response.Data.PrevPage,
		NextPage:   response.Data.NextPage,
		SubPages:   subPages,
	}, nil
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
