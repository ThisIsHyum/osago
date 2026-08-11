package osago

import (
	"context"
	"time"

	"github.com/ThisIsHyum/osago/client/admin"
	"github.com/ThisIsHyum/osago/models"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
)

type AdminClient struct {
	*Client
	auth runtime.ClientAuthInfoWriter
}

func NewAdminClient(url, token string, timeout time.Duration) *AdminClient {
	return &AdminClient{Client: NewClient(url, timeout), auth: httptransport.BearerToken(token)}
}

func (c *AdminClient) NewParser(ctx context.Context, collegeName string, campusNames []string) (string, error) {
	params := admin.NewPostAdminParserParams().WithNewParserRequest(
		&models.DtoNewParserRequest{CollegeName: collegeName, CampusNames: campusNames})

	resp, err := c.c.Admin.PostAdminParserContext(ctx, params, c.auth)
	if err != nil {
		return "", err
	}
	return resp.Payload.Token, nil
}

func (c *AdminClient) DeleteParser(ctx context.Context, id int64) error {
	_, err := c.c.Admin.DeleteAdminParserIDContext(ctx,
		admin.NewDeleteAdminParserIDParams().WithID(id), c.auth)
	return err
}
