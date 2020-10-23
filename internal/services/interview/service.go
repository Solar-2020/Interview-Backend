package interview

import (
	"database/sql"
	"errors"
	"github.com/Solar-2020/Interview-Backend/pkg/api"
	"github.com/Solar-2020/Interview-Backend/pkg/models"
)

type Service interface {
	Create(request api.CreateRequest) (response api.CreateResponse, err error)
	Get(postIds api.GetRequest) (response api.GetResponse, err error)
	Remove(interviewIds api.RemoveRequest) (response api.RemoveRequest, err error)

	GetResult(interviewID models.InterviewID) (response models.InterviewResult, err error)
	GetResults(interviewIDs []models.InterviewID) (response []models.InterviewResult, err error)

	SetAnswers(answers models.UserAnswers) (response models.InterviewResult, err error)
}

type service struct {
	interviewStorage interviewStorage
	answerStorage    answerStorage
}

func NewService(interviewStorage interviewStorage, answerStorage answerStorage) Service {
	return &service{
		interviewStorage: interviewStorage,
		answerStorage:    answerStorage,
	}
}

// POST /models/interview/create
func (s *service) Create(request api.CreateRequest) (response api.CreateResponse, err error) {
	err = s.interviewStorage.InsertInterviews(request.Interviews, request.PostID)
	if err != nil {
		return
	}
	response = api.CreateResponse{Interviews: request.Interviews}
	for i := range request.Interviews {
		response.Interviews[i].PostID = request.PostID
	}
	return
}

// POST /interview
func (s *service) Get(postIds api.GetRequest) (response api.GetResponse, err error) {
	resp, err := s.interviewStorage.SelectInterviews(postIds.Ids)
	if err != nil {
		return
	}
	response.Interviews = resp
	return
}

// POST /interview/remove
func (s *service) Remove(interviewIds api.RemoveRequest) (response api.RemoveRequest, err error) {
	ids, err := s.interviewStorage.RemoveInterviews(interviewIds.Ids)
	if err == nil {
		response.Ids = ids
	}
	return
}

func (s *service) GetResult(interviewID models.InterviewID) (response models.InterviewResult, err error) {
	response.InterviewFrame, err = s.interviewStorage.SelectInterview(interviewID)
	if err != nil {
		if err == sql.ErrNoRows {
			return response, errors.New("Опрос не найден")
		}
		return
	}

	response.Answers, err = s.answerStorage.SelectAnswersResult(interviewID)

	return
}

func (s *service) GetResults(interviewIDs []models.InterviewID) (response []models.InterviewResult, err error) {
	panic("implement me")
}

func (s *service) SetAnswers(answers models.UserAnswers) (response models.InterviewResult, err error) {
	//TODO CHECK PERMISSION
	//TODO CHECK REPEATED VOTE
	err = s.answerStorage.InsertUserAnswers(answers)
	if err != nil {
		return
	}

	response, err = s.GetResult(answers.InterviewID)

	return
}
