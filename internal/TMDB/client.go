package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

type SearchResponse struct {
	Page         int            `json:"page"`
	Results      []SearchResult `json:"results"`
	TotalPages   int            `json:"total_pages"`
	TotalResults int            `json:"total_results"`
}

type SearchResult struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Overview    string `json:"overview"`
	ReleaseDate string `json:"release_date"`
	PosterPath  string `json:"poster_path"`
}

func (c *Client) SearchMovie(title string) (SearchResult, error) {
	baseURL := "https://api.themoviedb.org/3/search/movie"

	params := url.Values{}
	params.Set("api_key", c.apiKey)
	params.Set("query", title)

	fullURL := baseURL + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		return SearchResult{}, fmt.Errorf("error calling TMDB search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SearchResult{}, fmt.Errorf("TMDB search returned status %d", resp.StatusCode)
	}

	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return SearchResult{}, fmt.Errorf("error decoding TMDB response: %w", err)
	}

	if len(result.Results) == 0 {
		return SearchResult{}, fmt.Errorf("no TMDB results found for %q", title)
	}

	return result.Results[0], nil
}
