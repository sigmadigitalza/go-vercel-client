package vercel_client

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

var (
	MissingTokenError  = errors.New("missing Vercel API token")
)

type transport struct {
	token               string
	teamId              string
	underlyingTransport http.RoundTripper
}

func newTransport(token string, teamId string) *transport {
	return &transport{
		token:               token,
		teamId:              teamId,
		underlyingTransport: http.DefaultTransport,
	}
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	req.Header.Add("Content-Type", "application/json")

	// Only add teamId if one was supplied
	if t.teamId != "" {
		q := req.URL.Query()
		q.Add("teamId", t.teamId)
		req.URL.RawQuery = q.Encode()
	}

	return t.underlyingTransport.RoundTrip(req)
}

func defaultClient() (*http.Client, error) {
	token := os.Getenv("VERCEL_TOKEN")
	if token == "" {
		return nil, MissingTokenError
	}

	teamId := os.Getenv("VERCEL_TEAM_ID")

	return &http.Client{
		Transport: newTransport(token, teamId),
	}, nil
}
