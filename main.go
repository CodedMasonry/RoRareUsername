package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type HTTPResult struct {
	Code int `json:"code"`
}

func getjson(client *retryablehttp.Client, url string, body []byte, target interface{}) error {
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error While Fetching: ", err)
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func incrementString(s string) string {
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] < 'z' {
			runes[i]++
			break
		} else {
			runes[i] = 'a'
			if i == 0 {
				runes = append([]rune{'a'}, runes...)
			}
		}
	}
	return string(runes)
}

func worker(id int, client *retryablehttp.Client, jobs <-chan string, results chan<- string) {
	for username := range jobs {
		var result HTTPResult
		url := "https://auth.roblox.com/v1/usernames/validate"

		body := map[string]interface{}{
			"username": username,
			"context":  "Signup",
			"birthday": "2000-06-29T04:00:00.000Z",
		}
		jsonBody, err := json.Marshal(body)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		err2 := getjson(client, url, jsonBody, &result)
		if err2 != nil {
			panic(err)
		}

		if result.Code != 0 {
			results <- username
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	jobs := make(chan string, 5)
	results := make(chan string, 5)
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.Backoff = retryablehttp.DefaultBackoff

	go func() {
		for result := range results {
			fmt.Printf(" Unique Username Found: %v\n", result)
		}
	}()

	for workers := 0; workers < 3; workers++ {
		go worker(workers, retryClient, jobs, results)
	}

	for s := "aaa"; s != "aaaaaa{"; s = incrementString(s) {
		jobs <- s
	}
	close(jobs)

	fmt.Println("All usernames checked")
}
