// A very simple proxy that given a GET request from a configurable host, you will be redirected to a configurable url

package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"net"
	"net/http"
)

var routeMappings map[string]string

func main() {
	routeMappings = make(map[string]string)

	m := martini.Classic()
	m.Get("/myip", MyIP)

	m.Group("/proxy", func(r martini.Router) {
		r.Post("/add_route", AddRoute)
		r.Get("/dump", Dump)
	})

	m.Get("/**", Redirect)

    http.ListenAndServe(":8000", m)
}

func MyIP(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return
	}

	fmt.Fprintln(w, host)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return
	}

	url, present := routeMappings[host]
	
	fmt.Println("Host: ", host)

	if !present {
		fmt.Println(" - No mapping to redirect to.")
		return
	}

	if present {
	    fmt.Println(" - redirecting to: ", url)
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func Dump(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%#v", routeMappings)
}


func AddRoute(w http.ResponseWriter, r *http.Request) {
	redirectUrl := r.FormValue("redirect_url")
	deviceEndpoint := r.FormValue("device_endpoint")

	if deviceEndpoint == "" {
		return
	}

	routeMappings[deviceEndpoint] = redirectUrl
	fmt.Println("Created the route: %s->%s", deviceEndpoint, redirectUrl)
}
