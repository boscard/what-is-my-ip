package main

import (
	"net/http"
	"whatismyip/whatismyip"
)

func main() {
	http.HandleFunc("/", whatismyip.RespondWithPublicIPAddress)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
