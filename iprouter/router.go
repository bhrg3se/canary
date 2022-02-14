package iprouter

import (
	"github.com/bhrg3se/reverse-proxy/config"
	"github.com/jpillora/ipfilter"
	log "github.com/sirupsen/logrus"
	"github.com/vulcand/oxy/roundrobin"
	"github.com/vulcand/oxy/testutils"
	"github.com/vulcand/oxy/utils"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func NewIPRouter(next http.Handler, conf *config.Config) *IPRouter {
	f := ipfilter.New(ipfilter.Options{
		BlockByDefault: conf.CanaryConfig.IPWhitelist != "",
	})

	f.BlockIP(conf.CanaryConfig.IPBlacklist)
	f.AllowIP(conf.CanaryConfig.IPWhitelist)

	return &IPRouter{
		next:       next,
		ipFilter:   f,
		canaryHost: testutils.ParseURI(conf.CanaryUpstreamHost),
		mainHost:   testutils.ParseURI(conf.MainUpstreamHost),
	}
}

type IPRouter struct {
	next                   http.Handler
	errHandler             utils.ErrorHandler
	canaryHost             *url.URL
	mainHost               *url.URL
	ipFilter               *ipfilter.IPFilter
	currentWeight          int
	requestRewriteListener roundrobin.RequestRewriteListener

	log *log.Logger
}

func (r *IPRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//if r.log.Level >= log.DebugLevel {
	//	logEntry := r.log.WithField("Request", utils.DumpHttpRequest(req))
	//	logEntry.Debug("vulcand/oxy/roundrobin/rr: begin ServeHttp on request")
	//	defer logEntry.Debug("vulcand/oxy/roundrobin/rr: completed ServeHttp on request")
	//}

	// make shallow copy of request before changing anything to avoid side effects
	newReq := *req

	ipAddr := getIp(req)

	if r.ipFilter.Allowed(ipAddr) {
		newReq.URL = r.canaryHost
	} else {
		newReq.URL = r.mainHost
	}

	//if r.log.Level >= log.DebugLevel {
	//	// log which backend URL we're sending this request to
	//	r.log.WithFields(log.Fields{"Request": utils.DumpHttpRequest(req), "ForwardURL": newReq.URL}).Debugf("vulcand/oxy/roundrobin/rr: Forwarding this request to URL")
	//}
	//
	//// Emit event to a listener if one exists
	//if r.requestRewriteListener != nil {
	//	r.requestRewriteListener(req, &newReq)
	//}

	r.next.ServeHTTP(w, &newReq)
}

// GetIp returns user's origin IP address
func getIp(r *http.Request) string {

	ip := r.Header.Get("X-Real-IP")
	if strings.Compare(ip, "") == 0 {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return "0.0.0.0"
		}
		return ip
	}

	return ip
}
