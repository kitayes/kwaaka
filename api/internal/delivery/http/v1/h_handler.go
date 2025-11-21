package v1

import (
	"context"
	"kwaaka/api/internal/application"
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
	Port         string        `env:"PORT"`
	ReadTimeOut  time.Duration `env:"READ_TIMEOUT"`
	WriteTimeOut time.Duration `env:"WRITE_TIMEOUT"`
}
type Handler struct {
	services   *application.Service
	router     *gin.Engine
	httpServer *http.Server
	cfg        *Config
	logger     Logger
}

func NewHandler(services *application.Service, cfg *Config, logger Logger) *Handler {
	return &Handler{
		services: services,
		cfg:      cfg,
		logger:   logger,
	}
}

func (h *Handler) Run(_ context.Context) {

	h.httpServer = &http.Server{
		Addr:         ":" + h.cfg.Port,
		Handler:      h.router,
		ReadTimeout:  h.cfg.ReadTimeOut,
		WriteTimeout: h.cfg.WriteTimeOut,
	}
	go func() {
		if err := h.httpServer.ListenAndServe(); err != nil {
			log.Println("listen: %s\n", err.Error())
			return
		}
	}()
}

func (h *Handler) Stop() {
	err := h.httpServer.Shutdown(context.Background())
	if err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *Handler) Init() error {
	router := gin.New()

	router.POST("/parse", h.CreateParseTask)
	router.GET("/parse/:task_id", h.GetParseTask)

	router.GET("/menu/:menu_id", h.GetMenu)

	router.PATCH("/products/:product_id/status", h.UpdateProductStatus)

	h.router = router
	return nil
}
