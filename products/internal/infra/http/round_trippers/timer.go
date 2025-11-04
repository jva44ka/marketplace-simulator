package round_trippers

import (
	"fmt"
	"net/http"
	"time"
)

type TimerRoundTripper struct {
	rt http.RoundTripper
}

func NewTimerRoundTipper(rt http.RoundTripper) http.RoundTripper {
	return &TimerRoundTripper{rt: rt}
}

func (trp *TimerRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	defer func(now time.Time) {
		fmt.Printf("SENT %s spent %s", r.URL.String(), time.Since(now))
	}(time.Now())

	return trp.rt.RoundTrip(r)
}
