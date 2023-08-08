package middlewares

import "net/http"

func Helmet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Xss-Protection", "0")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("X-Frame-Options", "SAMEORIGIN")
		w.Header().Add("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Add("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Add("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Add("Origin-Agent-Cluster", "?1")
		w.Header().Add("Referrer-Policy", "no-referrer")
		w.Header().Add("X-Dns-Prefetch-Control", "off")
		w.Header().Add("X-Download-Options", "noopen")
		w.Header().Add("X-Permitted-Cross-Domain-Policies", "none")

		next.ServeHTTP(w, r)
	})
}
