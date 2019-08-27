package main

import (
	"github.com/bannovd/evaluator/repository"
	"os"

	"github.com/go-kit/kit/log"

	"github.com/bannovd/evaluator/application"
	"github.com/bannovd/evaluator/models"
	"github.com/bannovd/evaluator/service"
)

var (
	appConfig models.Config
)

func init() {
	models.LoadConfig(&appConfig)
}

func main() {
	logger := log.With(
		log.NewJSONLogger(os.Stderr),
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)

	rep := repository.NewRepository(appConfig.ServerOpt.CacheCleanupInterval)
	svc := service.NewService(rep)

	app := application.NewApplication(svc, appConfig.HashSum, appConfig.ServerOpt, logger)
	app.Start()
}
