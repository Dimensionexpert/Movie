package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

// First call to get ID NAME and RELEASE DATE
type SearchResponse struct {
	Page         int           `json:"page"`
	Results      []SearchMatch `json:"results"`
	TotalPages   int           `json:"total_pages"`
	TotalResults int           `json:"total_results"`
}

type SearchMatch struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
}

// Second call to get other informations
type MovieDetails struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	ReleaseDate string  `json:"release_date"`
	Runtime     int     `json:"runtime"`
	PosterPath  string  `json:"poster_path"` // relative path only, e.g. "/xlaY2....jpg" — prepend https://image.tmdb.org/t/p/{size}/ to get a loadable image
	Genres      []Genre `json:"genres"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Client) SearchMovie(title string) (SearchMatch, error) {
	baseURL := "https://api.themoviedb.org/3/search/movie"

	params := url.Values{}
	params.Set("api_key", c.apiKey)
	params.Set("query", title)

	fullURL := baseURL + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		return SearchMatch{}, fmt.Errorf("error calling TMDB search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SearchMatch{}, fmt.Errorf("TMDB search returned status %d", resp.StatusCode)
	}

	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return SearchMatch{}, fmt.Errorf("error decoding TMDB response: %w", err)
	}

	if len(result.Results) == 0 {
		return SearchMatch{}, fmt.Errorf("no TMDB results found for %q", title)
	}

	return result.Results[0], nil
}

func (c *Client) GetMovieDetails(id int) (MovieDetails, error) {
	idStr := strconv.Itoa(id)
	baseURL := "https://api.themoviedb.org/3/movie/"
	params := url.Values{}
	params.Set("api_key", c.apiKey)
	fullURL := baseURL + idStr + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		return MovieDetails{}, fmt.Errorf("error calling TMDB details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return MovieDetails{}, fmt.Errorf("TMDB details returned status %d", resp.StatusCode)
	}

	var result MovieDetails
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return MovieDetails{}, fmt.Errorf("error decoding TMDB response: %w", err)
	}

	return result, nil
}
