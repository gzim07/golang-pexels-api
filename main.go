package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	PhotoApi = "https://api.pexels.com/v1/photos"
	VideoApi = "https://api.pexels.com/v1/videos"
)

type Client struct {
	TOKEN          string
	hc             http.Client
	RemainingTimes int32
}

func NewClient(token string) *Client {
	c := http.Client{}
	return &Client{TOKEN: token, hc: c}
}

type SearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Photos       []Photo `json:"photos"`
}

type Photo struct {
	Id              int32       `json:"id"`
	Width           int32       `json:"width"`
	Height          int32       `json:"height"`
	Url             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	PhotographerUrl string      `json:"photographer_url"`
	Src             PhotoSource `json:"src"`
}

type PhotoSource struct {
	Original  string `json:"original"`
	Large     string `json:"large"`
	Large2x   string `json:"large2x"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Potrait   string `json:"portrait"`
	Square    string `json:"square"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

func (c *Client) SearchPhotos(query string, perPage int, page int) (*SearchResult, error) {
	url := fmt.Sprintf(PhotoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)

	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result SearchResult
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, err
}
func (c *Client) CuratedPhotos(perPage, page int) (*CuratedResult, error) {}
func (c *Client) requestDoWithAuth(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.TOKEN)
	resp, err := c.hc.Do(req)
	if err != nil {
		return resp, err
	}
	times, err := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaing"))
	if err != nil {
		return resp, err
	} else {
		c.RemainingTimes = int32(times)
	}
	return resp, nil
}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	TOKEN := os.Getenv("APIKEY")
	var c = NewClient(TOKEN)

	result, err := c.SearchPhotos("waves", 15, 1)
	if err != nil {
		fmt.Errorf("Search errord:%v", err)
	}
	fmt.Println(result)
}
