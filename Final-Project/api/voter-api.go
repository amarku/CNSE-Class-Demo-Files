package api

import (
	"Votes-HATEOS/voter"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type VoterApi struct {
	voterList *voter.VoterList
}

func NewVoterApi() *VoterApi {
	return &VoterApi{
		voter.NewVoterList(),
	}
}

func (v *VoterApi) AddVoter(c *gin.Context) {
	//v.voterList.Voters[voterID] = *voter.NewVoter(voterID, firstName, lastName)
	var newVoter voter.Voter

	if err := c.ShouldBindJSON(&newVoter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := v.voterList.AddVoter(newVoter.VoterID, newVoter); err != nil {
		log.Println("Error adding voter: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
	}

	c.JSON(http.StatusOK, newVoter)
}

func (v *VoterApi) AddPoll(c *gin.Context, voterID, pollID uint) {
	voter, err := v.voterList.GetVoter(voterID)
	if err != nil {
		log.Println("error getting voter ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	voter.AddPoll(pollID)
}

func (v *VoterApi) GetVoter(c *gin.Context, voterID uint) *voter.Voter {
	voter, err := v.voterList.GetVoter(voterID)
	if err != nil {
		log.Println("error getting voter: " + err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	return voter
}

func (v *VoterApi) GetVoterJson(c *gin.Context, voterID uint) string {
	voter, err := v.voterList.GetVoter(voterID)
	if err != nil {
		log.Println("error getting voter: " + err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	return voter.ToJson()
}

func (v *VoterApi) GetVoterList() *voter.VoterList {
	return v.voterList
}

func (v *VoterApi) GetVoterListJson() string {
	b, _ := json.Marshal(v.voterList)
	return string(b)
}

func (v *VoterApi) ListAllVoters(c *gin.Context) {
	c.JSON(http.StatusOK, v.voterList)
}

func (v *VoterApi) ListVoter(c *gin.Context) {
	id := getIDFromContext(c, "id")
	voter := v.GetVoter(c, id)

	c.JSON(http.StatusOK, voter)
}

func (v *VoterApi) ListPollHistory(c *gin.Context) {
	id := getIDFromContext(c, "id")
	voter := v.GetVoter(c, id)
	voteHistory := voter.GetVoteHistory()

	c.JSON(http.StatusOK, voteHistory)
}

func (v *VoterApi) ListSinglePollData(c *gin.Context) {
	id := getIDFromContext(c, "id")
	pollId := getIDFromContext(c, "pollid")

	voter := v.GetVoter(c, id)
	pollData := voter.GetPollById(pollId)

	c.JSON(http.StatusOK, pollData)
}

func (v *VoterApi) AddPollData(c *gin.Context) {
	id := getIDFromContext(c, "id")
	pollId := getIDFromContext(c, "pollid")

	voter := v.GetVoter(c, id)
	voter.AddPoll(pollId)
}

func (v *VoterApi) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, "API is running")
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
