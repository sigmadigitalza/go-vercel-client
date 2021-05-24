package vercel_client

import (
	"net/url"
)

type Client struct {
	Project *ProjectApi
}

func New() (*Client, error) {
	c, err := defaultClient()
	if err != nil {
		return nil, err
	}

	baseUrl := &url.URL{
		Scheme: "https",
		Host: "api.vercel.com",
	}

	return &Client{
		Project: &ProjectApi{
			c:       c,
			baseUrl: baseUrl,
		},
	}, nil
}
