package api

import (
	"net/http"
	"strings"
)

type Visibility int

const (
	Public Visibility = iota
	Internal
)

type Route struct {
	visibility Visibility
	pattern    string
	handler    http.Handler
}

func NewRouteFunc(v Visibility, p string, f func(http.ResponseWriter, *http.Request)) *Route {
	return NewRoute(v, p, http.HandlerFunc(f))
}

func NewRoute(v Visibility, p string, h http.Handler) *Route {
	return &Route{
		visibility: v,
		pattern:    p,
		handler:    h,
	}
}

func (r *Route) getMethod() string {
	parts := strings.Split(r.pattern, " ")
	if len(parts) >= 2 {
		return parts[0]
	}
	return ""
}

func (r *Route) getUrl() (url string) {
	parts := strings.Split(r.pattern, " ")
	switch len(parts) {
	case 1:
		url = parts[0]
	case 2:
		url = parts[1]
	default:
		url = "/"
	}
	return url
}
