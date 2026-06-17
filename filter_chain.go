package ksi

import (
	"errors"
	"log"
	"net/http"
)

type funChain struct {
	Filters []ChainLink
}

type ChainLink func(http.ResponseWriter, *http.Request) error

func NewFunChain(f ...ChainLink) *funChain {
	return &funChain{Filters: f}
}

func (d *funChain) runAll(w http.ResponseWriter, r *http.Request) bool {
	for _, f := range d.Filters {
		if err := f(w, r); err != nil {
			if httpErr, ok := errors.AsType[HTTPError](err); ok {
				log.Print(httpErr.Status, " : ", httpErr.Message)
				WriteError(w, httpErr)
				return false
			}
		}
	}
	return true
}
