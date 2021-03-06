package interview

import (
	"encoding/json"
	"errors"
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

	GetResultDecode(ctx *fasthttp.RequestCtx) (interviewID models.InterviewID, userID int, err error)
	GetResultEncode(response models.InterviewResult, ctx *fasthttp.RequestCtx) (err error)

	GetUniversalDecode(ctx *fasthttp.RequestCtx) (request api.GetUniversalRequest, err error)
	GetUniversalEncode(response api.GetUniversalResponse, ctx *fasthttp.RequestCtx) (err error)

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

func (t transport) GetResultDecode(ctx *fasthttp.RequestCtx) (request models.InterviewID, userID int, err error) {
	//userID := ctx.Value("UserID").(int)
	userID = ctx.UserValue("userID").(int)
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

func (t transport) GetUniversalDecode(ctx *fasthttp.RequestCtx) (request api.GetUniversalRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	return
}

func (t transport) GetUniversalEncode(response api.GetUniversalResponse, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)

	return
}

func (t transport) SetAnswerDecode(ctx *fasthttp.RequestCtx) (request models.UserAnswers, err error) {
	var userAnswers models.UserAnswers
	err = json.Unmarshal(ctx.Request.Body(), &userAnswers)
	if err != nil {
		return
	}

	interviewIDStr := ctx.UserValue("interviewID").(string)
	interviewID, err := strconv.Atoi(interviewIDStr)
	if err != nil {
		return
	}

	userAnswers.InterviewID = models.InterviewID(interviewID)


	userID, ok := ctx.UserValue("userID").(int)
	if ok {
		userAnswers.UserID = userID
		request = userAnswers
		return
	}
	return request, errors.New("userID not found")
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
