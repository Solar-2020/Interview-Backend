package handlers

import (
	httputils "github.com/Solar-2020/GoUtils/http"
	interviewHandler "github.com/Solar-2020/Interview-Backend/cmd/handlers/interview"
	"github.com/buaazp/fasthttprouter"
)

func NewFastHttpRouter(interview interviewHandler.Handler, middleware httputils.Middleware) *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.PanicHandler = httputils.PanicHandler
	router.Handle("GET", "/health", middleware.Log(httputils.HealthCheckHandler))

	clientside := httputils.ClientsideChain(middleware)
	serverside := httputils.ServersideChain(middleware)

	router.Handle("POST", "/interview/create", serverside(interview.Create))
	router.Handle("POST", "/interview", serverside(interview.Get))
	router.Handle("POST", "/interview/remove", clientside(interview.Remove))

	router.Handle("POST", "/interview/result/:interviewID",clientside(interview.SetAnswer))
	router.Handle("GET", "/interview/result/:interviewID", clientside(interview.GetResult))

	router.Handle("POST", "/interview/list", serverside(interview.GetUniversal))

	//router.Handle("POST", "/interview/interview", middleware.CORS(interview.Create))

	return router
}
