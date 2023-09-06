package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"votes-api/schema"
)

type VoteApi struct {
	voteList    *schema.VoteList
	pollAPIURL  string
	voterAPIURL string
}

func NewVotesApi() *VoteApi {
	return &VoteApi{
		schema.NewVoteList(),
		os.Getenv("POLL_API_URL"),
		os.Getenv("VOTER_API_URL"),
	}
}

func (p *VoteApi) AddVote(c *gin.Context) {
	var newVote schema.Vote

	if err := c.ShouldBindJSON(&newVote); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := p.voteList.AddVote(newVote); err != nil {
		log.Println("Error adding vote: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, newVote)
}

func (p *VoteApi) GetVote(c *gin.Context, voteID uint) *schema.Vote {
	vote, err := p.voteList.GetVote(voteID)
	if err != nil {
		log.Println("error getting vote: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return nil
	}
	return vote
}

func (p *VoteApi) ListAllVotes(c *gin.Context) {
	voteList, err := p.voteList.GetAllVotes()
	if err != nil {
		log.Println("Error getting all vote: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if voteList == nil {
		voteList = make([]schema.Vote, 0)
	}
	c.JSON(http.StatusOK, voteList)
}

func (p *VoteApi) ListVote(c *gin.Context) {
	id := getIDFromContext(c, "id")
	vote := p.GetVote(c, id)

	c.JSON(http.StatusOK, vote)
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
