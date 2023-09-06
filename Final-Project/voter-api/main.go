package main

import (
	"voter-api/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	voterApi := api.NewVoterApi()

	r.GET("/voters", voterApi.ListAllVoters)
	r.GET("/voters/:id", voterApi.ListVoter)
	r.POST("/voters", voterApi.AddVoter)
	r.GET("/voters/health", voterApi.GetHealth)

	serverPath := "0.0.0.0:1080"
	r.Run(serverPath)
}
