package main

import (
	"poll-api/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	pollApi := api.NewPollApi()

	r.GET("/polls", pollApi.ListAllPolls)
	r.GET("/polls/:id", pollApi.ListPoll)
	r.POST("/polls", pollApi.AddPoll)

	serverPath := "0.0.0.0:2080"
	r.Run(serverPath)
}
