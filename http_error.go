package ksi

import "net/http"

type HTTPError struct {
	Status  int
	Message string
}

func (e HTTPError) Error() string {
	return e.Message
}

func WriteError(w http.ResponseWriter, e HTTPError) {
	http.Error(w, e.Message, e.Status)
}
