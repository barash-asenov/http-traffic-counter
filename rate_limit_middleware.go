package main

import (
	"net/http"
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start doing the request, take one spot from the channel
		serverConfig.RequestLimit <- struct{}{}

		next.ServeHTTP(w, r)
		// release
		<-serverConfig.RequestLimit
	})
}
