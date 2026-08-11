package osago

import (
	"context"
	"fmt"

	"github.com/ThisIsHyum/osago/client/campuses"
	"github.com/ThisIsHyum/osago/models"
)

func (c *Client) GetCampus(ctx context.Context, id int64) (*models.DtoCampusResponse, error) {
	resp, err := c.c.Campuses.GetCampusesIDContext(ctx, campuses.NewGetCampusesIDParams().WithID(id))
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func (c *Client) GetCampusByName(ctx context.Context, collegeID int64, name string) (*models.DtoCampusResponse, error) {
	resp, err := c.c.Campuses.GetCollegesCollegeIDCampusesContext(
		ctx, campuses.NewGetCollegesCollegeIDCampusesParams().
			WithCollegeID(int64(collegeID)).WithName(&name))
	if err != nil {
		return nil, err
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("campus %q: %w", name, ErrNotFound)
	}
	return resp.Payload[0], nil
}

func (c *Client) GetCampuses(ctx context.Context, collegeID int64) ([]*models.DtoCampusResponse, error) {
	resp, err := c.c.Campuses.GetCollegesCollegeIDCampusesContext(
		ctx, campuses.NewGetCollegesCollegeIDCampusesParams().WithCollegeID(int64(collegeID)))
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}
