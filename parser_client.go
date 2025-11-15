package osago

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ThisIsHyum/osago/dto"
	"github.com/ThisIsHyum/osago/types"
)

type Parser interface {
	SendLessons(groups map[string]uint, lessons chan<- []types.Lesson) error
	GetStudentGroupNames(campusName string) (groupNames []string, _ error)
	GetCalls() ([]types.Call, error)
}

type ParserClient struct {
	*Client
	Parser    Parser
	token     string
	collegeId uint
}

func NewParserClient(ctx context.Context, baseURL, token string, timeout time.Duration) (*ParserClient, error) {
	if token == "" {
		return nil, errors.New("token cannot be empty")
	}

	c := NewClient(baseURL, timeout)
	c.httpClient.Transport = newAuthTransport(token)

	var parserResp dto.ParserResponse
	err := c.doReq(ctx, http.MethodGet, "/parser", nil, &parserResp)
	if err != nil {
		return nil, err
	}

	return &ParserClient{
		Client:    c,
		token:     token,
		collegeId: parserResp.CollegeID,
	}, nil
}

func (c *ParserClient) SetParser(parser Parser) { c.Parser = parser }

func (c *ParserClient) UpdateGroups(ctx context.Context, campusID uint, studentGroupNames []string) error {
	request := dto.GroupsRequest{
		CampusID:          campusID,
		StudentGroupNames: studentGroupNames,
	}
	return c.doReq(ctx, http.MethodPost, "/parser/groups", request, nil)
}
func (c *ParserClient) UpdateCalls(ctx context.Context, calls []types.Call) error {
	return c.doReq(ctx, http.MethodPost, "/parser/calls", calls, nil)
}

func (c *ParserClient) AddLessons(ctx context.Context, lessons []types.Lesson) error {
	return c.doReq(ctx, http.MethodPost, "/parser/lessons", lessons, nil)
}

func (c *ParserClient) GetCollege(ctx context.Context) (types.College, error) {
	return c.Client.GetCollege(ctx, c.collegeId)
}
func (c *ParserClient) GetCampuses(ctx context.Context) ([]types.Campus, error) {
	return c.Client.GetCampuses(ctx, c.collegeId)
}
func (c *ParserClient) GetCampusByName(ctx context.Context, name string) (types.Campus, error) {
	return c.Client.GetCampusByName(ctx, c.collegeId, name)
}

func (c *ParserClient) GetGroups(ctx context.Context) ([]types.StudentGroup, error) {
	return c.Client.GetGroupsByCollegeID(ctx, c.collegeId)
}

func (c *ParserClient) Run(ctx context.Context) error {
	if c.Parser == nil {
		return errors.New("parser is not set")
	}

	college, err := c.GetCollege(ctx)
	if err != nil {
		return fmt.Errorf("unable to get college: %w", err)
	}

	for _, campus := range college.Campuses {
		groupNames, err := c.Parser.GetStudentGroupNames(campus.Name)
		if err != nil {
			return fmt.Errorf("unable to get student group names: %w", err)
		}
		if err := c.UpdateGroups(ctx, campus.CampusID, groupNames); err != nil {
			return fmt.Errorf("unable to update groups: %w", err)
		}
	}

	groups, err := c.GetGroups(ctx)
	if err != nil {
		return fmt.Errorf("unable to get groups: %w", err)
	}

	mapGroups := make(map[string]uint, len(groups))
	for _, group := range groups {
		mapGroups[group.Name] = group.StudentGroupID
	}

	calls, err := c.Parser.GetCalls()
	if err != nil {
		return fmt.Errorf("unable to get calls: %w", err)
	}
	if err = c.UpdateCalls(ctx, calls); err != nil {
		return fmt.Errorf("unable to update calls: %w", err)
	}
	errChan := make(chan error, 1)
	lessonsChan := make(chan []types.Lesson)
	go func() {
		defer close(lessonsChan)
		errChan <- c.Parser.SendLessons(mapGroups, lessonsChan)
	}()
	for lessons := range lessonsChan {
		err := c.AddLessons(ctx, lessons)
		if err != nil {
			return fmt.Errorf("unable to add lessons: %w", err)
		}
	}
	return <-errChan
}
