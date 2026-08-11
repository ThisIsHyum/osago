package osago

import (
	"errors"
	"net/http"
	"time"

	"github.com/ThisIsHyum/osago/client"
	httptransport "github.com/go-openapi/runtime/client"
)

var ErrNotFound = errors.New("not found")

type Client struct {
	c *client.OpenScheduleAPI
}

func NewClient(url string, timeout time.Duration) *Client {
	c := httptransport.NewWithClient(url, "", client.DefaultSchemes, &http.Client{
		Timeout: timeout,
	})
	return &Client{c: client.New(c, nil)}
}
