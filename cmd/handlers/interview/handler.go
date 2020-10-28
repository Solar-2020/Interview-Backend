package interviewHandler

import (
	"github.com/Solar-2020/GoUtils/context"
	http "github.com/Solar-2020/GoUtils/http"
	"github.com/Solar-2020/Interview-Backend/internal/services/interview"
)

type Handler interface {
	Create(ctx context.Context)
	Get(ctx context.Context)
	GetUniversal(ctx context.Context)
	Remove(ctx context.Context)

	GetResult(ctx context.Context)
	SetAnswer(ctx context.Context)
}

type handler struct {
	interviewService   interview.Service
	interviewTransport interviewTransport
	errorWorker        errorWorker
}

func NewHandler(interviewService interview.Service, interviewTransport interviewTransport, errorWorker errorWorker) Handler {
	return &handler{
		interviewService:   interviewService,
		interviewTransport: interviewTransport,
		errorWorker:        errorWorker,
	}
}

func (h *handler) Create(ctx context.Context) {
	poll, err := h.interviewTransport.CreateDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	pollReturn, err := h.interviewService.Create(poll)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.interviewTransport.CreateEncode(pollReturn, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) Get(ctx context.Context) {
	list, err := h.interviewTransport.GetDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	listReturn, err := h.interviewService.Get(list)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.interviewTransport.GetEncode(listReturn, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) Remove(ctx context.Context) {
	list, err := h.interviewTransport.RemoveDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	listReturn, err := h.interviewService.Remove(list)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = http.EncodeDefault(&listReturn, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) GetResult(ctx context.Context) {
	interviewID, userID, err := h.interviewTransport.GetResultDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	interviewResult, err := h.interviewService.GetResult(interviewID, userID)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.interviewTransport.GetResultEncode(interviewResult, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) SetAnswer(ctx context.Context) {
	userAnswers, err := h.interviewTransport.SetAnswerDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	interviewResult, err := h.interviewService.SetAnswers(userAnswers)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.interviewTransport.SetAnswerEncode(interviewResult, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) GetUniversal(ctx context.Context) {
	request, err := h.interviewTransport.GetUniversalDecode(ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	response, err := h.interviewService.GetUniversal(request)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.interviewTransport.GetUniversalEncode(response, ctx.RequestCtx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) handleError(err error, ctx context.Context) {
	err = h.errorWorker.ServeJSONError(ctx.RequestCtx, err)
	if err != nil {
		h.errorWorker.ServeFatalError(ctx.RequestCtx)
	}
	return
}