package interview

import (
	"encoding/json"
	"github.com/Solar-2020/Interview-Backend/pkg/api"
	"github.com/Solar-2020/Interview-Backend/pkg/models"
	"github.com/go-playground/validator"
	"github.com/valyala/fasthttp"
	"strconv"
)

type Transport interface {
	CreateDecode(ctx *fasthttp.RequestCtx) (request api.CreateRequest, err error)
	CreateEncode(response api.CreateResponse, ctx *fasthttp.RequestCtx) (err error)

	GetDecode(ctx *fasthttp.RequestCtx) (request api.GetRequest, err error)
	GetEncode(response api.GetResponse, ctx *fasthttp.RequestCtx) (err error)

	RemoveDecode(ctx *fasthttp.RequestCtx) (request api.RemoveRequest, err error)
	//RemoveEncode(response models.RemoveResponse, ctx *fasthttp.RequestCtx) (err error)

	GetResultDecode(ctx *fasthttp.RequestCtx) (interviewID models.InterviewID, err error)
	GetResultEncode(response models.InterviewResult, ctx *fasthttp.RequestCtx) (err error)

	GetResultsDecode(ctx *fasthttp.RequestCtx) (interviewIDs []models.InterviewID, err error)
	GetResultsEncode(response []models.InterviewResult, ctx *fasthttp.RequestCtx) (err error)

	SetAnswerDecode(ctx *fasthttp.RequestCtx) (request models.UserAnswers, err error)
	SetAnswerEncode(response models.InterviewResult, ctx *fasthttp.RequestCtx) (err error)
}

type transport struct {
	validator *validator.Validate
}

func NewTransport() Transport {
	return &transport{
		validator: validator.New(),
	}
}

func (t *transport) CreateDecode(ctx *fasthttp.RequestCtx) (request api.CreateRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	return
}

func (t *transport) CreateEncode(response api.CreateResponse, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t *transport) GetDecode(ctx *fasthttp.RequestCtx) (request api.GetRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	return
}

func (t transport) GetEncode(response api.GetResponse, ctx *fasthttp.RequestCtx) (err error) {
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

func (t transport) RemoveDecode(ctx *fasthttp.RequestCtx) (request api.RemoveRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	return
}

func (t transport) GetResultDecode(ctx *fasthttp.RequestCtx) (request models.InterviewID, err error) {
	//userID := ctx.Value("UserID").(int)
	//userID := 1
	//TODO THINK ABOUT CHECK PERMISSION
	interviewIDStr := ctx.UserValue("interviewID").(string)
	tmp, err := strconv.Atoi(interviewIDStr)
	if err != nil {
		return
	}
	request = models.InterviewID(tmp)
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

func (t transport) GetResultsDecode(ctx *fasthttp.RequestCtx) (request []models.InterviewID, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	return
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

	userAnswers.InterviewID = models.InterviewID(interviewID)

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
