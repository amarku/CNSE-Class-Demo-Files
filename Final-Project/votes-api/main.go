package main

import (
	"votes-api/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	votesApi := api.NewVotesApi()

	r.GET("/votes", votesApi.ListAllVotes)
	r.GET("/votes/:id", votesApi.ListVote)
	r.POST("/votes", votesApi.AddVote)
	r.GET("/votes/:id/poll", votesApi.ListPoll)
	r.GET("/votes/:id/voter", votesApi.ListVoter)

	serverPath := "0.0.0.0:3080"
	r.Run(serverPath)
}
