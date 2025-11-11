package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shegai01/project1/internal"
)

func main() {
	r := mux.NewRouter()
	h := internal.NewHandler(r)
	r.HandleFunc("/status", h.Status)
	log.Fatal(http.ListenAndServe(":8080", r))

}

// 	// logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
// 	// 	Level: slog.LevelInfo,
// 	// }))

// 	// appPort := os.Getenv("APP_PORT")
// 	// if appPort == "" {
// 	// 	fmt.Println("env is missing: APP_PORT")
// 	// }

// 	// urls := []string{
// 	// 	"http://google.com",
// 	// 	"http://vk.com",
// 	// 	"http://yandex.com",
// 	// 	"http://f.com",
// 	// 	"http://eeee",
// 	// }
// 	result := make(chan internal.Result)
// 	// wg := &sync.WaitGroup{}
// 	// wg.Add(len(urls))
// 	// for _, url := range urls {
// 	// 	// defer wg.Done()
// 	// 	go internal.CheckUrls(url, result)
// 	// // }
// 	// for i := 0; i < len(urls); i++ {
// 	res := <-result
// 	if res.Success {
// 		fmt.Printf("%s - ok\n", res.Url)
// 	} else {
// 		fmt.Printf("%s - not ok\n", res.Url)
// 	}

// 	// }
// 	r := mux.NewRouter()
// 	internal.NewHandler(r)
// 	// wg.Wait()

// 	// handler := internal.NewHandler()

// 	// server := &http.Server{
// 	// 	Addr:              ":" + appPort,
// 	// 	Handler:           handler,
// 	// 	ReadHeaderTimeout: time.Second,
// 	// 	ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelInfo),
// 	// }

// 	// signalChan := make(chan os.Signal, 1)

// 	// signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

// 	// wg := &sync.WaitGroup{}
// 	// go func() {
// 	// 	if err := server.ListenAndServe(); err != nil {
// 	// 		return
// 	// 	}
// 	// }()
// 	// wg.Add(1)
// }
