package internal

import "net/http"

type Result struct {
	Url     string
	Success bool
}

func CheckUrls(url string, res chan<- Result) {
	response, err := http.Get(url)
	if err != nil {
		res <- Result{Url: url, Success: false}
		return
	}
	success := response.StatusCode == 200
	response.Body.Close()
	res <- Result{Url: url, Success: success}
}
