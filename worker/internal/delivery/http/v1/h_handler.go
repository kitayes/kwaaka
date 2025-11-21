package v1

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type Config struct {
	Port         string        `env:"WORKER_HTTP_PORT" envDefault:"8081"`
	ReadTimeOut  time.Duration `env:"WORKER_HTTP_READ_TIMEOUT" envDefault:"5s"`
	WriteTimeOut time.Duration `env:"WORKER_HTTP_WRITE_TIMEOUT" envDefault:"5s"`
}

type Handler struct {
	router     *gin.Engine
	httpServer *http.Server
	cfg        *Config
	logger     Logger
}

func NewHandler(cfg *Config, logger Logger) *Handler {
	return &Handler{
		cfg:    cfg,
		logger: logger,
	}
}

func (h *Handler) Init() error {
	router := gin.New()

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"role":   "worker",
		})
	})

	h.router = router
	return nil
}

func (h *Handler) Run(_ context.Context) {
	h.httpServer = &http.Server{
		Addr:           ":" + h.cfg.Port,
		Handler:        h.router,
		ReadTimeout:    h.cfg.ReadTimeOut,
		WriteTimeout:   h.cfg.WriteTimeOut,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		if err := h.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("worker http listen error:", err.Error())
			return
		}
	}()
}

func (h *Handler) Stop() {
	if h.httpServer == nil {
		return
	}
	if err := h.httpServer.Shutdown(context.Background()); err != nil {
		h.logger.Error("worker http shutdown err: %v", err)
	}
}
