package interview

import (
	"github.com/Solar-2020/Interview-Backend/internal/models"
)

type Service interface {
	Create(request models.Interview) (response models.Interview, err error)

	Get(interviewIDs []int) (response []models.Interview, err error)

	GetResult(interviewID int) (response models.InterviewResult, err error)
	GetResults(interviewIDs []int) (response []models.InterviewResult, err error)

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

func (s *service) Create(request models.Interview) (response models.Interview, err error) {
	panic("implement me")
}

func (s *service) Get(interviewIDs []int) (response []models.Interview, err error) {
	panic("implement me")
}

func (s *service) GetResult(interviewID int) (response models.InterviewResult, err error) {
	response.InterviewFrame, err = s.interviewStorage.SelectInterview(interviewID)
	if err != nil {
		return
	}

	response.Answers, err = s.answerStorage.SelectAnswersResult(interviewID)

	return
}

func (s *service) GetResults(interviewIDs []int) (response []models.InterviewResult, err error) {
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
