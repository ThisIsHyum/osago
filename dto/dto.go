package dto

import "github.com/ThisIsHyum/osago/types"

type NewParserRequest struct {
	CollegeName string   `json:"collegeName"`
	CampusNames []string `json:"campusNames"`
}

type NewParserResponse struct {
	Token string `json:"token"`
}

type GroupsRequest struct {
	CampusID          uint     `json:"campusId"`
	StudentGroupNames []string `json:"studentGroupNames"`
}

type LessonsRequest struct {
	CampusID uint `json:"campusId"`
	Lessons  []types.Lesson
}

type ParserResponse struct {
	ParserID  uint `json:"parserId"`
	CollegeID uint `json:"collegeId"`
}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}
