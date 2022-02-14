package main

import (
	"github.com/bhrg3se/reverse-proxy/config"
	"github.com/bhrg3se/reverse-proxy/forward"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log.SetLevel(log.InfoLevel)
	m := http.NewServeMux()

	conf := config.GetConfig()

	m.Handle(conf.MainUpstreamPattern, forward.ToMain(conf.MainUpstreamHost))
	m.Handle(conf.CanaryUpstreamPattern, forward.ToCanary(conf))

	log.Info("Listening on: 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", m))

}
