package serve

import (
	"encoding/json"
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
	password   string
}

type response struct {
	IP        net.IP
	City      string
	ASN       string
	Latitude  float64
	Longitude float64
}

func (h *geoIPHandler) getRemoteAddr(r *http.Request) {

	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		h.remoteAddr = net.ParseIP(strings.Split(forwarded, ", ")[0])
	} else {
		h.remoteAddr = net.ParseIP(r.RemoteAddr)
	}

	if h.remoteAddr.To4() != nil {
		h.family = "v4"
	} else {
		h.family = "v6"
	}
}

func (h geoIPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.getRemoteAddr(r)

	if len(h.password) > 0 {
		r.ParseForm()
		if password, ok := r.Form["p"]; ok {
			if password[0] != h.password {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	recordCity, err := h.cityDB.City(h.remoteAddr)
	recordASN, err := h.asnDB.ASN(h.remoteAddr)

	if err != nil {
		fmt.Println(err)
	}

	city := fmt.Sprintf("%s - %s", recordCity.City.Names["pt-BR"], recordCity.Country.Names["pt-BR"])
	lat := recordCity.Location.Latitude
	lon := recordCity.Location.Longitude
	asn := fmt.Sprintf("%s (%d)", recordASN.AutonomousSystemOrganization, recordASN.AutonomousSystemNumber)
	out := response{
		IP:        h.remoteAddr,
		City:      city,
		ASN:       asn,
		Latitude:  lat,
		Longitude: lon,
	}
	j, err := json.Marshal(out)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write(j)
}

func openGeoIPDB(path string) (*geoip2.Reader, error) {
	return geoip2.Open(path)
}

func (p *ServeCommand) Execute(args []string) error {
	if len(p.CityPath) == 0 {
		return fmt.Errorf("Please specify GetLite2 City mmdb file path")
	}
	if len(p.ASNPath) == 0 {
		return fmt.Errorf("Please specify GetLite2 ASN mmdb file path")
	}
	cityDB, err := openGeoIPDB(p.CityPath)
	asnDB, err := openGeoIPDB(p.ASNPath)
	if err != nil {
		return err
	}

	h := geoIPHandler{}
	h.cityDB = cityDB
	h.asnDB = asnDB
	h.password = p.Password

	r := mux.NewRouter()

	r.PathPrefix("/").Handler(h)

	fmt.Printf("Listening on %s\n", p.Listen)
	log.Fatal(http.ListenAndServe(p.Listen, r))
	return nil
}
