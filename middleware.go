package inertia

import "net/http"

// Middleware function.
func (i *Inertia) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Inertia") == "" {
			next.ServeHTTP(w, r)

			return
		}

		if r.Method == "GET" && r.Header.Get("X-Inertia-Version") != i.version {
			w.Header().Set("X-Inertia-Location", i.Url+r.RequestURI)
			w.WriteHeader(http.StatusConflict)

			return
		}

		next.ServeHTTP(w, r)
	})
}
