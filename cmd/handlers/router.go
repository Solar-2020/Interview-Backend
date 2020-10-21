package handlers

import (
	"fmt"
	interviewHandler "github.com/Solar-2020/Interview-Backend/cmd/handlers/interview"
	"github.com/Solar-2020/Interview-Backend/internal/errorWorker"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"runtime/debug"
)

func NewFastHttpRouter(interview interviewHandler.Handler, middleware Middleware) *fasthttprouter.Router {
	router := fasthttprouter.New()

	//router.Handle("GET", "/health", check)

	router.PanicHandler = panicHandler

	router.Handle("POST", "/interview/interview", middleware.CORS(interview.Create))
	//router.Handle("GET", "/interview/interview/interview", middleware.CORS(interview.GetList))

	router.Handle("POST", "/interview/interview/result/:interviewID", middleware.CORS(interview.SetAnswer))
	router.Handle("GET", "/interview/interview/result/:interviewID", middleware.CORS(interview.GetResult))

	//router.Handle("POST", "/interview/interview", middleware.CORS(interview.Create))

	return router
}

func panicHandler(ctx *fasthttp.RequestCtx, err interface{}) {
	fmt.Printf("Request falied with panic: %s, error: %v\nTrace:\n", string(ctx.Request.RequestURI()), err)
	fmt.Println(string(debug.Stack()))
	errorWorker.NewErrorWorker().ServeFatalError(ctx)
}
