package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	Id              int32
	Width           int32
	Height          int32
	Url             string
	Photographer    string
	PhotographerUrl string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	TOKEN := os.Getenv("APIKEY")
	var c = NewClient(TOKEN)

	result, err := c.SearchPhotos("waves")
	if err != nil {
		fmt.Errorf("Search errord:%v", err)
	}
	fmt.Println(result)
}
