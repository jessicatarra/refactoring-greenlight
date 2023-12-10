package service

import (
	"log/slog"
	"net/http"
)

func (s service) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		s.logger.Info("request", slog.Group("properties",
			"request_remote_addr", request.RemoteAddr,
			"request_proto", request.Proto,
			"request_method", request.Method,
			"request_url_request_uri", request.URL.RequestURI()),
		)

		next.ServeHTTP(writer, request)
	})
}
