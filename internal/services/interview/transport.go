package interview

import (
	"encoding/json"
	"github.com/Solar-2020/Interview-Backend/internal/models"
	"github.com/valyala/fasthttp"
	"strconv"
)

type Transport interface {
	CreateDecode(ctx *fasthttp.RequestCtx) (request models.InterviewsRequest, err error)
	CreateEncode(response models.InterviewsRequest, ctx *fasthttp.RequestCtx) (err error)

	GetDecode(ctx *fasthttp.RequestCtx) (interviewIDs []int, err error)
	GetEncode(response []models.Interview, ctx *fasthttp.RequestCtx) (err error)

	GetResultDecode(ctx *fasthttp.RequestCtx) (interviewID int, err error)
	GetResultEncode(response models.InterviewResult, ctx *fasthttp.RequestCtx) (err error)

	GetResultsDecode(ctx *fasthttp.RequestCtx) (interviewIDs []int, err error)
	GetResultsEncode(response []models.InterviewResult, ctx *fasthttp.RequestCtx) (err error)

	SetAnswerDecode(ctx *fasthttp.RequestCtx) (request models.UserAnswers, err error)
	SetAnswerEncode(response models.InterviewResult, ctx *fasthttp.RequestCtx) (err error)
}

type transport struct {
}

func NewTransport() Transport {
	return &transport{}
}

func (t *transport) CreateDecode(ctx *fasthttp.RequestCtx) (request models.InterviewsRequest, err error) {
	var inputPost models.InterviewsRequest
	err = json.Unmarshal(ctx.Request.Body(), &inputPost)
	if err != nil {
		return
	}
	request = inputPost
	return
}

func (t *transport) CreateEncode(response models.InterviewsRequest, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t *transport) GetDecode(ctx *fasthttp.RequestCtx) (interviewIDs []int, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &interviewIDs)
	return
}

func (t transport) GetEncode(response []models.Interview, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		ctx.Response.Header.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) GetResultDecode(ctx *fasthttp.RequestCtx) (interviewID int, err error) {
	//userID := ctx.Value("UserID").(int)
	//userID := 1
	//TODO THINK ABOUT CHECK PERMISSION
	interviewIDStr := ctx.UserValue("interviewID").(string)
	interviewID, err = strconv.Atoi(interviewIDStr)
	return
}

func (t transport) GetResultEncode(response models.InterviewResult, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) GetResultsDecode(ctx *fasthttp.RequestCtx) (interviewIDs []int, err error) {
	panic("implement me")
}

func (t transport) GetResultsEncode(response []models.InterviewResult, ctx *fasthttp.RequestCtx) (err error) {
	panic("implement me")
}

func (t transport) SetAnswerDecode(ctx *fasthttp.RequestCtx) (request models.UserAnswers, err error) {
	//userID := ctx.Value("UserID").(int)
	userID := 1
	var userAnswers models.UserAnswers
	err = json.Unmarshal(ctx.Request.Body(), &userAnswers)
	if err != nil {
		return
	}
	userAnswers.UserID = userID
	interviewIDStr := ctx.UserValue("interviewID").(string)
	interviewID, err := strconv.Atoi(interviewIDStr)
	if err != nil {
		return
	}

	userAnswers.InterviewID = interviewID

	request = userAnswers

	return
}

func (t transport) SetAnswerEncode(response models.InterviewResult, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}
