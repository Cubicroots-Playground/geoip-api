# GeoIP API

Dead simple API to get geo information for IP addresses. Including prometheus metrics.

## Usage

Download a GeoIP database containing country codes, e.g. from [db-ip.com](https://db-ip.com/db/download/ip-to-country-lite) and add it to the main folder as `geodata.mmdb`. 

Run the webserver with `go run main.go`. 

Call `curl localhost:8080 --data-raw "8.8.8.8"` to get the country code for an IP.

Metrics are served at `localhost:8080/metrics`.

### Docker Compose

```yaml
version: "3.3"

services:
  geoip-api:
    image: cubicrootxyz/geoip-api:beta
    container_name: geoip-api
    volumes:
      - "./geodata.mmdb:/run/geodata.mmdb"
    environment:
      HTTP_PORT: "8080"
```
