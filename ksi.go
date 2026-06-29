// Package ksi is an opinionated wrapper around net/http
package ksi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
)

type ksi struct {
	addr      string
	mux       *http.ServeMux
	preChain  *funChain // pre funChain is a slice of funcs that are run before the given handler
	postChain *funChain // post funChain is a slice of funcs that are run after the given handler
}

type KsiFunc = any

func NewKsi(addr string) *ksi {
	k := ksi{addr: addr, mux: http.NewServeMux()}
	k.preChain = NewFunChain()
	k.postChain = NewFunChain()
	return &k
}

func (k *ksi) SetPreChain(f ...ChainLink) {
	k.preChain = NewFunChain(f...)
}

func (k *ksi) SetPostChain(f ...ChainLink) {
	k.postChain = NewFunChain(f...)
}

func (k *ksi) Start() error {
	log.Print(`
 ___  __        ________       ___     
|\  \|\  \     |\   ____\     |\  \    
\ \  \/  /|_   \ \  \___|_    \ \  \   
 \ \   ___  \   \ \_____  \    \ \  \  
  \ \  \\ \  \ __\|____|\  \  __\ \  \ 
   \ \__\\ \__|\__\____\_\  \|\__\ \__\
    \|__| \|__\|__|\_________\|__|\|__|
                  \|_________|         `)
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

		isFunc(f)
		validateHandler(reflect.TypeOf(f))

		if !k.preChain.runAll(w, r) {
			return
		}
		res, err := injectAndRun(w, r, f)
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

		mergeHeaders(w.Header(), res.Headers)
		w.WriteHeader(res.Status)
		if err := json.NewEncoder(w).Encode(res.Body); err != nil && res.Body != nil {
			log.Panic("Couldn't parse body into JSON")
		}

		if !k.postChain.runAll(w, r) {
			return
		}
	}
}

func isFunc(f any) {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		panic("ksi: Handler must be a function")
	}
}

func validateHandler(t reflect.Type) {
	structCount := 0
	for p := range t.Ins() {
		if p != reflect.TypeFor[*http.Request]() &&
			p != reflect.TypeFor[http.ResponseWriter]() {
			if p.Kind() != reflect.Struct {
				panic(fmt.Sprintf("ksi: unsupported parameter type: %s", p.Name()))
			}
			structCount++
			if structCount > 1 {
				panic("ksi: only one struct body parameter is allowed")
			}
		}
	}
}

func injectAndRun(w http.ResponseWriter, r *http.Request, f KsiFunc) (Response, error) {
	t := reflect.TypeOf(f)
	args := []reflect.Value{}
	for p := range t.Ins() {
		if p == reflect.TypeFor[*http.Request]() {
			args = append(args, reflect.ValueOf(r))
		} else if p == reflect.TypeFor[http.ResponseWriter]() {
			args = append(args, reflect.ValueOf(w))
		} else {
			ptr := reflect.New(p)
			if err := json.NewDecoder(r.Body).Decode(ptr.Interface()); err != nil {
				return Response{}, HTTPError{Status: 400, Message: "invalid request body"}
			}
			args = append(args, ptr.Elem())
		}
	}
	v := reflect.ValueOf(f)
	results := v.Call(args)
	resp := results[0].Interface().(Response)
	err, _ := results[1].Interface().(error)
	return resp, err
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
