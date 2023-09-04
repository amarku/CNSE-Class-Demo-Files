package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"poll-api/schema"
	"strconv"
)

type PollApi struct {
	pollList *schema.PollList
}

func NewPollApi() *PollApi {
	return &PollApi{
		schema.NewPollList(),
	}
}

func (p *PollApi) AddPoll(c *gin.Context) {
	var newPoll schema.Poll

	if err := c.ShouldBindJSON(&newPoll); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := p.pollList.AddPoll(newPoll); err != nil {
		log.Println("Error adding poll: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, newPoll)
}

func (p *PollApi) GetPoll(c *gin.Context, pollID uint) *schema.Poll {
	poll, err := p.pollList.GetPoll(pollID)
	if err != nil {
		log.Println("error getting poll: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return nil
	}
	return poll
}

func (p *PollApi) ListAllPolls(c *gin.Context) {
	pollList, err := p.pollList.GetAllPolls()
	if err != nil {
		log.Println("Error getting all poll: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if pollList == nil {
		pollList = make([]schema.Poll, 0)
	}
	c.JSON(http.StatusOK, pollList)
}

func (p *PollApi) ListPoll(c *gin.Context) {
	id := getIDFromContext(c, "id")
	poll := p.GetPoll(c, id)

	c.JSON(http.StatusOK, poll)
}

func getIDFromContext(c *gin.Context, param string) uint {
	idString := c.Param(param)
	id, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		log.Println("Error converting id to integer: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
	return uint(id)
}
