package ksi

import "net/http"

type response struct {
	status  int
	body    any
	headers http.Header
}

func (r *response) Ok(body any) *response {
	r.status = 200
	r.body = body
	return r
}

func (r *response) Created(body any) *response {
	r.status = 201
	r.body = body
	return r
}

func (r *response) NoContent() *response {
	r.status = 204
	r.body = nil
	return r
}

// NewResponse is a default response with status 200, body nil, and no headers
func NewResponse() *response {
	r := response{status: 200, body: nil, headers: http.Header{}}
	return &r
}

func (r *response) Status(s int) *response {
	r.status = s
	return r
}

func (r *response) Body(b any) *response {
	r.body = b
	return r
}

func (r *response) SetHeader(key string, val string) *response {
	r.headers.Add(key, val)
	return r
}
