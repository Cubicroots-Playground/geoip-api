package main

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"

	"github.com/oschwald/geoip2-golang"
)

func main() {
	geoDB, err := geoip2.Open("dbip-country-lite.mmdb")
	if err != nil {
		panic(err)
	}

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawIP, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Info("failed to read body: " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

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

	slog.Info("listening on :8080")
	err = http.ListenAndServe(":8080", nil)
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error(err.Error())
	}
	slog.Info("shutting down, bye")
}
