package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type SearchResult struct {
	Name                  string  `json:"Name"`
}

func worker(id int, jobs <-chan string, results chan <- SearchResult) {
	resp, err := retryablehttp.Get("https://www.roblox.com/search/users/results?keyword=abc&maxRows=2")
	if err != nil {
		fmt.Println("Error While Fetching: ", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error While Reading Response Body: ", err)
	}

	var result SearchResult
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		fmt.Println("Error Getting Json: ", err2)
	}

	fmt.Println(result)
}

// Main URL https://www.roblox.com/search/users/results?keyword=%v&maxRows=2
func main() {
	jobs := make(chan string, 5)
	results := make(chan SearchResult, 5)

	go worker(1, jobs, results)

	time.Sleep(2 * time.Second)
}