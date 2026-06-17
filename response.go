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
