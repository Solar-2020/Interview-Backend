package interview

import (
	"github.com/Solar-2020/Interview-Backend/internal/models"
)

type interviewStorage interface {
	InsertInterviews(interviews []models.Interview, postID int) (err error)
	SelectInterviews(postIDs []int) (interviews []models.Interview, err error)

	SelectInterview(interviewID int) (interview models.InterviewFrame, err error)
}

type answerStorage interface {
	InsertUserAnswers(answers models.UserAnswers) (err error)

	SelectAnswersResult(interviewID int) (answers []models.AnswerResult, err error)
	SelectAnswersResults(interviewIDs []int) (answers []models.AnswerResult, err error)
}
