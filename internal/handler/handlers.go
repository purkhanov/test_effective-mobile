package handler

import (
	"music/internal/controller"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	controller *controller.Controller
}

func NewHandler(controller *controller.Controller) *Handler {
	return &Handler{controller: controller}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://*", "https://*"},
		AllowMethods:  []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:  []string{"Accept", "Authorization", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.GET("", h.controller.GetMusics)
	router.GET(":music_id", h.controller.GetSongLyricsByVerses)
	router.POST("", h.controller.AddMusic)
	router.PATCH(":music_id", h.controller.UpdateMusic)
	router.DELETE(":music_id", h.controller.DeleteMusic)

	return router
}
