package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"

	"voter-api/schema"
)

type VoterApi struct {
	voterList *schema.VoterList
}

func NewVoterApi() *VoterApi {
	return &VoterApi{
		schema.NewVoterList(),
	}
}

func (v *VoterApi) AddVoter(c *gin.Context) {
	var newVoter schema.Voter

	if err := c.ShouldBindJSON(&newVoter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := v.voterList.AddVoter(newVoter); err != nil {
		log.Println("Error adding voter: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, newVoter)
}

func (v *VoterApi) GetVoter(c *gin.Context, voterID uint) *schema.Voter {
	voter, err := v.voterList.GetVoter(voterID)
	if err != nil {
		log.Println("error getting voter-api: " + err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return nil
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

func (v *VoterApi) ListAllVoters(c *gin.Context) {
	voterList, err := v.voterList.GetAllVoters()
	if err != nil {
		log.Println("Error getting all voter: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if voterList == nil {
		voterList = make([]schema.Voter, 0)
	}
	c.JSON(http.StatusOK, voterList)
}

func (v *VoterApi) ListVoter(c *gin.Context) {
	id := getIDFromContext(c, "id")
	voter := v.GetVoter(c, id)

	c.JSON(http.StatusOK, voter)
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
