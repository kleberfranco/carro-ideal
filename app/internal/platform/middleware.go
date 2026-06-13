package platform

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"log/slog"
	"math"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"carro-ideal/app/internal/response"

	"github.com/go-chi/chi/v5/middleware"
)

const csrfCookieName = "carro_csrf"

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusRecorder) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusRecorder) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	count, err := w.ResponseWriter.Write(body)
	w.bytes += count
	return count, err
}

func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			recorder := &statusRecorder{ResponseWriter: w}
			next.ServeHTTP(recorder, r)
			if recorder.status == 0 {
				recorder.status = http.StatusOK
			}
			logger.InfoContext(r.Context(), "http request",
				"request_id", middleware.GetReqID(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
				"status", recorder.status,
				"bytes", recorder.bytes,
				"duration_ms", time.Since(started).Milliseconds(),
				"remote_ip", clientIP(r),
			)
		})
	}
}

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.ErrorContext(r.Context(), "panic recovered",
						"request_id", middleware.GetReqID(r.Context()),
						"panic", recovered,
						"stack", string(debug.Stack()),
					)
					response.Error(w, http.StatusInternalServerError, "erro interno do servidor", "INTERNAL_ERROR")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func RequestIDHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", middleware.GetReqID(r.Context()))
		next.ServeHTTP(w, r)
	})
}

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' https://cdn.jsdelivr.net; script-src 'self' https://code.jquery.com https://cdn.jsdelivr.net; connect-src 'self'; font-src 'self' https://cdn.jsdelivr.net; img-src 'self' data:")
		next.ServeHTTP(w, r)
	})
}

func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[origin] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && allowed[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			}
			if r.Method == http.MethodOptions {
				if origin == "" || !allowed[origin] {
					response.Error(w, http.StatusForbidden, "origem não permitida", "CORS_FORBIDDEN")
					return
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type visitor struct {
	tokens   float64
	lastSeen time.Time
}

func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	var mu sync.Mutex
	visitors := map[string]visitor{}
	refillPerSecond := float64(limit) / window.Seconds()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			ip := clientIP(r)
			mu.Lock()
			current := visitors[ip]
			if current.lastSeen.IsZero() {
				current = visitor{tokens: float64(limit), lastSeen: now}
			} else {
				current.tokens = math.Min(float64(limit), current.tokens+now.Sub(current.lastSeen).Seconds()*refillPerSecond)
				current.lastSeen = now
			}
			allowed := current.tokens >= 1
			if allowed {
				current.tokens--
			}
			visitors[ip] = current
			mu.Unlock()

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			remaining := int(current.tokens)
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			if !allowed {
				retryAfter := int(math.Ceil(1 / refillPerSecond))
				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				response.Error(w, http.StatusTooManyRequests, "limite de requisições excedido", "RATE_LIMITED")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CSRF(secureCookie bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(csrfCookieName)
			if err != nil || cookie.Value == "" {
				token, tokenErr := randomToken()
				if tokenErr != nil {
					response.Error(w, http.StatusInternalServerError, "falha ao iniciar proteção CSRF", "INTERNAL_ERROR")
					return
				}
				cookie = &http.Cookie{
					Name: csrfCookieName, Value: token, Path: "/", Secure: secureCookie,
					HttpOnly: false, SameSite: http.SameSiteStrictMode,
				}
				http.SetCookie(w, cookie)
			}

			if isStateChanging(r.Method) {
				header := r.Header.Get("X-CSRF-Token")
				if header == "" || subtle.ConstantTimeCompare([]byte(header), []byte(cookie.Value)) != 1 {
					response.Error(w, http.StatusForbidden, "token CSRF inválido", "CSRF_INVALID")
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func randomToken() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func isStateChanging(method string) bool {
	return method == http.MethodPost || method == http.MethodPut ||
		method == http.MethodPatch || method == http.MethodDelete
}

func clientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]); forwarded != "" {
		return forwarded
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		response.Error(w, http.StatusNotFound, "recurso não encontrado", "NOT_FOUND")
		return
	}
	http.Error(w, "Página não encontrada.", http.StatusNotFound)
}
