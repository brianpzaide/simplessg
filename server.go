package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const OutputDir = "./dist"
const OutputDirPrefix = "/dist/"

func genRSSFeedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		rssStruct := RSSChannel{
			Title:         Title,
			Description:   Description,
			LastBuildDate: time.Now().Format("Fri, 30 Jan 2026"),
			AtomLink:      AtomLink{},
			Items:         []RSSItem{},
		}

		rssFeed, err := xml.Marshal(rssStruct)
		if err != nil {
			fmt.Println("could not marshal to xml encoding", err)
			return
		}

		w.Header().Set("Content-Type", "application/xml")
		w.Write(rssFeed)
		return
	}

	env := envelope{"error": fmt.Sprintf("the %s method is not supported for this resource", r.Method)}
	err := writeJSON(w, http.StatusMethodNotAllowed, env, nil)
	if err != nil {
		w.WriteHeader(500)
	}
}

func writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func serve() error {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(OutputDir))
	mux.Handle(OutputDirPrefix, http.StripPrefix(OutputDirPrefix, fs))
	mux.HandleFunc("/rss.xml", genRSSFeedHandler)

	srv := &http.Server{
		Addr:    ":4000",
		Handler: mux,
	}
	shutdownErr := make(chan error, 1)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		log.Printf("caught signal %s", s.String())
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdownErr <- srv.Shutdown(ctx)
	}()

	log.Printf("starting server addr: %s\n", srv.Addr)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return <-shutdownErr
}
