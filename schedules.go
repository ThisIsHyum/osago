package osago

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ThisIsHyum/osago/types"
)

func (c *Client) GetScheduleForToday(ctx context.Context, groupId uint) (types.Schedule, error) {
	return c.getSchedule(ctx, groupId, "day=today")
}

func (c *Client) GetScheduleForTomorrow(ctx context.Context, groupId uint) (types.Schedule, error) {
	return c.getSchedule(ctx, groupId, "day=tomorrow")
}

func (c *Client) GetScheduleForDate(ctx context.Context, groupId uint, date time.Time) (types.Schedule, error) {
	return c.getSchedule(ctx, groupId, "date="+url.QueryEscape(date.Format("02-01-2006")))
}

func (c *Client) GetScheduleForWeekdayOfPreviousWeek(ctx context.Context, groupId uint, weekday time.Weekday) (types.Schedule, error) {
	return c.getSchedule(ctx, groupId, fmt.Sprintf("week=previous&weekday=%s", weekday))
}

func (c *Client) GetScheduleForWeekday(ctx context.Context, groupId uint, weekday time.Weekday) (types.Schedule, error) {
	return c.getSchedule(ctx, groupId, fmt.Sprintf("week=current&weekday=%s", weekday))
}

func (c *Client) GetScheduleForWeekdayOfNextWeek(ctx context.Context, groupId uint, weekday time.Weekday) (types.Schedule, error) {
	return c.getSchedule(ctx, groupId, fmt.Sprintf("week=next&weekday=%s", weekday))
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

func (c *Client) getSchedules(ctx context.Context, groupId uint, query string) ([]types.Schedule, error) {
	var schedules []types.Schedule
	path := fmt.Sprintf("/groups/%d/schedules?%s", groupId, query)
	if err := c.doReq(ctx, http.MethodGet, path, nil, &schedules); err != nil {
		return nil, err
	}
	return schedules, nil
}

func (c *Client) getSchedule(ctx context.Context, groupId uint, query string) (types.Schedule, error) {
	schedules, err := c.getSchedules(ctx, groupId, query)
	if err != nil {
		return types.Schedule{}, err
	}
	if len(schedules) == 0 {
		return types.Schedule{}, errors.New("schedule not found")
	}
	return schedules[0], nil
}
