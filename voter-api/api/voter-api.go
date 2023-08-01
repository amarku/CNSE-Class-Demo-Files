package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"voter-api/voter"

	"github.com/gin-gonic/gin"
)

type VoterApi struct {
	voterList voter.VoterList
}

func NewVoterApi() *VoterApi {
	return &VoterApi{
		voter.NewVoterList(),
	}
}

func (v *VoterApi) AddVoter(voterID uint, firstName, lastName string) {
	v.voterList.Voters[voterID] = *voter.NewVoter(voterID, firstName, lastName)
}

func (v *VoterApi) AddPoll(voterID, pollID uint) {
	voter := v.voterList.Voters[voterID]
	voter.AddPoll(pollID)
	v.voterList.Voters[voterID] = voter
}

func (v *VoterApi) GetVoter(voterID uint) voter.Voter {
	voter := v.voterList.Voters[voterID]
	return voter
}

func (v *VoterApi) GetVoterJson(voterID uint) string {
	voter := v.voterList.Voters[voterID]
	return voter.ToJson()
}

func (v *VoterApi) GetVoterList() voter.VoterList {
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
	voter := v.GetVoter(id)

	c.JSON(http.StatusOK, voter)
}

func (v *VoterApi) ListPollHistory(c *gin.Context) {
	id := getIDFromContext(c, "id")
	voter := v.GetVoter(id)
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
