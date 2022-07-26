package main

import (
	"net/http"
	"strings"
)

func IpLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIp := r.RemoteAddr
		realIP := strings.Split(clientIp, ":")[0]

		// add to processing ips
		mu.Lock()
		if _, ok := serverConfig.ProcessingIps[realIP]; !ok {
			serverConfig.ProcessingIps[realIP] = make(chan struct{}, serverConfig.IpLimit)
		}
		mu.Unlock()

		serverConfig.ProcessingIps[realIP] <- struct{}{}

		next.ServeHTTP(w, r)

		<-serverConfig.ProcessingIps[realIP]
	})
}
