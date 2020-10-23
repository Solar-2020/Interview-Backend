package handlers

import (
	httputils "github.com/Solar-2020/GoUtils/http"
	interviewHandler "github.com/Solar-2020/Interview-Backend/cmd/handlers/interview"
	"github.com/buaazp/fasthttprouter"
)

func NewFastHttpRouter(interview interviewHandler.Handler, middleware httputils.Middleware) *fasthttprouter.Router {
	router := fasthttprouter.New()

	//router.Handle("GET", "/health", check)

	router.PanicHandler = httputils.PanicHandler

	middlewareChain := httputils.NewLogCorsChain(middleware)

	router.Handle("POST", "/interview/create", middlewareChain(interview.Create))
	router.Handle("POST", "/interview", middlewareChain(interview.Get))
	router.Handle("POST", "/interview/remove", middlewareChain(interview.Remove))

	router.Handle("POST", "/interview/result/:interviewID", middlewareChain(interview.SetAnswer))
	router.Handle("GET", "/interview/result/:interviewID", middlewareChain(interview.GetResult))

	//router.Handle("POST", "/interview/interview", middleware.CORS(interview.Create))

	return router
}

