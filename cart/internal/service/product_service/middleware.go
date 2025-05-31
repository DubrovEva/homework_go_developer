package product_service

import (
	"net/http"
)

type Middleware struct {
	rt http.RoundTripper

	retryCount uint
}

func NewMiddleware(rt http.RoundTripper, retryCount uint) http.RoundTripper {
	return &Middleware{rt: rt, retryCount: retryCount}
}

func (m *Middleware) RoundTrip(r *http.Request) (*http.Response, error) {
	for i := uint(0); i < m.retryCount-1; i++ {
		resp, err := m.rt.RoundTrip(r)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != 420 && resp.StatusCode != 429 {
			return resp, err
		}
	}

	return m.rt.RoundTrip(r)
}
