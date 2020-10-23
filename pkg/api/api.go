package api

import "github.com/Solar-2020/Interview-Backend/internal/models"

// POST /interview/create
type CreateRequest struct {
	Interviews []models.Interview `json:"interviews"`
	PostID int `json:"postID"`
}

type CreateResponse struct {
	models.Answer
}

// POST /interview
type GetRequest struct {
	Ids []int	`json:"posts"`
}

type GetResponse struct {
	Interviews map[int][]models.Interview
}

// GET /interview/result/:id
type ResultRequest struct {
	Id models.InterviewID `json:"id"`
}

type ResultResponse struct {
	models.InterviewResult
}

// POST /interview/result/:id
type SetVoteRequest struct {
	models.UserAnswer
}

type SetVoteResponse struct {
	models.InterviewResult
}
