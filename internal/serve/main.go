package serve

import (
	flags "github.com/jessevdk/go-flags"
)

type ServeCommand struct {
	Listen   string `short:"l" long:"listen" env:"MYIP_LISTEN" default:"0.0.0.0:8000" description:"What IP and port to listen to connections"`
	CityPath string `long:"geoip-city" env:"MYIP_GEOIP_CITY" description:"Path to MaxMind GeoIP City database"`
	ASNPath  string `long:"geoip-asn" env:"MYIP_GEOIP_ASN" description:"Path to MaxMind GeoIP ASN database"`
	Password string `short:"p" long:"password" env:"MYIP_GEOIP_PASSWORD" description:"Optional password to restrict requests"`
}

var (
	serveCommand ServeCommand
)

func Init(parser *flags.Parser) {
	_, err := parser.AddCommand(
		"serve",
		"Serve",
		"Start web service",
		&serveCommand)

	if err != nil {
		return
	}
}
