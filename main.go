package main

import (
	"github.com/bhrg3se/reverse-proxy/config"
	"github.com/bhrg3se/reverse-proxy/routes"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log.SetLevel(log.InfoLevel)
	m := http.NewServeMux()

	conf := config.GetConfig()
	for _, backend := range conf.Backends {
		routes.AddRoute(backend, m)
	}
	log.Info("Listening on: 0.0.0.0:8000")
	log.Fatal(http.ListenAndServe(":8080", m))

}
