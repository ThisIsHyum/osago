package types

import (
	"fmt"
	"strings"
	"time"
)

type (
	College struct {
		CollegeID uint     `json:"collegeId"`
		Name      string   `json:"name"`
		Calls     []Call   `json:"calls"`
		Campuses  []Campus `json:"campuses"`
	}

	Campus struct {
		CampusID      uint           `json:"campusId"`
		Name          string         `json:"name"`
		CollegeID     uint           `json:"collegeId"`
		StudentGroups []StudentGroup `json:"studentGroups"`
	}

	StudentGroup struct {
		StudentGroupID uint   `json:"studentGroupId,omitempty"`
		Name           string `json:"name"`
		CampusID       uint   `json:"campusId"`
	}

	Lesson struct {
		LessonID       uint         `json:"lessonId"`
		Title          string       `json:"title"`
		Cabinet        string       `json:"cabinet"`
		Date           time.Time    `json:"date"`
		Teacher        string       `json:"teacher"`
		Order          uint         `json:"order"`
		StudentGroupID uint         `json:"studentGroupID"`
		StudentGroup   StudentGroup `json:"studentGroup,omitempty"`
	}

	Call struct {
		CallID    uint         `json:"callId"`
		Weekday   time.Weekday `json:"weekday"`
		Begins    CallTime     `json:"begins"`
		Ends      CallTime     `json:"ends"`
		Order     uint         `json:"order"`
		CollegeID uint         `json:"-"`
	}
	CallTime struct {
		time.Time
	}

	Parser struct {
		ParserID  uint
		Token     string
		CollegeID uint
		College   College
	}
	Admin struct {
		AdminID uint
		Token   string
	}

	Schedule struct {
		GroupID uint             `json:"groupId"`
		Date    time.Time        `json:"date"`
		Lessons []ScheduleLesson `json:"lessons"`
	}

	ScheduleLesson struct {
		Title     string   `json:"title"`
		Cabinet   string   `json:"cabinet"`
		Teacher   string   `json:"teacher"`
		Order     uint     `json:"order"`
		StartTime CallTime `json:"startTime"`
		EndTime   CallTime `json:"endTime"`
	}
	Weekdays []time.Weekday
)

func (c College) String() string      { return c.Name }
func (c Campus) String() string       { return c.Name }
func (g StudentGroup) String() string { return g.Name }

func NewGroup(name string, campusID uint) StudentGroup {
	return StudentGroup{
		Name:     name,
		CampusID: campusID,
	}
}

func NewCalls(weekdays []time.Weekday, order uint, begins, ends time.Time) []Call {
	calls := []Call{}
	for _, weekday := range weekdays {
		calls = append(calls, Call{
			Weekday: weekday,
			Begins:  CallTime{begins},
			Ends:    CallTime{ends},
			Order:   order,
		})
	}
	return calls
}

func (c *CallTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	if s != "null" {
		c.Time, err = time.Parse("15:04:05", s)
	}
	return
}
func (c CallTime) MarshalJSON() ([]byte, error) {
	if c.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, c.Time.Format("15:04"))), nil
}
func NewLesson(title, cabinet, teacher string, date time.Time, groupID, order uint) Lesson {
	return Lesson{
		Title:          title,
		Cabinet:        cabinet,
		Teacher:        teacher,
		Date:           date,
		Order:          order,
		StudentGroupID: groupID,
	}
}
