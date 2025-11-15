package osago

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ThisIsHyum/osago/types"
)

func (c *Client) GetSchedulesForToday(ctx context.Context, groupId uint) ([]types.Schedule, error) {
	return c.getSchedules(ctx, groupId, "day=today")
}

func (c *Client) GetSchedulesForTomorrow(ctx context.Context, groupId uint) ([]types.Schedule, error) {
	return c.getSchedules(ctx, groupId, "day=tomorrow")
}

func (c *Client) GetSchedulesForDate(ctx context.Context, groupId uint, date time.Time) ([]types.Schedule, error) {
	return c.getSchedules(ctx, groupId, "date="+url.QueryEscape(date.Format("02-01-2006")))
}

func (c *Client) GetSchedulesForPreviousWeek(ctx context.Context, groupId uint) ([]types.Schedule, error) {
	return c.getSchedules(ctx, groupId, "week=previous")
}

func (c *Client) GetSchedulesForCurrentWeek(ctx context.Context, groupId uint) ([]types.Schedule, error) {
	return c.getSchedules(ctx, groupId, "week=current")
}

func (c *Client) GetSchedulesForNextWeek(ctx context.Context, groupId uint) ([]types.Schedule, error) {
	return c.getSchedules(ctx, groupId, "week=next")
}

func (c *Client) GetSchedulesForWeekdayOfPreviousWeek(ctx context.Context, groupId uint, weekday time.Weekday) ([]types.Schedule, error) {
	return c.getSchedules(ctx, groupId, fmt.Sprintf("week=previous&weekday=%s", weekday))
}

func (c *Client) GetSchedulesForWeekday(ctx context.Context, groupId uint, weekday time.Weekday) ([]types.Schedule, error) {
	return c.getSchedules(ctx, groupId, fmt.Sprintf("week=current&weekday=%s", weekday))
}

func (c *Client) GetSchedulesForWeekdayOfNextWeek(ctx context.Context, groupId uint, weekday time.Weekday) ([]types.Schedule, error) {
	return c.getSchedules(ctx, groupId, fmt.Sprintf("week=next&weekday=%s", weekday))
}

func (c *Client) getSchedules(ctx context.Context, groupId uint, query string) ([]types.Schedule, error) {
	var schedules []types.Schedule
	path := fmt.Sprintf("/groups/%d/schedules?%s", groupId, query)
	if err := c.doReq(ctx, http.MethodGet, path, nil, &schedules); err != nil {
		return nil, err
	}
	return schedules, nil
}
