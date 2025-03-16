# GeoIP API

Dead simple API to get geo information for IP addresses.

## Usage

Download a GeoIP database containing country codes, e.g. from [db-ip.com](https://db-ip.com/db/download/ip-to-country-lite) and add it to the main folder as `dbip-country-lite.mmdb`. 

Run the webserver with `go run main.go`. 

Call `curl localhost:8080 --data-raw "8.8.8.8"` to get the country code for an IP.
