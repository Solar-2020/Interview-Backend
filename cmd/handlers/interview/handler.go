package interviewHandler

import (
	http "github.com/Solar-2020/GoUtils/http"
	"github.com/Solar-2020/Interview-Backend/internal/services/interview"
	"github.com/valyala/fasthttp"
)

type Handler interface {
	Create(ctx *fasthttp.RequestCtx)
	Get(ctx *fasthttp.RequestCtx)
	Remove(ctx *fasthttp.RequestCtx)

	GetResult(ctx *fasthttp.RequestCtx)
	SetAnswer(ctx *fasthttp.RequestCtx)
}

type handler struct {
	interviewService   interview.Service
	interviewTransport interviewTransport
	errorWorker   errorWorker
}

func NewHandler(interviewService interview.Service, interviewTransport interviewTransport, errorWorker errorWorker) Handler {
	return &handler{
		interviewService:   interviewService,
		interviewTransport: interviewTransport,
		errorWorker:   errorWorker,
	}
}

func (h *handler) Create(ctx *fasthttp.RequestCtx) {
	poll, err := h.interviewTransport.CreateDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	pollReturn, err := h.interviewService.Create(poll)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	err = h.interviewTransport.CreateEncode(pollReturn, ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}
}

func (h *handler) Get(ctx *fasthttp.RequestCtx) {
	list, err := h.interviewTransport.GetDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	listReturn, err := h.interviewService.Get(list)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	err = h.interviewTransport.GetEncode(listReturn, ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}
}

func (h *handler) Remove(ctx *fasthttp.RequestCtx) {
	list, err := h.interviewTransport.RemoveDecode(ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	listReturn, err := h.interviewService.Remove(list)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	err = http.EncodeDefault(&listReturn, ctx)
	if err != nil {
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}
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

