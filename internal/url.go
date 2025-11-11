package internal

import (
	"net/http"
	"time"
)

func CheckUrl(url string) Result {
	client := http.Client{Timeout: 5 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		return Result{
			URL:     url,
			Success: false,
			Error:   err.Error(),
		}
	}
	defer res.Body.Close()
	return Result{
		URL:     url,
		Success: res.StatusCode == http.StatusOK,
	}

}
