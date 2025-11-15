package osago

import "net/http"

type AuthTransport struct {
	transport http.RoundTripper
	token     string
}

func (at *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+at.token)
	return at.transport.RoundTrip(req)
}

func newAuthTransport(token string) *AuthTransport {
	return &AuthTransport{token: token, transport: http.DefaultTransport}
}
