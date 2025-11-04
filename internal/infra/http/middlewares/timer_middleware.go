package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

type TimerMiddleware struct {
	h http.Handler
}

func NewTimerMiddleware(h http.Handler) http.Handler {
	return &TimerMiddleware{h: h}
}

func (m *TimerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(now time.Time) {
		fmt.Printf("%s spent %s", r.URL.String(), time.Since(now))
	}(time.Now())

	//body, err := io.ReadAll(r.Body)
	// Replace the body with a new reader after reading from the original
	//r.Body = io.NopCloser(bytes.NewBuffer(body))

	m.h.ServeHTTP(w, r)
}
