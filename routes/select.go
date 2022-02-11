package routes

import (
	"github.com/bhrg3se/reverse-proxy/config"
	"github.com/jpillora/ipfilter"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/roundrobin"
	"github.com/vulcand/oxy/testutils"
	"net/http"
)

func AddRoute(backend config.Backend, m *http.ServeMux) {

	f := ipfilter.New(ipfilter.Options{
		BlockByDefault: true,
	})

	// initialize a reverse proxy and pass the actual backend server url here

	// Forwards incoming requests to whatever location URL points to, adds proper forwarding headers
	fwd, _ := forward.New()

	lb, _ := roundrobin.New(fwd)

	for _, upstream := range backend.Upstreams {
		if upstream.Routing.IPWhitelist != "" {
			f.AllowIP(upstream.Routing.IPWhitelist)
		}
		if upstream.Routing.IPBlacklist != "" {
			f.BlockIP(upstream.Routing.IPWhitelist)
		}

		lb.UpsertServer(testutils.ParseURI(upstream.Host), roundrobin.Weight(upstream.Weight))
	}

	// handle all requests to your server using the proxy

	m.HandleFunc(backend.UrlPattern, lb.ServeHTTP)

}
