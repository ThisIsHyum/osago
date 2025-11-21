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

var ErrScheduleNotFound = errors.New("schedule not found")

func (c *Client) GetScheduleForToday(ctx context.Context, groupID uint) (types.Schedule, error) {
	return c.getSchedule(ctx, groupID, "day=today")
}

func (c *Client) GetScheduleForTomorrow(ctx context.Context, groupID uint) (types.Schedule, error) {
	return c.getSchedule(ctx, groupID, "day=tomorrow")
}

func (c *Client) GetScheduleForDate(ctx context.Context, groupID uint, date time.Time) (types.Schedule, error) {
	return c.getSchedule(ctx, groupID, "date="+url.QueryEscape(date.Format("02-01-2006")))
}

func (c *Client) GetScheduleForWeekdayOfPreviousWeek(ctx context.Context, groupID uint, weekday time.Weekday) (types.Schedule, error) {
	return c.getSchedule(ctx, groupID, fmt.Sprintf("week=previous&weekday=%s", weekday))
}

func (c *Client) GetScheduleForWeekday(ctx context.Context, groupID uint, weekday time.Weekday) (types.Schedule, error) {
	return c.getSchedule(ctx, groupID, fmt.Sprintf("week=current&weekday=%s", weekday))
}

func (c *Client) GetScheduleForWeekdayOfNextWeek(ctx context.Context, groupID uint, weekday time.Weekday) (types.Schedule, error) {
	return c.getSchedule(ctx, groupID, fmt.Sprintf("week=next&weekday=%s", weekday))
}

func (c *Client) GetSchedulesForPreviousWeek(ctx context.Context, groupID uint) ([]types.Schedule, error) {
	return c.getSchedules(ctx, groupID, "week=previous")
}

func (c *Client) GetSchedulesForCurrentWeek(ctx context.Context, groupID uint) ([]types.Schedule, error) {
	return c.getSchedules(ctx, groupID, "week=current")
}

func (c *Client) GetSchedulesForNextWeek(ctx context.Context, groupID uint) ([]types.Schedule, error) {
	return c.getSchedules(ctx, groupID, "week=next")
}

func (c *Client) getSchedules(ctx context.Context, groupID uint, query string) ([]types.Schedule, error) {
	var schedules []types.Schedule
	path := fmt.Sprintf("/groups/%d/schedules?%s", groupID, query)
	if err := c.doReq(ctx, http.MethodGet, path, nil, &schedules); err != nil {
		return nil, err
	}
	return schedules, nil
}

func (c *Client) getSchedule(ctx context.Context, groupID uint, query string) (types.Schedule, error) {
	schedules, err := c.getSchedules(ctx, groupID, query)
	if err != nil {
		return types.Schedule{}, err
	}
	if len(schedules) == 0 {
		return types.Schedule{}, ErrScheduleNotFound
	}
	return schedules[0], nil
}
