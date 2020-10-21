package interviewHandler

import (
	"github.com/Solar-2020/Interview-Backend/internal/models"
	"github.com/valyala/fasthttp"
)

type interviewService interface {
	Create(request models.Interview) (response models.Interview, err error)

	Get(interviewIDs []int) (response []models.Interview, err error)

	GetResult(interviewID int) (response models.InterviewResult, err error)
	GetResults(interviewIDs []int) (response []models.InterviewResult, err error)

	SetAnswers(answers models.UserAnswers) (response models.InterviewResult, err error)
	//GetList(request models.GetPostListRequest) (response []models.Post, err error)
}

type interviewTransport interface {
	CreateDecode(ctx *fasthttp.RequestCtx) (request models.Interview, err error)
	CreateEncode(response models.Interview, ctx *fasthttp.RequestCtx) (err error)

	GetDecode(ctx *fasthttp.RequestCtx) (interviewIDs []int, err error)
	GetEncode(response []models.Interview, ctx *fasthttp.RequestCtx) (err error)

	GetResultDecode(ctx *fasthttp.RequestCtx) (interviewID int, err error)
	GetResultEncode(response models.InterviewResult, ctx *fasthttp.RequestCtx) (err error)

	GetResultsDecode(ctx *fasthttp.RequestCtx) (interviewIDs []int, err error)
	GetResultsEncode(response []models.InterviewResult, ctx *fasthttp.RequestCtx) (err error)

	SetAnswerDecode(ctx *fasthttp.RequestCtx) (request models.UserAnswers, err error)
	SetAnswerEncode(response models.InterviewResult, ctx *fasthttp.RequestCtx) (err error)

	//GetListDecode(ctx *fasthttp.RequestCtx) (request models.GetPostListRequest, err error)
	//GetListEncode(response []models.Post, ctx *fasthttp.RequestCtx) (err error)
}

type errorWorker interface {
	ServeJSONError(ctx *fasthttp.RequestCtx, serveError error) (err error)
	ServeFatalError(ctx *fasthttp.RequestCtx)
}
