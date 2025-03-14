package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

const (
	DefaultAllowHeader = "Access-Control-Allow-Origin, Content-Type, Authorization"
	AllowOriginHeader  = "Access-Control-Allow-Origin"
	AllowMethodsHeader = "Access-Control-Allow-Methods"
	AllowHeadersHeader = "Access-Control-Allow-Headers"
	AllowOriginValue   = "*"
)

type MuxBuilder struct {
	logger       *slog.Logger
	allowOrigin  string
	allowHeaders string
}

func NewMuxBuilder(logger *slog.Logger, opts ...MuxOption) *MuxBuilder {
	mux := &MuxBuilder{
		logger:       logger,
		allowOrigin:  AllowOriginValue,
		allowHeaders: DefaultAllowHeader,
	}
	for _, opt := range opts {
		opt(mux)
	}
	return mux
}

type MuxOption func(*MuxBuilder)

func WithAllowOrigin(o string) MuxOption {
	return func(b *MuxBuilder) {
		b.allowOrigin = o
	}
}

func WithAllowHeaders(h []string) MuxOption {
	return func(b *MuxBuilder) {
		b.allowHeaders = strings.Join(h, ", ")
	}
}

type optionMap map[string][]string

func (o optionMap) addOption(url, method string) {
	if _, ok := o[url]; !ok {
		o[url] = make([]string, 0)
	}
	o[url] = append(o[url], method)
}

func (b *MuxBuilder) Build(r []*Route) *http.ServeMux {
	mux := http.NewServeMux()
	optionsRoutes := make(optionMap)
	for _, route := range r {
		url := route.getUrl()
		b.logger.Info("new route added", "route", route.pattern)
		mux.Handle(route.pattern, route.handler)
		if route.visibility == Public {
			optionsRoutes.addOption(url, route.getMethod())
			b.logger.Debug("Route for OPTIONS method as route is public", "route", route.pattern)
		}
	}
	b.handleOptions(mux, optionsRoutes)
	return mux
}

func (b *MuxBuilder) handleOptions(mux *http.ServeMux, optionsRoutes optionMap) {
	for url, methods := range optionsRoutes {
		h := newOptionsHandler(methods, b.allowHeaders, b.allowOrigin)
		mux.HandleFunc(fmt.Sprintf("%s %s", http.MethodOptions, url), h)
		mux.HandleFunc(fmt.Sprintf("%s %s", http.MethodHead, url), h)
	}
}

func newOptionsHandler(m []string, h, a string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(AllowMethodsHeader, strings.Join(m, ", "))
		w.Header().Set(AllowHeadersHeader, h)
		w.Header().Set(AllowOriginHeader, a)
		w.WriteHeader(http.StatusOK)
	}
}
