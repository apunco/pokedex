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
	val, ok, err := checkCache[RespLocations](c, url)
	if err != nil {
		return RespLocations{}, err
	}
	if ok {
		return val, nil
	}

	locationsResponse, err := getClientCall[RespLocations](c, url)
	if err != nil {
		return RespLocations{}, err
	}

	return locationsResponse, nil
}

func (c *Client) ExploreLocation(location string) (RespLocation, error) {
	url := baseUrl + "/location-area/" + location

	//Return cached response
	val, ok, err := checkCache[RespLocation](c, url)

	if err != nil {
		return RespLocation{}, err
	}
	if ok {
		return val, nil
	}

	locationResponse, err := getClientCall[RespLocation](c, url)
	if err != nil {
		return RespLocation{}, err
	}

	return locationResponse, nil
}

func getClientCall[T any](c *Client, url string) (T, error) {
	var result T

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	//Cache response
	c.cache.Add(url, dat)

	err = json.Unmarshal(dat, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func checkCache[T any](c *Client, url string) (T, bool, error) {
	var genericVal T

	//Return cached response
	if val, ok := c.cache.Get(url); ok {
		err := json.Unmarshal(val, &genericVal)
		if err != nil {
			return genericVal, false, err
		}

		return genericVal, true, nil
	}
	return genericVal, false, nil
}
