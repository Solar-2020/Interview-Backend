package interview

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Solar-2020/Interview-Backend/pkg/api"
	"github.com/Solar-2020/Interview-Backend/pkg/models"
)

type Service interface {
	Create(request api.CreateRequest) (response api.CreateResponse, err error)
	Get(postIds api.GetRequest) (response api.GetResponse, err error)
	Remove(interviewIds api.RemoveRequest) (response api.RemoveRequest, err error)

	GetResult(interviewID models.InterviewID, userID int) (response models.InterviewResult, err error)
	GetUniversal(request api.GetUniversalRequest) (response api.GetUniversalResponse, err error)

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

func (s *service) GetResult(interviewID models.InterviewID, userID int) (response models.InterviewResult, err error) {
	response.InterviewFrame, err = s.interviewStorage.SelectInterviewWithStatus(interviewID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return response, errors.New("Опрос не найден")
		}
		return
	}

	response.Answers, err = s.answerStorage.SelectAnswersResult(interviewID)
	if err != nil {
		return
	}

	userAnswers, err := s.answerStorage.SelectUserAnswer(response.ID, userID)
	if err != nil {
		return
	}

	for i := range response.Answers {
		for _, ans := range userAnswers {
			if int(response.Answers[i].ID) == int(ans.ID) {
				response.Answers[i].IsMyAnswer = true
			}
		}
	}

	return
}

func (s *service) GetUniversal(request api.GetUniversalRequest) (response api.GetUniversalResponse, err error) {
	interviews, err := s.interviewStorage.SelectInterviewsWithStatus(request.PostIDs, request.UserID)
	if err != nil {
		return
	}

	interviewIDs := make([]models.InterviewID, 0)
	for i, _ := range interviews {
		interviewIDs = append(interviewIDs, interviews[i].ID)
	}

	answers, err := s.answerStorage.SelectAnswersResults(interviewIDs)
	if err != nil {
		return
	}

	userAnswers, err := s.answerStorage.SelectUserAnswers(interviewIDs, request.UserID)
	if err != nil {
		return
	}

	for i, _ := range answers {
		for j, _ := range userAnswers {
			if int(answers[i].ID) == int(userAnswers[j].ID) {
				answers[i].IsMyAnswer = true
			}
		}
	}

	for i, _ := range interviews {
		for j, _ := range answers {
			if interviews[i].ID == answers[j].InterviewID {
				interviews[i].Answers = append(interviews[i].Answers, answers[j])
			}
		}
	}

	response.Interviews = interviews
	return
}

func (s *service) SetAnswers(answers models.UserAnswers) (response models.InterviewResult, err error) {
	//TODO CHECK PERMISSION
	//TODO CHECK REPEATED VOTE
	err = s.answerStorage.InsertUserAnswers(answers)
	if err != nil {
		fmt.Println(err)
		err = fmt.Errorf("cannot set vote")
		return
	}

	response, err = s.GetResult(answers.InterviewID, answers.UserID)

	return
}
