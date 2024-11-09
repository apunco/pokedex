package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	pokecache "github.com/apunco/go/pokedex/internal/cache"
)

type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(time.Duration(5 * time.Second)),
	}
}

func (c *Client) GetLocations(pageUrl *string) (RespLocations, error) {
	url := baseUrl + "/location-area"

	if pageUrl != nil {
		url = *pageUrl
	}

	//Return cached response
	if val, ok := c.cache.Get(url); ok {
		var locations RespLocations
		err := json.Unmarshal(val, &locations)
		if err != nil {
			return RespLocations{}, err
		}

		return locations, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespLocations{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RespLocations{}, err
	}

	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RespLocations{}, err
	}

	locationResp := RespLocations{}
	err = json.Unmarshal(dat, &locationResp)
	if err != nil {
		return RespLocations{}, err
	}

	//Cache response
	c.cache.Add(url, dat)
	return locationResp, nil
}
