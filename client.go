package osago

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/ThisIsHyum/osago/client"
	httptransport "github.com/go-openapi/runtime/client"
)

var ErrNotFound = errors.New("not found")

type Client struct {
	c *client.OpenScheduleAPI
}

func NewClient(URL string, timeout time.Duration) *Client {
	u, err := url.Parse(URL)
	if err != nil {
		panic(err)
	}
	c := httptransport.NewWithClient(u.Host, u.Path, []string{u.Scheme}, &http.Client{
		Timeout: timeout,
	})
	return &Client{c: client.New(c, nil)}
}
