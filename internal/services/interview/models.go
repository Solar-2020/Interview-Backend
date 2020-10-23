package interview

import (
	"github.com/Solar-2020/Interview-Backend/pkg/models"
)

type interviewStorage interface {
	InsertInterviews(interviews []models.Interview, postID int) (err error)
	SelectInterviews(postIDs []int) (interviews []models.Interview, err error)
	RemoveInterviews(ids []models.InterviewID) (removed []models.InterviewID, err error)

	SelectInterview(interviewID models.InterviewID) (interview models.InterviewFrame, err error)
}

type answerStorage interface {
	InsertUserAnswers(answers models.UserAnswers) (err error)

	SelectAnswersResult(interviewID models.InterviewID) (answers []models.AnswerResult, err error)
	SelectAnswersResults(interviewIDs []models.InterviewID) (answers []models.AnswerResult, err error)
}
