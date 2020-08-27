package main

import (
	"net/http"
	pckg "github.com/boscard/what-is-my-ip/package"
)

func main() {
	http.HandleFunc("/",pckg.RespondWithPublicIPAddress)
	if err := http.ListenAndServe(":8080", nil); err!= nil {
		panic(err)
	}
}
