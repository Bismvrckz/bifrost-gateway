package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
)

type contextKey string

const requestIDKey = contextKey("requestID")

var SessionManager *scs.SessionManager

func InitSession() {
	// Initialize the session manager.
	SessionManager = scs.New()
	SessionManager.Lifetime = 24 * time.Hour // Set session lifetime
	SessionManager.Cookie.Persist = true     // Keep session after browser close
	SessionManager.Cookie.Secure = true      // Use true in production
}

func GetRequestID(ctx context.Context) string {
	// .Value() returns an interface{}, so we need to assert the type.
	id, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return "unknown"
	}
	return id
}

func RequestIdMiddleware(next http.Handler) http.Handler {
	// This http.HandlerFunc is the actual middleware
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the X-Request-ID header from the request.
		// This is useful if a proxy (like Nginx) already set an ID.
		requestID := r.Header.Get("X-Request-ID")

		// 2. If it's empty, generate a new one.
		if requestID == "" {
			requestID = uuid.NewString()
		}

		// 3. Set the ID on the response header.
		// This lets the client see what the ID was.
		w.Header().Set("X-Request-ID", requestID)

		// 4. Create a new context with this request ID
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		// 5. Call the next handler in the chain, passing the new context.
		// r.WithContext(ctx) returns a shallow copy of r with the new context.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
