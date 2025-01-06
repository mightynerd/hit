package discogs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/mightynerd/hit/db"
)

type DiscogsConfig struct {
	apiKey string
}

type DiscogsSearchResult struct {
	Country  string   `json:"country"`
	Year     string   `json:"year,omitempty"`
	Format   []string `json:"format"`
	Label    []string `json:"label"`
	Type     string   `json:"type"`
	Genre    []string `json:"genre"`
	Style    []string `json:"style"`
	ID       int      `json:"id"`
	Barcode  []any    `json:"barcode"`
	UserData struct {
		InWantlist   bool `json:"in_wantlist"`
		InCollection bool `json:"in_collection"`
	} `json:"user_data"`
	MasterID    int    `json:"master_id"`
	MasterURL   any    `json:"master_url"`
	URI         string `json:"uri"`
	Catno       string `json:"catno"`
	Title       string `json:"title"`
	Thumb       string `json:"thumb"`
	CoverImage  string `json:"cover_image"`
	ResourceURL string `json:"resource_url"`
	Community   struct {
		Want int `json:"want"`
		Have int `json:"have"`
	} `json:"community"`
	FormatQuantity int `json:"format_quantity,omitempty"`
	Formats        []struct {
		Name         string   `json:"name"`
		Qty          string   `json:"qty"`
		Descriptions []string `json:"descriptions"`
	} `json:"formats,omitempty"`
}

type DiscogsSearchResults struct {
	Pagination struct {
		Page    int `json:"page"`
		Pages   int `json:"pages"`
		PerPage int `json:"per_page"`
		Items   int `json:"items"`
		Urls    struct {
		} `json:"urls"`
	} `json:"pagination"`
	Results []DiscogsSearchResult `json:"results"`
}

func NewDiscogsConfig(apiKey string) *DiscogsConfig {
	config := &DiscogsConfig{
		apiKey: apiKey,
	}

	return config
}

func (config *DiscogsConfig) GetEarliestReleaseYear(artist string, track string) (int, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.discogs.com/database/search?type=release&artist=%s&track=%s", artist, track), nil)

	if err != nil {
		return 0, err
	}

	request.Header.Add("User-Agent", "github.com/mightynerd")
	request.Header.Add("Authorization", fmt.Sprintf("Discogs token=%s", config.apiKey))

	resp, err := client.Do(request)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, err
	}

	var result DiscogsSearchResults
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	var years []int
	for _, res := range result.Results {
		year, err := strconv.Atoi(res.Year)
		if err == nil {
			years = append(years, year)
		}

	}

	if len(years) < 1 {
		return 0, fmt.Errorf("no results found for artist '%s', track '%s'", artist, track)
	}

	slices.Sort(years)

	return years[0], nil
}

func (config *DiscogsConfig) EnhanceYears(tracks *([]db.Track)) {
	for i, track := range *tracks {
		dcYear, err := config.GetEarliestReleaseYear(track.Artist, track.Title)
		time.Sleep(1 * time.Second)
		if err != nil {
			fmt.Println("Enhance error", err)
			continue
		}

		fmt.Printf("Track '%s', artist '%s', org year %d, dc year %d\n", track.Title, track.Artist, track.Year, dcYear)

		if track.Year > dcYear {
			fmt.Println("DC year is better for above")
			(*tracks)[i].Year = dcYear
		}

	}
}
