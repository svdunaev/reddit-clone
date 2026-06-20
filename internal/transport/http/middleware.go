package httptransport

import (
	"net/http"
	"time"
)

func (s *Server) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		s.logger.Info("new request", "method:", r.Method, "path: ", r.URL.Path)
		next.ServeHTTP(w, r)
		s.logger.Info("request completed", "method:", r.Method, "path: ", r.URL.Path, "duration: ", time.Since(start))
	})
}
