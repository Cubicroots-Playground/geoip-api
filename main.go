package main

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/oschwald/geoip2-golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	slog.Info("reading geo data file")
	geoDB, err := geoip2.Open("geodata.mmdb")
	if err != nil {
		panic(err)
	}

	slog.Info("initializing metrics")
	metricsRequestsByCountry := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Help: "Counts requests by country ISO code.",
			Name: "geoip_requests_country_total",
		},
		[]string{"country"},
	)

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

		var country string
		if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			country = "private"
		} else {
			obj, err := geoDB.Country(ip)
			if err != nil {
				slog.Info("failed to lookup IP: " + err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			country = obj.Country.IsoCode
		}

		metricsRequestsByCountry.WithLabelValues(country).Inc()

		_, _ = w.Write([]byte(country))
	}))
	http.Handle("/metrics", promhttp.Handler())

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
