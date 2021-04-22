package serve

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/oschwald/geoip2-golang"
)

type geoIPHandler struct {
	cityDB     *geoip2.Reader
	asnDB      *geoip2.Reader
	remoteAddr net.IP
	family     string
}

func (h geoIPHandler) getRemoteAddr(r *http.Request) {
	var remoteAddr string

	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded == "" {
		remoteAddr = strings.Split(forwarded, ", ")[0]
	} else {
		remoteAddr = r.RemoteAddr
	}

	h.remoteAddr = net.ParseIP(remoteAddr)
	if h.remoteAddr.To4() != nil {
		h.family = "v4"
	} else {
		h.family = "v6"
	}
}

func (h geoIPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	getRemoteAddr(r)
	fmt.Println(r.Header)
	remoteAddr := strings.Split(r.RemoteAddr, ":")[0]

	recordCity, err := h.cityDB.City(ip)
	recordASN, err := h.asnDB.ASN(ip)

	if err != nil {
		fmt.Println(err)
	}

	city := fmt.Sprintf("%s - %s", recordCity.City.Names["pt-BR"], recordCity.Country.Names["pt-BR"])
	asn := fmt.Sprintf("%s (%d)", recordASN.AutonomousSystemOrganization, recordASN.AutonomousSystemNumber)
	fmt.Println(recordASN)
	out := fmt.Sprintf("IP: %s\n"+"City: %s\n"+"ASN: %s\n", ip, city, asn)

	w.Write([]byte(out))
}

func openGeoIPDB(path string) (*geoip2.Reader, error) {
	return geoip2.Open(path)
}

func (p *ServeCommand) Execute(args []string) error {
	cityDB, err := openGeoIPDB(p.CityPath)
	asnDB, err := openGeoIPDB(p.ASNPath)
	if err != nil {
		return err
	}

	h := geoIPHandler{}
	h.cityDB = cityDB
	h.asnDB = asnDB

	r := mux.NewRouter()

	r.PathPrefix("/").Handler(h)

	log.Fatal(http.ListenAndServe(p.Listen, r))
	return nil
}
