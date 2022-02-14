package config

import "github.com/spf13/viper"

type Config struct {
	MainUpstreamHost      string       `json:"mainUpstreamHost,omitempty"`    // old cluster
	MainUpstreamPattern   string       `json:"mainUpstreamPattern,omitempty"` // old cluster
	CanaryUpstreamHost    string       `json:"canaryUpstreamHost,omitempty"`
	CanaryUpstreamPattern string       `json:"canaryUpstreamPattern,omitempty"`
	CanaryConfig          CanaryConfig `json:"canaryConfig,omitempty"`
}

type CanaryConfig struct {
	RouteType      string `json:"routeType"`
	TrafficPercent int    `json:"trafficPercent,omitempty"`
	IPWhitelist    string `json:"ipWhitelist,omitempty"`
	IPBlacklist    string `json:"ipBlacklist,omitempty"`
}

func GetConfig() *Config {
	viper.SetConfigFile("config/config.json")
	must(viper.ReadInConfig())
	var config Config
	must(viper.Unmarshal(&config))
	return &config
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
