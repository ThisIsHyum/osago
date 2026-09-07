package osago

import (
	"context"
	"time"

	"github.com/ThisIsHyum/osago/client/schedules"
	"github.com/ThisIsHyum/osago/models"
)

type ScheduleFilter struct {
	Teacher, Cabinet, Title *string
}

func (c *Client) GetScheduleForToday(ctx context.Context, groupID int64, filter ScheduleFilter) (*models.DtoScheduleResponse, error) {
	return c.getSchedule(ctx, groupID, nil, nil, nil, ptr("today"), filter)
}

func (c *Client) GetScheduleForTomorrow(ctx context.Context, groupID int64, filter ScheduleFilter) (*models.DtoScheduleResponse, error) {
	return c.getSchedule(ctx, groupID, nil, nil, nil, ptr("tomorrow"), filter)
}

func (c *Client) GetScheduleForDate(ctx context.Context, groupID int64, date time.Time, filter ScheduleFilter) (*models.DtoScheduleResponse, error) {
	return c.getSchedule(ctx, groupID, &date, nil, nil, nil, filter)
}

func (c *Client) GetScheduleForWeekdayOfPreviousWeek(ctx context.Context, groupID int64, weekday time.Weekday, filter ScheduleFilter) (*models.DtoScheduleResponse, error) {
	return c.getSchedule(ctx, groupID, nil, &weekday, ptr("previous"), nil, filter)
}

func (c *Client) GetScheduleForWeekday(ctx context.Context, groupID int64, weekday time.Weekday, filter ScheduleFilter) (*models.DtoScheduleResponse, error) {
	return c.getSchedule(ctx, groupID, nil, &weekday, ptr("current"), nil, filter)
}

func (c *Client) GetScheduleForWeekdayOfNextWeek(ctx context.Context, groupID int64, weekday time.Weekday, filter ScheduleFilter) (*models.DtoScheduleResponse, error) {
	return c.getSchedule(ctx, groupID, nil, &weekday, ptr("next"), nil, filter)
}

func (c *Client) GetSchedulesForPreviousWeek(ctx context.Context, groupID int64, filter ScheduleFilter) ([]*models.DtoScheduleResponse, error) {
	return c.getSchedules(ctx, groupID, ptr("previous"), filter)
}

func (c *Client) GetSchedulesForCurrentWeek(ctx context.Context, groupID int64, filter ScheduleFilter) ([]*models.DtoScheduleResponse, error) {
	return c.getSchedules(ctx, groupID, ptr("current"), filter)
}

func (c *Client) GetSchedulesForNextWeek(ctx context.Context, groupID int64, filter ScheduleFilter) ([]*models.DtoScheduleResponse, error) {
	return c.getSchedules(ctx, groupID, ptr("next"), filter)
}

func (c *Client) getSchedules(ctx context.Context, groupID int64, week *string, filter ScheduleFilter) ([]*models.DtoScheduleResponse, error) {
	resp, err := c.c.Schedules.GetGroupsGroupIDSchedulesContext(ctx,
		schedules.NewGetGroupsGroupIDSchedulesParams().
			WithGroupID(groupID).WithWeek(week).
			WithCabinet(filter.Cabinet).WithTeacher(filter.Teacher).WithTitle(filter.Title))
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func (c *Client) getSchedule(ctx context.Context, groupID int64,
	date *time.Time, weekday *time.Weekday, week, day *string, filter ScheduleFilter) (*models.DtoScheduleResponse, error) {
	var dateString *string
	if date != nil {
		dateString = ptr(date.Format(time.DateOnly))
	}
	var w *string
	if weekday != nil {
		w = ptr(weekday.String())
	}
	resp, err := c.c.Schedules.GetGroupsGroupIDSchedulesContext(ctx,
		schedules.NewGetGroupsGroupIDSchedulesParams().WithGroupID(groupID).
			WithDate(dateString).WithDay(day).WithWeekday(w).WithWeek(week).
			WithCabinet(filter.Cabinet).WithTeacher(filter.Teacher).WithTitle(filter.Title))

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
