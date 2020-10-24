package interviewStorage

import (
	"database/sql"
	sqlutils "github.com/Solar-2020/GoUtils/sql"
	"github.com/Solar-2020/Interview-Backend/pkg/models"
	sqltools "github.com/go-park-mail-ru/2019_2_Next_Level/pkg/sqlTools"
	"strconv"
	"strings"
)

type Storage interface {
	InsertInterviews(interviews []models.Interview, postID int) (err error)
	SelectInterviews(postIDs []int) (interviews []models.Interview, err error)
	SelectInterviewsWithStatus(postIDs []int, userID int) (interviews []models.InterviewResult, err error)

	RemoveInterviews(ids []models.InterviewID) (removed []models.InterviewID, err error)

	SelectInterview(interviewID models.InterviewID) (interview models.InterviewFrame, err error)
}

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return &storage{
		db: db,
	}
}

func (s *storage) InsertInterviews(interviews []models.Interview, postID int) (err error) {
	if len(interviews) == 0 {
		return
	}

	const sqlQuery = `
	INSERT INTO interviews(text, type, post_id)
	VALUES ($1, $2, $3)
	RETURNING id;`

	tx, err := s.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	for i, _ := range interviews {
		var currentInterviewID models.InterviewID
		err = s.db.QueryRow(sqlQuery, interviews[i].Text, interviews[i].Type, postID).Scan(&currentInterviewID)
		if err != nil {
			return
		}
		interviews[i].ID = currentInterviewID

		err = s.insertAnswers(tx, interviews[i].Answers, currentInterviewID)
		if err != nil {
			return
		}
	}

	err = tx.Commit()

	return
}

func (s *storage) insertAnswers(tx *sql.Tx, answers []models.Answer, interviewID models.InterviewID) (err error) {
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

func (s *storage) SelectInterviews(postIDs []int) (interviews []models.Interview, err error) {
	interviews = make([]models.Interview, 0)
	if len(postIDs) == 0 {
		return
	}
	const sqlQueryTemplate = `
	SELECT i.id, i.text, i.type, i.post_id
	FROM interviews AS i
	WHERE i.post_id IN `

	sqlQuery := sqlQueryTemplate + sqlutils.CreateIN(len(postIDs))

	var params []interface{}

	for i, _ := range postIDs {
		params = append(params, postIDs[i])
	}

	for i := 1; i <= len(postIDs)*1; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	rows, err := s.db.Query(sqlQuery, params...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tempInterview models.Interview
		err = rows.Scan(&tempInterview.ID, &tempInterview.Text, &tempInterview.Type, &tempInterview.PostID)
		if err != nil {
			return
		}
		tempInterview.Answers = make([]models.Answer, 0)
		interviews = append(interviews, tempInterview)
	}

	interviewIDs := make([]models.InterviewID, 0)
	for i, _ := range interviews {
		interviewIDs = append(interviewIDs, interviews[i].ID)
	}

	answers, err := s.selectAnswers(interviewIDs)
	if err != nil {
		return
	}

	for _, answer := range answers {
		for i, _ := range interviews {
			if answer.InterviewID == interviews[i].ID {
				interviews[i].Answers = append(interviews[i].Answers, answer)
			}
		}
	}
	return
}

func (s *storage) SelectInterviewsWithStatus(postIDs []int, userID int) (interviews []models.InterviewResult, err error) {
	interviews = make([]models.InterviewResult, 0)
	if len(postIDs) == 0 {
		return
	}
	const sqlQueryTemplate = `
	SELECT i.id,
		   i.text,
		   i.type,
		   i.post_id,
		   (SELECT count(*)
			FROM users_answers AS ua
			WHERE ua.user_id = $1
			  AND ua.interview_id = i.post_id)
	FROM interviews AS i
	WHERE i.post_id IN `

	sqlQuery := sqlQueryTemplate + sqlutils.CreateIN(len(postIDs))

	var params []interface{}

	params = append(params, userID)

	for i, _ := range postIDs {
		params = append(params, postIDs[i])
	}

	for i := 2; i <= len(postIDs)*1+1; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	rows, err := s.db.Query(sqlQuery, params...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tempInterview models.InterviewResult
		err = rows.Scan(&tempInterview.ID, &tempInterview.Text, &tempInterview.Type, &tempInterview.PostID, &tempInterview.Status)
		if err != nil {
			return
		}

		if tempInterview.Status != 0 {
			tempInterview.Status = 1
		}

		tempInterview.Answers = make([]models.AnswerResult, 0)
		interviews = append(interviews, tempInterview)
	}

	return
}

func (s *storage) RemoveInterviews(ids []models.InterviewID) (removed []models.InterviewID, err error) {
	removed = make([]models.InterviewID, 0, len(ids))
	if len(ids) == 0 {
		return
	}

	const sqlQueryTemplate = `
	DELETE FROM interviews AS i
	WHERE i.id IN `
	const sqlQueryPostfix = ` RETURNING i.id`

	sqlQuery := sqltools.CreatePacketQuery(sqlQueryTemplate, len(ids), 1, sqlQueryPostfix)

	var params []interface{}

	for i, _ := range ids {
		params = append(params, ids[i])
	}

	rows, err := s.db.Query(sqlQuery, params...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id models.InterviewID
		err = rows.Scan(&id)
		if err != nil {
			return
		}
		removed = append(removed, id)
	}
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

func (s *storage) SelectInterview(interviewID models.InterviewID) (interview models.InterviewFrame, err error) {
	const sqlQuery = `
	SELECT i.id, i.text, i.type, i.post_id
	FROM interviews AS i
	WHERE i.id = $1;`

	err = s.db.QueryRow(sqlQuery, interviewID).Scan(&interview.ID, &interview.Text, &interview.Type, &interview.PostID)
	if err != nil {
		return
	}

	return
}
