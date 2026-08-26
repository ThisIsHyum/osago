package osago

import (
	"context"
	"fmt"
	"time"

	"github.com/ThisIsHyum/osago/client/colleges"
	"github.com/ThisIsHyum/osago/models"
)

func (c *Client) GetCollege(ctx context.Context, id int64) (*models.DtoCollegeResponse, error) {
	resp, err := c.c.Colleges.GetCollegesIDContext(ctx,
		colleges.NewGetCollegesIDParams().WithID(id))
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func (c *Client) GetColleges(ctx context.Context) ([]*models.DtoCollegeResponse, error) {
	resp, err := c.c.Colleges.GetCollegesContext(ctx, colleges.NewGetCollegesParams())
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func (c *Client) GetCollegeByName(ctx context.Context, name string) (*models.DtoCollegeResponse, error) {
	resp, err := c.c.Colleges.GetCollegesContext(ctx,
		colleges.NewGetCollegesParams().WithName(&name))
	if err != nil {
		return nil, err
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("campus %q: %w", name, ErrNotFound)
	}
	return resp.Payload[0], nil
}

func (c *Client) GetCollegeScheduleLast(ctx context.Context, id int64) (time.Time, error) {
	resp, err := c.c.Colleges.GetCollegesIDSchedulesLastContext(ctx,
		colleges.NewGetCollegesIDSchedulesLastParams().WithID(id))
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.DateOnly, resp.Payload.Date)
}
