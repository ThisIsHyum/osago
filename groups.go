package osago

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ThisIsHyum/osago/types"
)

func (c *Client) GetGroup(ctx context.Context, groupID uint) (types.StudentGroup, error) {
	var group types.StudentGroup
	err := c.doReq(ctx, http.MethodGet, fmt.Sprintf("/groups/%d", groupID), nil, &group)
	return group, err
}

func (c *Client) GetGroupByName(ctx context.Context, campusID uint, name string) (types.StudentGroup, error) {
	var groups []types.StudentGroup
	err := c.doReq(ctx, http.MethodGet, fmt.Sprintf("/campuses/%d/groups?name=%s",
		campusID, url.QueryEscape(name)), nil, &groups)
	if err != nil {
		return types.StudentGroup{}, err
	}
	if len(groups) == 0 {
		return types.StudentGroup{}, err
	}
	return groups[0], err
}

func (c *Client) GetGroups(ctx context.Context, campusID uint) ([]types.StudentGroup, error) {
	var groups []types.StudentGroup
	err := c.doReq(ctx, http.MethodGet, fmt.Sprintf("/campuses/%d/groups", campusID), nil, &groups)
	return groups, err
}

func (c *Client) GetGroupsByCollegeID(ctx context.Context, collegeID uint) ([]types.StudentGroup, error) {
	var groups []types.StudentGroup
	err := c.doReq(ctx, http.MethodGet, fmt.Sprintf("/colleges/%d/groups", collegeID), nil, &groups)
	return groups, err
}
