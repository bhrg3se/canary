package main

import (
	"github.com/bhrg3se/reverse-proxy/config"
	"github.com/bhrg3se/reverse-proxy/routes"
	"log"
	"net/http"
)

func main() {

	m := http.NewServeMux()

	conf := config.GetConfig()
	for _, backend := range conf.Backends {
		routes.AddRoute(backend, m)
	}
	log.Fatal(http.ListenAndServe(":8080", m))

}
