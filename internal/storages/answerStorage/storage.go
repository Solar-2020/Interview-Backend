package answerStorage

import (
	"database/sql"
	sqlutils "github.com/Solar-2020/GoUtils/sql"
	"github.com/Solar-2020/Interview-Backend/pkg/models"
	"strconv"
	"strings"
)

type Storage interface {
	InsertUserAnswers(answers models.UserAnswers) (err error)

	SelectAnswersResult(interviewID models.InterviewID) (answers []models.AnswerResult, err error)
	SelectAnswersResults(interviewIDs []models.InterviewID) (answers []models.AnswerResult, err error)

	SelectUserAnswer(interviewIDs models.InterviewID, userID int) (answers []models.UserAnswer, err error)
	SelectUserAnswers(interviewIDs []models.InterviewID, userID int) (answers []models.UserAnswer, err error)
}

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return &storage{
		db: db,
	}
}

func (s *storage) insertAnswers(tx *sql.Tx, answers []models.Answer, interviewID int) (err error) {
	if len(answers) == 0 {
		return
	}

	sqlQueryTemplate := `
	INSERT INTO answers(interview_id, text)
	VALUES `

	if len(answers) == 0 {
		return
	}

	var params []interface{}

	sqlQuery := sqlQueryTemplate + sqlutils.CreateInsertQuery(len(answers), 2)

	for i, _ := range answers {
		params = append(params, interviewID, answers[i].Text)
	}

	for i := 1; i <= len(answers)*2; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	_, err = tx.Exec(sqlQuery, params...)
	return
}

func (s *storage) selectAnswers(interviewIDs []models.InterviewID) (answers []models.Answer, err error) {
	const sqlQueryTemplate = `
	SELECT a.id, a.text, a.interview_id
	FROM answers AS a
	WHERE a.interview_id IN `

	sqlQuery := sqlQueryTemplate + sqlutils.CreateIN(len(interviewIDs))

	var params []interface{}

	for i, _ := range interviewIDs {
		params = append(params, interviewIDs[i])
	}

	for i := 1; i <= len(interviewIDs)*1; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	rows, err := s.db.Query(sqlQuery, params...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tempAnswer models.Answer
		err = rows.Scan(&tempAnswer.ID, &tempAnswer.Text, &tempAnswer.InterviewID)
		if err != nil {
			return
		}
		answers = append(answers, tempAnswer)
	}

	return
}

func (s *storage) InsertUserAnswers(answers models.UserAnswers) (err error) {
	if len(answers.AnswerIDs) == 0 {
		return
	}

	sqlQueryTemplate := `
	INSERT INTO users_answers(interview_id, user_id, answer_id, post_id)
	VALUES `

	var params []interface{}

	sqlQuery := sqlQueryTemplate + sqlutils.CreateInsertQuery(len(answers.AnswerIDs), 4)

	for i, _ := range answers.AnswerIDs {
		params = append(params, answers.InterviewID, answers.UserID, answers.AnswerIDs[i], answers.PostID)
	}

	for i := 1; i <= len(answers.AnswerIDs)*4; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	_, err = s.db.Exec(sqlQuery, params...)

	return
}

func (s *storage) SelectAnswersResults(interviewIDs []models.InterviewID) (answers []models.AnswerResult, err error) {
	if len(interviewIDs) == 0 {
		return
	}
	const sqlQueryTemplate = `
	SELECT a.id,
		   a.text,
		   a.interview_id,
		   (SELECT count(*)
			FROM users_answers AS ua
			WHERE ua.answer_id = a.id)
	FROM answers AS a
	WHERE a.interview_id IN `

	sqlQuery := sqlQueryTemplate + sqlutils.CreateIN(len(interviewIDs))

	var params []interface{}

	for i, _ := range interviewIDs {
		params = append(params, interviewIDs[i])
	}

	for i := 1; i <= len(interviewIDs)*1; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	rows, err := s.db.Query(sqlQuery, params...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tempAnswer models.AnswerResult
		err = rows.Scan(&tempAnswer.ID, &tempAnswer.Text, &tempAnswer.InterviewID, &tempAnswer.AnswerCount)
		if err != nil {
			return
		}
		answers = append(answers, tempAnswer)
	}

	return
}

func (s *storage) SelectAnswersResult(interviewID models.InterviewID) (answers []models.AnswerResult, err error) {
	interviewIDs := make([]models.InterviewID, 1)
	interviewIDs = append(interviewIDs, interviewID)
	answers, err = s.SelectAnswersResults(interviewIDs)
	return
}

func (s *storage) SelectUserAnswers(interviewIDs []models.InterviewID, userID int) (answers []models.UserAnswer, err error) {
	if len(interviewIDs) == 0 {
		return
	}
	const sqlQueryTemplate = `
	SELECT ua.answer_id,
		   ua.interview_id
	FROM users_answers AS ua
	WHERE ua.user_id = $1 AND ua.interview_id IN `

	sqlQuery := sqlQueryTemplate + sqlutils.CreateIN(len(interviewIDs))

	var params []interface{}

	params = append(params, userID)

	for i, _ := range interviewIDs {
		params = append(params, interviewIDs[i])
	}

	for i := 2; i <= len(interviewIDs)*1+1; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	rows, err := s.db.Query(sqlQuery, params...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tempAnswer models.UserAnswer
		err = rows.Scan(&tempAnswer.ID, &tempAnswer.InterviewID)
		if err != nil {
			return
		}
		answers = append(answers, tempAnswer)
	}

	return
}

func (s *storage) SelectUserAnswer(interviewID models.InterviewID, userID int) (answers []models.UserAnswer, err error) {
	interviewIDs := make([]models.InterviewID, 1)
	interviewIDs = append(interviewIDs, interviewID)
	answers, err = s.SelectUserAnswers(interviewIDs, userID)
	return
}
