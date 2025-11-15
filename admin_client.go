package osago

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ThisIsHyum/osago/dto"
)

type AdminClient struct {
	*Client
	token string
}

func NewAdminClient(url, token string, timeout time.Duration) *AdminClient {
	c := NewClient(url, timeout)
	c.httpClient.Transport = newAuthTransport(token)
	return &AdminClient{
		Client: c,
		token:  token,
	}
}

func (c *AdminClient) NewParser(ctx context.Context, collegeName string, campusNames []string) (string, error) {
	request := dto.NewParserRequest{
		CollegeName: collegeName,
		CampusNames: campusNames,
	}
	var response dto.NewParserResponse
	err := c.doReq(ctx, http.MethodPost, "/admin/parser", request, &response)
	if err != nil {
		return "", err
	}
	return response.Token, nil
}

func (c *AdminClient) DeleteParser(ctx context.Context, id uint) error {
	return c.doReq(ctx, http.MethodDelete, fmt.Sprintf("/admin/parser/%d", id), nil, nil)
}
