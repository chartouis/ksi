# ksi

A lightweight HTTP wrapper around Go's `net/http`. No external dependencies.

## Install

```bash
go get github.com/chartouis/ksi
```

Requires Go 1.26+.

## Usage

### Basic setup

```go
package main

import (
    "log"
    "github.com/chartouis/ksi"
)

func main() {
    k := ksi.NewKsi(":8080")
    k.Get("/hello", helloHandler)
    log.Fatal(k.Start())
}
```

### Handlers

A handler takes a `*http.Request` and returns a `Response` and an error.

```go
type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

func helloHandler(r *http.Request) (ksi.Response, error) {
    id := r.PathValue("id")
    return ksi.Ok(User{ID: id, Name: "Ernazar"}), nil
}
```

The response body is automatically encoded as JSON.

### Response helpers

```go
ksi.Ok(body)       // 200
ksi.Created(body)  // 201
ksi.NoContent()    // 204
```

To include custom headers:

```go
func handler(r *http.Request) (ksi.Response, error) {
    h := http.Header{}
    h.Set("X-Custom", "value")
    return ksi.WithHeaders(h).Ok(body), nil
}
```

### Returning errors

Return an `HTTPError` to send a specific status code and message:

```go
func handler(r *http.Request) (ksi.Response, error) {
    return ksi.Response{}, ksi.HTTPError{Status: 404, Message: "not found"}
}
```

Any other error results in a 400 response.

### HTTP methods

```go
k.Get("/users/{id}", getUser)
k.Post("/users", createUser)
k.Put("/users/{id}", updateUser)
k.Patch("/users/{id}", patchUser)
k.Delete("/users/{id}", deleteUser)
```

Path parameters use Go 1.22+ stdlib syntax and are accessed via `r.PathValue("id")`.

### Filter chains

Filters run before or after your handler. A filter returning an error stops the chain immediately.

```go
func authFilter(w http.ResponseWriter, r *http.Request) error {
    token := r.Header.Get("Authorization")
    if token == "" {
        return ksi.HTTPError{Status: 401, Message: "unauthorized"}
    }
    return nil
}

func loggingFilter(w http.ResponseWriter, r *http.Request) error {
    log.Printf("%s %s", r.Method, r.URL.Path)
    return nil
}

k.SetPreChain(authFilter, loggingFilter)   // runs before handler
k.SetPostChain(loggingFilter)              // runs after handler
```

Pre-chain filters can short-circuit the request. Post-chain filters are meant for logging, metrics, and cleanup — they cannot modify the response.
