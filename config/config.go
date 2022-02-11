package config

import "github.com/spf13/viper"

type Config struct {
	Backends  []Backend `json:"backends,omitempty"`
	RouteType string
}

type Backend struct {
	Name       string     `json:"name,omitempty"`
	UrlPattern string     `json:"urlPattern,omitempty"`
	Upstreams  []Upstream `json:"upstreams,omitempty"`
}

type Routing struct {
	UrlPattern  string `json:"urlPattern,omitempty"`
	IPWhitelist string `json:"ipWhitelist,omitempty"`
	IPBlacklist string `json:"ipBlacklist,omitempty"`
	RouterType  string `json:"routerType,omitempty"`
}

type Upstream struct {
	Name    string  `json:"name,omitempty"`
	Host    string  `json:"host,omitempty"`
	Weight  int     `json:"weight,omitempty"`
	Routing Routing `json:"routing"`
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
