package main

import (
	"poll-api/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	voterApi := api.NewPollApi()

	r.GET("/voters", voterApi.ListAllPolls)
	r.GET("/voters/:id", voterApi.ListPoll)
	r.POST("/voters", voterApi.AddPoll)
	r.GET("/voters/health", voterApi.GetHealth)

	serverPath := "0.0.0.0:1080"
	r.Run(serverPath)
}
