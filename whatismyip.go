package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func getClientsIPAddress(r *http.Request) string {
	addressFromProxyHeader := net.ParseIP(r.Header.Get("X-Forwarded-For"))
	if addressFromProxyHeader != nil {
		return addressFromProxyHeader.String()
	} else {
		if strings.HasPrefix(r.RemoteAddr, "[") {
			return strings.Split(strings.Trim(r.RemoteAddr, "["), "]:")[0]
		} else {
			return strings.Split(r.RemoteAddr, ":")[0]
		}
	}
}

func respondWithPublicIPAddress(w http.ResponseWriter, r *http.Request) {
	clinetsPublicIpAddress := getClientsIPAddress(r)
	fmt.Fprintf(w, "%v\n", clinetsPublicIpAddress)
}

func main() {
	http.HandleFunc("/", respondWithPublicIPAddress)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
