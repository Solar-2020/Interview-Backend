package models

type InterviewFrame struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Type   int    `json:"type"`
	PostID int    `json:"postID"`
	Status int    `json:"status"` //Проголосовал юзер или нет
}

type Interview struct {
	InterviewFrame
	Answers []Answer `json:"answers"`
}

type Answer struct {
	ID          int    `json:"id"`
	Text        string `json:"text"`
	InterviewID int    `json:"interviewID"`
}

type InterviewResult struct {
	InterviewFrame
	Answers []AnswerResult `json:"answers"`
}

type AnswerResult struct {
	Answer
	AnswerCount int `json:"answerCount"`
}

type UserAnswer struct {
	ID          int `json:"id"`
	InterviewID int `json:"interviewID"`
}

type UserAnswers struct {
	PostID      int   `json:"postID"`
	InterviewID int   `json:"interviewID"`
	UserID      int   `json:"-"`
	AnswerIDs   []int `json:"answers"`
}
