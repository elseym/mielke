package main

import (
	"net/http"
	"strings"
)

type route struct {
	path    string
	method  string
	headers map[string]string
	handler http.HandlerFunc
}

type router struct {
	routes []route
}

func (t *router) add(path, method string, headers map[string]string, handler http.HandlerFunc) {
	t.routes = append(t.routes, route{path, method, headers, handler})
}

func (t router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, o := range t.routes {
		matches := o.path == r.URL.Path && o.method == r.Method
		for n, v := range o.headers {
			matches = matches && strings.Contains(r.Header.Get(n), v)
		}
		if matches {
			o.handler(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
