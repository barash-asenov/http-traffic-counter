package main

import (
	"log"
	"net"
	"net/http"
	"strings"
)

func IpLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIp := r.RemoteAddr
		realIP := string(net.ParseIP(strings.Replace(clientIp, " ", "", -1)))

		// add to processing ips
		_, ok := serverConfig.ProcessingIps[realIP]

		if !ok {
			log.Println(serverConfig.ProcessingIps)
			serverConfig.ProcessingIps[realIP] = make(chan struct{}, serverConfig.IpLimit)
			log.Println(serverConfig.ProcessingIps)
		}

		serverConfig.ProcessingIps[realIP] <- struct{}{}

		next.ServeHTTP(w, r)

		<-serverConfig.ProcessingIps[realIP]
	})
}
