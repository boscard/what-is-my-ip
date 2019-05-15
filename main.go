package main

import (
	"fmt"
	"os"
	"net"
	"net/http"
	"strings"
)

type Configuration struct {
	ListenPort	string
	ListenAddress	string
}

func GetClientsIPAddress(r *http.Request) string {
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

func RespondWithPublicIPAddress(w http.ResponseWriter, r *http.Request) {
        clinetsPublicIpAddress := GetClientsIPAddress(r)
        fmt.Fprintf(w, "%v\n", clinetsPublicIpAddress)
}

func getConfigurationFromEnvVars() Configuration {
	var conf Configuration
	conf.ListenPort = os.Getenv("WIMI_PORT")
	if conf.ListenPort == "" {
		conf.ListenPort = "8080"
	}
	conf.ListenAddress = os.Getenv("WIMI_LISTEN_ADDRESS")

	return conf
}

func main() {
	conf := getConfigurationFromEnvVars()
	http.HandleFunc("/", RespondWithPublicIPAddress)
	fmt.Println ("Starting server at port", conf.ListenAddress + ":" + conf.ListenPort)
	if err := http.ListenAndServe(conf.ListenAddress + ":" + conf.ListenPort, nil); err != nil {
		panic(err)
	}
}
