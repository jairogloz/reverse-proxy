package middleware

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(next http.Handler) http.HandlerFunc
type MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc

// RequestLoggerMiddleware logs every incoming request to the console.
// It logs the start time, HTTP method, and URL path when the request starts,
// and logs the completion time and elapsed duration after the request is processed.
func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s - %s %s %s %v", r.Method, r.RemoteAddr, r.URL.Path, r.Proto, time.Since(start))
	}
}

// RequireAuthMiddleware ensures that the request is authenticated before proceeding.
// It checks for the presence of an "Authorization" header with a specific token value ("Bearer token").
// If the token does not match, it responds with an HTTP 401 Unauthorized status and terminates the request.
// If the token is valid, it allows the request to pass through to the next handler.
func RequireAuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check if the user is authenticated
		token := r.Header.Get("Authorization")
		if token != "Bearer token" {
			http.Error(w, "Unautorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// MiddlewareChain creates a chain of middlewares, allowing multiple middleware functions to be applied in sequence.
// It takes a variadic number of middleware functions as input and applies them in reverse order,
// meaning that the last middleware in the list will be executed first and the first middleware will be executed last.
// It returns a new http.HandlerFunc that applies the chain of middlewares to the given next handler.
func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}
}

func MiddlewareFuncChain(middlewares ...MiddlewareFunc) MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}
}

// CORSMiddleware añade las cabeceras necesarias para permitir CORS
func CORSMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") 
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
		
        next.ServeHTTP(w, r)
    })
}

// CORSMiddleware añade las cabeceras necesarias para permitir CORS
func CORSMiddlewareHandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") 
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
		
		next.ServeHTTP(w, r)
    })
}


var MiddlewaresList = map[string]Middleware{
"Logger": RequestLoggerMiddleware,
"Auth":   RequireAuthMiddleware,
}

