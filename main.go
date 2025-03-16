package main

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/oschwald/geoip2-golang"
)

func main() {
	slog.Info("reading geo data file")
	geoDB, err := geoip2.Open("geodata.mmdb")
	if err != nil {
		panic(err)
	}

	slog.Info("setting up webserver")
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawIP, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Info("failed to read body: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		ip := net.ParseIP(string(rawIP))
		if ip == nil {
			slog.Info("failed to parse IP")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		country, err := geoDB.Country(ip)
		if err != nil {
			slog.Info("failed to lookup IP: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, _ = w.Write([]byte(country.Country.IsoCode))
	}))

	listenPort := os.Getenv("HTTP_PORT")
	if listenPort == "" {
		listenPort = "8080"
	}

	slog.Info("listening on :" + listenPort)
	err = http.ListenAndServe(":"+listenPort, nil)
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error(err.Error())
	}
	slog.Info("shutting down, bye")
}
