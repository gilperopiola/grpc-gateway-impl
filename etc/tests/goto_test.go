package tests

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func allMyServerWithGotos() {
	requestHandlerFn := newRequestHandler()

	routes := []Route{
		{Path: "/hello", Auth: RouteAccessPublic},
		{Path: "/goodbye", Auth: RouteAccessUser},
	}

	for _, route := range routes {
		http.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			requestHandlerFn(route, w, r)
		})
	}

	http.ListenAndServe(":8080", nil)
}

type Route struct {
	Path string
	Auth RouteAccess
}

type RouteAccess string

const (
	RouteAccessInvalid RouteAccess = "invalid"
	RouteAccessPublic  RouteAccess = "public"
	RouteAccessUser    RouteAccess = "user"
	RouteAccessAdmin   RouteAccess = "admin"
)

func newRequestHandler() func(route Route, w http.ResponseWriter, r *http.Request) {

	// Here go all the non-request-specific stufs
	var limiter = rate.NewLimiter(rate.Limit(1), 1)

	// And this is the actual handler for every request
	return func(route Route, w http.ResponseWriter, r *http.Request) {

		// ⭐
		fmt.Printf("[%s] - New request for %s\n", time.Now().Format(time.RFC3339Nano), route.Path)

		var err error
		var errCode int

		goto LIMIT_REQUEST

	LIMIT_REQUEST:
		if !limiter.Allow() {
			err = fmt.Errorf("Sorry! We're at full capacity, go take a break and try again later.")
			errCode = http.StatusTooManyRequests
			goto END_REQUEST
		}
		goto AUTH_REQUEST

	AUTH_REQUEST:
		if route.Auth != RouteAccessPublic {
			err = fmt.Errorf("You are not authorized to access this.")
			errCode = http.StatusUnauthorized
			goto END_REQUEST
		}
		goto VALIDATE_REQUEST

	VALIDATE_REQUEST:
		goto HANDLE_REQUEST

	HANDLE_REQUEST:
		if route.Path == "/hello" {
			fmt.Fprintln(w, "Hello, world!")
		} else if route.Path == "/goodbye" {
			fmt.Fprintln(w, "Goodbye, world!")
		}

	END_REQUEST:
		if err != nil {
			http.Error(w, err.Error(), errCode)
		}
		return
	}
}

type HelloRequest struct {
	Name string `json:"name"`
}

type GoodbyeRequest struct {
	Name string `json:"name"`
}
