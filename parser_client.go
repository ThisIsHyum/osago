package osago

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ThisIsHyum/osago/client/parser"
	"github.com/ThisIsHyum/osago/models"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
)

type Parser interface {
	SendLessons(groups map[string]int64, lessons chan<- []*models.DtoLesson) error
	SendReplaces(groups map[string]int64, replaces chan<- []*models.DtoReplace) error
	GetStudentGroupNames(campusName string) (groupNames []string, _ error)
	GetCalls() ([]*models.DtoCall, error)
}

type ParserClient struct {
	*Client
	Parser    Parser
	auth      runtime.ClientAuthInfoWriter
	collegeID int64
}

func NewParserClient(ctx context.Context, url, token string, timeout time.Duration) (*ParserClient, error) {
	c := &ParserClient{
		Client: NewClient(url, timeout),
		auth:   httptransport.BearerToken(token),
	}

	resp, err := c.c.Parser.GetParserContext(ctx, parser.NewGetParserParams(), c.auth)
	if err != nil {
		return nil, err
	}
	c.collegeID = resp.Payload.CollegeID
	return c, nil
}

func (c *ParserClient) SetParser(parser Parser) { c.Parser = parser }

func (c *ParserClient) UpdateGroups(ctx context.Context, campusID int64, studentGroupNames []string) error {
	req := models.DtoUpdateGroupsRequest{
		CampusID:          campusID,
		StudentGroupNames: studentGroupNames,
	}
	_, err := c.c.Parser.PostParserGroupsContext(ctx,
		parser.NewPostParserGroupsParams().WithUpdateGroupsRequest(&req), c.auth)
	return err
}
func (c *ParserClient) UpdateCalls(ctx context.Context, calls []*models.DtoCall) error {
	_, err := c.c.Parser.PostParserCallsContext(ctx,
		parser.NewPostParserCallsParams().WithCalls(calls), c.auth)
	return err
}

func (c *ParserClient) AddLessons(ctx context.Context, lessons []*models.DtoLesson) error {
	_, err := c.c.Parser.PostParserLessonsContext(ctx,
		parser.NewPostParserLessonsParams().WithLessons(lessons), c.auth)
	return err
}

func (c *ParserClient) AddReplaces(ctx context.Context, replaces []*models.DtoReplace) error {
	_, err := c.c.Parser.PostParserReplacesContext(ctx,
		parser.NewPostParserReplacesParams().WithReplaces(replaces), c.auth)
	return err
}

func (c *ParserClient) GetCollege(ctx context.Context) (*models.DtoCollegeResponse, error) {
	return c.Client.GetCollege(ctx, c.collegeID)
}
func (c *ParserClient) GetCampuses(ctx context.Context) ([]*models.DtoCampusResponse, error) {
	return c.Client.GetCampuses(ctx, c.collegeID)
}
func (c *ParserClient) GetCampusByName(ctx context.Context, name string) (*models.DtoCampusResponse, error) {
	return c.Client.GetCampusByName(ctx, c.collegeID, name)
}

func (c *ParserClient) GetGroups(ctx context.Context) ([]*models.DtoStudentGroupResponse, error) {
	return c.Client.GetGroupsByCollegeID(ctx, c.collegeID, nil)
}
func (c *ParserClient) GetGroupsByName(ctx context.Context, name string) ([]*models.DtoStudentGroupResponse, error) {
	return c.Client.GetGroupsByCollegeID(ctx, c.collegeID, &name)
}

func (c *ParserClient) Run(ctx context.Context, errorCh chan<- error) error {
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

	mapGroups := make(map[string]int64, len(groups))
	for _, group := range groups {
		mapGroups[group.Name] = group.StudentGroupID
	}

	calls, err := c.Parser.GetCalls()
	if err != nil {
		return fmt.Errorf("unable to get calls: %w", err)
	}
	if err := c.UpdateCalls(ctx, calls); err != nil {
		return fmt.Errorf("unable to update calls: %w", err)
	}

	lessonsChan := make(chan []*models.DtoLesson)
	go func() {
		defer close(lessonsChan)
		if err := c.Parser.SendLessons(mapGroups, lessonsChan); err != nil {
			reportError(errorCh, fmt.Errorf("unable to send lessons: %w", err))
		}
	}()

	replacesChan := make(chan []*models.DtoReplace)
	go func() {
		defer close(replacesChan)
		if err := c.Parser.SendReplaces(mapGroups, replacesChan); err != nil {
			reportError(errorCh, fmt.Errorf("unable to send replaces: %w", err))
		}
	}()

	go func() {
		for lessons := range lessonsChan {
			if err := c.AddLessons(ctx, lessons); err != nil {
				reportError(errorCh, fmt.Errorf("unable to add lessons: %w", err))
			}
		}
	}()
	go func() {
		for replaces := range replacesChan {
			if err := c.AddReplaces(ctx, replaces); err != nil {
				reportError(errorCh, fmt.Errorf("unable to add replaces: %w", err))
			}
		}
	}()

	<-ctx.Done()
	return ctx.Err()
}

func reportError(errs chan<- error, err error) {
	if errs != nil {
		errs <- err
	}
}
