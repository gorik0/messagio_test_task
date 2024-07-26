package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"log"
	"messagio/internal/models"
	"messagio/internal/web"
	"net/http"
)

type Handler struct {
	Mux  *gin.Engine
	Srvr Servicer
}

type Servicer interface {
	CloseDB() error
	HandleMessage(msg models.JsonMsg) error
	GiveMeStats() (int, int, error)
	StartConsumeMessages()
}

func NewHandler(mux *gin.Engine, srvr Servicer) *Handler {
	return &Handler{
		Mux:  mux,
		Srvr: srvr,
	}
}
func (h *Handler) SetupRoutes() {

	h.Mux.GET("/", h.Index)
	h.Mux.POST("/message", h.MessageRecieve)
	h.Mux.GET("/stat", h.GetStat)
	h.Mux.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (h *Handler) Index(ctx *gin.Context) {
	web.Render(ctx.Writer, "test.html.gohtml")
}

func (h *Handler) MessageRecieve(ctx *gin.Context) {
	var mes models.JsonMsg
	err := ctx.BindJSON(&mes)
	if err != nil {
		log.Printf("Error parsing json: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error to marshal request body"})
		return
	}

	err = h.Srvr.HandleMessage(mes)
	if err != nil {

		log.Printf("Error handling message: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error handle message"})
		return

	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": "received"})
}

func (h *Handler) GetStat(ctx *gin.Context) {
	total, processed, err := h.Srvr.GiveMeStats()
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error getting stats"})
		return
	}
	msg := models.GetStat{
		Total:     total,
		Processed: processed,
	}
	ctx.JSON(http.StatusOK, msg)
}
