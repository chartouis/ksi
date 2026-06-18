package ksi

import "net/http"

type Response struct {
	Status  int
	Body    any
	Headers http.Header
}

func Ok(body any) Response      { return Response{Status: 200, Body: body} }
func Created(body any) Response { return Response{Status: 201, Body: body} }
func NoContent() Response       { return Response{Status: 204} }

func (r *Response) Ok(body any) Response {
	r.Status = 200
	r.Body = body
	return *r
}

func (r *Response) Created(body any) Response {
	r.Status = 201
	r.Body = body
	return *r
}

func (r *Response) NoContent() Response {
	r.Status = 204
	r.Body = nil
	return *r
}

func WithHeaders(h http.Header) Response {
	return Response{Headers: h}
}
