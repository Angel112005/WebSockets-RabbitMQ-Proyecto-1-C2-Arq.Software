package main

import (
	"my-notification-api/config"
	"my-notification-api/infrastructure/api"
	"my-notification-api/infrastructure/persistence"
	"my-notification-api/application"
	"my-notification-api/server"
)

func main() {
	config.ConnectDatabase()

	repo := persistence.NewNotificationMySQL(config.DB)
	useCase := application.NewCreateNotification(repo)
	controller := api.NewNotificationController(useCase)

	r := server.SetupRouter(controller)
	r.Run(":8081")
}
