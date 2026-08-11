package osago

import (
	"context"
	"time"

	"github.com/ThisIsHyum/osago/client/schedules"
	"github.com/ThisIsHyum/osago/models"
)

func (c *Client) GetScheduleForToday(ctx context.Context, groupID int64) (*models.DtoScheduleResponse, error) {
	return c.getSchedule(ctx, groupID, nil, nil, nil, ptr("today"))
}

func (c *Client) GetScheduleForTomorrow(ctx context.Context, groupID int64) (*models.DtoScheduleResponse, error) {
	return c.getSchedule(ctx, groupID, nil, nil, nil, ptr("tomorrow"))
}

func (c *Client) GetScheduleForDate(ctx context.Context, groupID int64, date time.Time) (*models.DtoScheduleResponse, error) {
	return c.getSchedule(ctx, groupID, &date, nil, nil, nil)
}

func (c *Client) GetScheduleForWeekdayOfPreviousWeek(ctx context.Context, groupID int64, weekday time.Weekday) (*models.DtoScheduleResponse, error) {
	return c.getSchedule(ctx, groupID, nil, &weekday, ptr("previous"), nil)
}

func (c *Client) GetScheduleForWeekday(ctx context.Context, groupID int64, weekday time.Weekday) (*models.DtoScheduleResponse, error) {
	return c.getSchedule(ctx, groupID, nil, &weekday, ptr("current"), nil)
}

func (c *Client) GetScheduleForWeekdayOfNextWeek(ctx context.Context, groupID int64, weekday time.Weekday) (*models.DtoScheduleResponse, error) {
	return c.getSchedule(ctx, groupID, nil, &weekday, ptr("next"), nil)
}

func (c *Client) GetSchedulesForPreviousWeek(ctx context.Context, groupID int64) ([]*models.DtoScheduleResponse, error) {
	return c.getSchedules(ctx, groupID, ptr("previous"))
}

func (c *Client) GetSchedulesForCurrentWeek(ctx context.Context, groupID int64) ([]*models.DtoScheduleResponse, error) {
	return c.getSchedules(ctx, groupID, ptr("current"))
}

func (c *Client) GetSchedulesForNextWeek(ctx context.Context, groupID int64) ([]*models.DtoScheduleResponse, error) {
	return c.getSchedules(ctx, groupID, ptr("next"))
}

func (c *Client) getSchedules(ctx context.Context, groupID int64, week *string) ([]*models.DtoScheduleResponse, error) {
	resp, err := c.c.Schedules.GetGroupsGroupIDSchedulesContext(ctx,
		schedules.NewGetGroupsGroupIDSchedulesParams().WithGroupID(groupID).WithWeek(week))
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func (c *Client) getSchedule(ctx context.Context, groupID int64,
	date *time.Time, weekday *time.Weekday, week, day *string) (*models.DtoScheduleResponse, error) {
	var dateString *string
	if date != nil {
		dateString = ptr(date.Format("02-01-2006"))
	}
	var w *string
	if weekday != nil {
		w = ptr(weekday.String())
	}
	resp, err := c.c.Schedules.GetGroupsGroupIDSchedulesContext(ctx,
		schedules.NewGetGroupsGroupIDSchedulesParams().WithGroupID(groupID).
			WithDate(dateString).WithDay(day).WithWeekday(w).WithWeek(week))

	if err != nil {
		return nil, err
	}

	if len(resp.Payload) == 0 {
		return nil, ErrNotFound
	}
	return resp.Payload[0], nil
}

func ptr[T any](v T) *T {
	return &v
}
