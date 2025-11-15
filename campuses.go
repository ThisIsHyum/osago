package osago

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ThisIsHyum/osago/types"
)

func (c *Client) GetCampus(ctx context.Context, campusID uint) (types.Campus, error) {
	var campus types.Campus
	err := c.doReq(ctx, http.MethodGet, fmt.Sprintf("/campuses/%d", campusID), nil, &campus)
	return campus, err
}

func (c *Client) GetCampusByName(ctx context.Context, collegeID uint, name string) (types.Campus, error) {
	var campuses []types.Campus
	err := c.doReq(ctx, http.MethodGet, fmt.Sprintf("/colleges/%d/campuses?name=%s",
		collegeID, url.QueryEscape(name)), nil, &campuses)
	if err != nil {
		return types.Campus{}, err
	}
	if len(campuses) == 0 {
		return types.Campus{}, fmt.Errorf("campus %q not found", name)
	}
	return campuses[0], nil
}

func (c *Client) GetCampuses(ctx context.Context, collegeID uint) ([]types.Campus, error) {
	var campuses []types.Campus
	err := c.doReq(ctx, http.MethodGet, fmt.Sprintf("/colleges/%d/campuses", collegeID), nil, &campuses)
	return campuses, err
}
