package util

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type aggregateHandler struct {
	handle     http.Handler
	middleware []Middleware
}

func AggregateHandler(h http.Handler, mws ...Middleware) http.Handler {
	agg := aggregateHandler{
		handle:     h,
		middleware: mws,
	}
	for _, mw := range agg.middleware {
		agg.handle = mw(agg.handle)
	}
	return &agg
}

func (agg *aggregateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	agg.handle.ServeHTTP(w, r)
}
