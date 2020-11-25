package handlers

import (
	httputils "github.com/Solar-2020/GoUtils/http"
	interviewHandler "github.com/Solar-2020/Interview-Backend/cmd/handlers/interview"
	"github.com/buaazp/fasthttprouter"
)

func NewFastHttpRouter(interview interviewHandler.Handler, middleware Middleware) *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.PanicHandler = httputils.PanicHandler
	router.Handle("GET", "/health", middleware.Log(httputils.HealthCheckHandler))

	router.Handle("POST", "/api/interview/result/:interviewID", middleware.Log(middleware.ExternalAuth(interview.SetAnswer)))

	router.Handle("POST", "/api/interview/create", middleware.Log(middleware.InternalAuth(interview.Create)))
	router.Handle("POST", "/api/interview/remove", middleware.Log(middleware.InternalAuth(interview.Remove)))
	router.Handle("POST", "/api/interview/list", middleware.Log(middleware.InternalAuth(interview.GetUniversal)))

	//NOT USED
	//router.Handle("GET", "/api/interview/result/:interviewID", middleware.Log(middleware.ExternalAuth(interview.GetResult)))
	//	router.Handle("POST", "/api/interview", middleware.Log(middleware.InternalAuth(interview.Get)))

	return router
}
