package osago

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ThisIsHyum/osago/types"
)

func (c *Client) GetCollege(ctx context.Context, collegeID uint) (types.College, error) {
	var college types.College
	err := c.doReq(ctx, http.MethodGet, fmt.Sprintf("/colleges/%d", collegeID), nil, &college)
	return college, err
}

func (c *Client) GetCollegeByName(ctx context.Context, name string) (types.College, error) {
	var colleges []types.College
	err := c.doReq(ctx, http.MethodGet, fmt.Sprintf("/colleges?name=%s", url.QueryEscape(name)), nil, &colleges)
	if err != nil {
		return types.College{}, err
	}
	if len(colleges) == 0 {
		return types.College{}, fmt.Errorf("college %q not found", name)
	}
	return colleges[0], nil
}

func (c *Client) GetColleges(ctx context.Context) ([]types.College, error) {
	var colleges []types.College
	err := c.doReq(ctx, http.MethodGet, "/colleges", nil, &colleges)
	return colleges, err
}
