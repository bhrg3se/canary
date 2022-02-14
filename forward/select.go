package forward

import (
	"github.com/bhrg3se/canary/config"
	"github.com/bhrg3se/canary/iprouter"
	log "github.com/sirupsen/logrus"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/roundrobin"
	"github.com/vulcand/oxy/testutils"
	"net/http"
)

func ToCanary(config *config.Config) http.Handler {
	switch config.CanaryConfig.RouteType {
	case "traffic":
		// Forwards incoming requests to whatever location URL points to, adds proper forwarding headers
		fwd, _ := forward.New()
		lb, _ := roundrobin.New(fwd)

		lb.UpsertServer(testutils.ParseURI(config.CanaryUpstreamHost), roundrobin.Weight(config.CanaryConfig.TrafficPercent))
		lb.UpsertServer(testutils.ParseURI(config.MainUpstreamHost), roundrobin.Weight(100-config.CanaryConfig.TrafficPercent))

		return lb

	case "content":
		fwd, _ := forward.New()
		ipRouter := iprouter.NewIPRouter(fwd, config)
		return ipRouter

	default:
		log.Fatal("invalid router type")
	}
	panic("invalid router type")
}

func ToMain(host string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// let us forward this request to another server
		fwd, _ := forward.New()
		req.URL = testutils.ParseURI(host)
		fwd.ServeHTTP(w, req)
	})
}
