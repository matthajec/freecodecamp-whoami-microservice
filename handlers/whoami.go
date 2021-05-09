package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type whoami struct {
	l *log.Logger
}

type userInfo struct {
	IP       string `json:"ipaddress"`
	Language string `json:"language"`
	Software string `json:"software"`
}

func NewWhoami(l *log.Logger) *whoami {
	return &whoami{l}
}

func (w *whoami) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	w.l.Printf("%s whoami", r.Method)
	addCORS(&rw)

	if r.Method == http.MethodGet {
		rw.Header().Add("Content-Type", "application/json")

		d := userInfo{
			IP:       getIp(r),
			Language: getLanguage(r),
			Software: getUserAgent(r),
		}

		e := json.NewEncoder(rw)
		err := e.Encode(d)
		if err != nil {
			http.Error(rw, "Failed to decode JSON", http.StatusInternalServerError)
		}
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func getIp(r *http.Request) string {
	ips := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
	return ips[len(ips)-1]
}

func getLanguage(r *http.Request) string {
	return r.Header.Get("Accept-Language")
}

func getUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

func addCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "https://www.freecodecamp.org")
}
