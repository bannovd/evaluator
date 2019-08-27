package application

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"

	"github.com/bannovd/evaluator/middleware"
	"github.com/bannovd/evaluator/models"
	"github.com/bannovd/evaluator/service"
)

// Application struct
type Application struct {
	r       *mux.Router
	opt     models.ServerOpt
	hashSum []byte
	svc     *service.Service
	logger  log.Logger
}

// NewApplication return new application
func NewApplication(svc *service.Service, hashSum []byte, opt models.ServerOpt, logger log.Logger) *Application {
	return &Application{
		r:       mux.NewRouter(),
		opt:     opt,
		hashSum: hashSum,
		svc:     svc,
		logger:  logger,
	}
}

func (app *Application) initRoutes() {
	// Basic routes
	app.r.HandleFunc("/health", app.HealthHandler).Methods(http.MethodGet)
	app.r.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	// Catcher routes
	app.r.Path("/api/v1/catch").HandlerFunc(app.CatchHandler) //.Methods(http.MethodPost)
}

// Start application
func (app *Application) Start() {
	app.initRoutes()
	pr := middleware.NewPrometheus("evaluator")
	app.r.Use(pr.PrometheusMiddleware)
	listenErr := make(chan error, 1)
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  time.Duration(app.opt.ReadTimeout),
		IdleTimeout:  time.Duration(app.opt.IdleTimeout),
		WriteTimeout: time.Duration(app.opt.WriteTimeout),
		Handler:      app.r,
	}
	go func() {
		listenErr <- server.ListenAndServe()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-listenErr:
		if err != nil {
			_ = app.logger.Log("err", err)
			os.Exit(1)
		}
	case <-osSignals:
		server.SetKeepAlivesEnabled(false)
		app.svc.Close()
		timeout := time.Second * 5
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		if err := server.Shutdown(ctx); err != nil {
			_ = app.logger.Log("err", err)
			os.Exit(1)
		}
		cancel()
	}
}
