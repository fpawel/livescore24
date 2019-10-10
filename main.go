package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"github.com/fpawel/livescore24/internal/livescore"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		write(w, []byte("Hello, world!"))
	})
	ctx := context.Background()
	workers := livescore.NewWorkers(ctx)

	http.Handle("/champs/", http.StripPrefix("/champs/",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			url := "https://24score.pro/" + r.URL.Path
			games,err := workers.Get(url).Champs()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			httpJSON(w, r, games)
	})))

	server := &http.Server{Addr: ":" + port}
	log.Info("https://localhost:" + port)

	go handleSigint(server.Shutdown)

	log.ErrIfFail(server.ListenAndServe)

	workers.Close()

	log.Info("all closed and canceled")

}

func httpJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	for _, h := range strings.Split(r.Header.Get("Accept-Encoding"), ","){
		if strings.TrimSpace(h) == "gzip" {
			httpGzipJson(w, data)
			return
		}
	}
	b,err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	write(w, b)
}

func httpGzipJson(w http.ResponseWriter, data interface{}){
	gz, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)
	if err != nil {
		panic(err)
	}
	defer log.ErrIfFail(gz.Close)
	w.Header().Set("Content-Encoding", "gzip")

	encoder := json.NewEncoder(gz)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(data); err != nil {
		log.Fatal(err)
	}
}
func write(w io.Writer, b []byte){
	log.ErrIfFail(func() error {
		_,err := w.Write(b)
		return err
	})
}

func handleSigint(closeFunc func(ctx context.Context) error ){
	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Info("SIGINT (pkill -2) accepted")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.ErrIfFail(func() error {
		return closeFunc(ctx)
	})
}