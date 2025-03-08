package server

import (
	"my-notification-api/infrastructure/api"

	"github.com/gin-gonic/gin"
)

func SetupRouter(notificationController *api.NotificationController) *gin.Engine {
	r := gin.Default()
	r.POST("/notificaciones", notificationController.Create)
	return r
}
