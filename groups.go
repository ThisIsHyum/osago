package osago

import (
	"context"
	"fmt"

	"github.com/ThisIsHyum/osago/client/groups"
	"github.com/ThisIsHyum/osago/models"
)

func (c *Client) GetGroup(ctx context.Context, id int64) (*models.DtoStudentGroupResponse, error) {
	resp, err := c.c.Groups.GetGroupsIDContext(ctx, groups.NewGetGroupsIDParams().WithID(id))
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func (c *Client) GetGroupByName(ctx context.Context, campusID int64, name string) (*models.DtoStudentGroupResponse, error) {
	resp, err := c.c.Groups.GetCampusesCampusIDGroupsContext(ctx,
		groups.NewGetCampusesCampusIDGroupsParams().WithCampusID(campusID).WithName(&name))
	if err != nil {
		return nil, err
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("group %q: %w", name, ErrNotFound)
	}
	return resp.Payload[0], nil
}

func (c *Client) GetGroups(ctx context.Context, campusID int64) ([]*models.DtoStudentGroupResponse, error) {
	resp, err := c.c.Groups.GetCampusesCampusIDGroupsContext(ctx,
		groups.NewGetCampusesCampusIDGroupsParams().WithCampusID(campusID))
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func (c *Client) GetGroupsByCollegeID(ctx context.Context, collegeID int64, name *string) ([]*models.DtoStudentGroupResponse, error) {
	resp, err := c.c.Groups.GetCollegesCollegeIDGroupsContext(ctx,
		groups.NewGetCollegesCollegeIDGroupsParams().WithCollegeID(collegeID).WithName(name))
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}
