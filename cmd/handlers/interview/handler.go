package interviewHandler

import (
	"github.com/valyala/fasthttp"
)

type Handler interface {
	Create(ctx *fasthttp.RequestCtx)
	Get(ctx *fasthttp.RequestCtx)
	GetResult(ctx *fasthttp.RequestCtx)
	SetAnswer(ctx *fasthttp.RequestCtx)
}

type handler struct {
	interviewService   interviewService
	interviewTransport interviewTransport
	errorWorker   errorWorker
}

func NewHandler(interviewService interviewService, interviewTransport interviewTransport, errorWorker errorWorker) Handler {
	return &handler{
		interviewService:   interviewService,
		interviewTransport: interviewTransport,
		errorWorker:   errorWorker,
	}
}

func (h *handler) Create(ctx *fasthttp.RequestCtx) {
	post, err := h.interviewTransport.CreateDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	postReturn, err := h.interviewService.Create(post)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	err = h.interviewTransport.CreateEncode(postReturn, ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}
}

func (h *handler) Get(ctx *fasthttp.RequestCtx) {
	panic("implement me")
}

func (h *handler) GetResult(ctx *fasthttp.RequestCtx) {
	interviewID, err := h.interviewTransport.GetResultDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	interviewResult, err := h.interviewService.GetResult(interviewID)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	err = h.interviewTransport.GetResultEncode(interviewResult, ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}
}

func (h *handler) SetAnswer(ctx *fasthttp.RequestCtx) {
	userAnswers, err := h.interviewTransport.SetAnswerDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	interviewResult, err := h.interviewService.SetAnswers(userAnswers)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	err = h.interviewTransport.SetAnswerEncode(interviewResult, ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

}

