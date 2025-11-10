package main

import (
	"fmt"

	"github.com/shegai01/project1/internal"
)

func main() {

	urls := []string{
		"http://google.com",
		"http://vk.com",
		"http://yandex.com",
		"http://f.com",
		"http://eeee",
	}
	result := make(chan internal.Result)
	for _, url := range urls {
		go internal.CheckUrls(url, result)
	}
	for i := 0; i < len(urls); i++ {
		res := <-result
		if res.Success {
			fmt.Printf("%s - ok\n", res.Url)
		} else {
			fmt.Printf("%s - not ok\n", res.Url)
		}

	}
}
