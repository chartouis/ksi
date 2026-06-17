// Package ksi is an opinionated wrapper around net/http
package ksi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime/debug"
)

type ksi struct {
	addr      string
	mux       *http.ServeMux
	preChain  *funChain // pre funChain is a slice of funcs that are run before the given handler
	postchain *funChain // post funChain is a slice of funcs that are run after the given handler
}

type KsiFunc func(*http.Request) (Response, error)

func NewKsi(addr string) *ksi {
	k := ksi{addr: addr, mux: http.NewServeMux()}
	k.preChain = NewFunChain()
	k.postchain = NewFunChain()
	return &k
}

func (k *ksi) SetPreChain(f ...ChainLink) {
	k.preChain = NewFunChain(f...)
}

func (k *ksi) SetPostChain(f ...ChainLink) {
	k.postchain = NewFunChain(f...)
}

func (k *ksi) Start() error {
	log.Print("Server started on - " + k.addr)
	return http.ListenAndServe(k.addr, k.mux)
}

func (k *ksi) Handle(pattern string, handler KsiFunc) {
	k.mux.HandleFunc(pattern, k.Middleware(handler))
	log.Print("Listening on - " + pattern)
}

func (k *ksi) Middleware(f KsiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer forceRecover(w)
		if !k.preChain.runAll(w, r) {
			return
		}
		res, err := f(r)
		if err != nil {
			if httpErr, ok := errors.AsType[HTTPError](err); ok {
				log.Print(httpErr.Status, " : ", httpErr.Message)
				WriteError(w, httpErr)
				return
			} else {
				WriteError(w, HTTPError{Status: 400, Message: ""})
				return
			}
		}

		w.WriteHeader(res.Status)
		mergeHeaders(w.Header(), res.Headers)
		if err := json.NewEncoder(w).Encode(res.Body); err != nil && res.Body != nil {
			log.Panic("Couldn't parse body into JSON")
		}

		if !k.postchain.runAll(w, r) {
			return
		}
	}
}

func mergeHeaders(dst, src http.Header) {
	for k, vv := range src {
		dst.Del(k)
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func forceRecover(w http.ResponseWriter) {
	if err := recover(); err != nil {
		log.Printf("error: %v\nstacktrace:\n%s", err, debug.Stack())
		http.Error(w, "Unexpected Error", http.StatusInternalServerError)
	}
}
