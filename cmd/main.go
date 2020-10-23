package main

import (
	"database/sql"
	httputils "github.com/Solar-2020/GoUtils/http"
	"github.com/Solar-2020/GoUtils/http/errorWorker"
	"github.com/Solar-2020/Interview-Backend/cmd/handlers"
	interviewHandler "github.com/Solar-2020/Interview-Backend/cmd/handlers/interview"
	"github.com/Solar-2020/Interview-Backend/internal/services/interview"
	"github.com/Solar-2020/Interview-Backend/internal/storages/answerStorage"
	"github.com/Solar-2020/Interview-Backend/internal/storages/interviewStorage"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})

	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	interviewDB, err := sql.Open("postgres", cfg.InterviewDataBaseConnectionString)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	interviewDB.SetMaxIdleConns(5)
	interviewDB.SetMaxOpenConns(10)

	//userDB, err := sql.Open("postgres", cfg.UserDataBaseConnectionString)
	//if err != nil {
	//	log.Fatal().Msg(err.Error())
	//	return
	//}

	//userDB.SetMaxIdleConns(5)
	//userDB.SetMaxOpenConns(10)

	errorWorker := errorWorker.NewErrorWorker()

	interviewStorage := interviewStorage.NewStorage(interviewDB)
	answerStorage := answerStorage.NewStorage(interviewDB)
	interviewService := interview.NewService(interviewStorage, answerStorage)
	interviewTransport := interview.NewTransport()

	interviewHandler := interviewHandler.NewHandler(interviewService, interviewTransport, errorWorker)

	middlewares := httputils.NewMiddleware()

	server := fasthttp.Server{
		Handler: handlers.NewFastHttpRouter(interviewHandler, middlewares).Handler,
	}

	go func() {
		log.Info().Str("msg", "start server").Str("port", cfg.Port).Send()
		if err := server.ListenAndServe(":" + cfg.Port); err != nil {
			log.Error().Str("msg", "server run failure").Err(err).Send()
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	defer func(sig os.Signal) {

		log.Info().Str("msg", "received signal, exiting").Str("signal", sig.String()).Send()

		if err := server.Shutdown(); err != nil {
			log.Error().Str("msg", "server shutdown failure").Err(err).Send()
		}

		//dbConnection.Shutdown()
		log.Info().Str("msg", "goodbye").Send()
	}(<-c)
}
