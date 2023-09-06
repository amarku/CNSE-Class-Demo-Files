package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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

func (p *VoteApi) ListPoll(c *gin.Context) {
	voteID := getIDFromContext(c, "id")
	poll := p.GetPoll(c, voteID)

	c.JSON(http.StatusOK, poll)
}

func (p *VoteApi) ListVoter(c *gin.Context) {
	voteID := getIDFromContext(c, "id")
	voter := p.GetVoter(c, voteID)

	c.JSON(http.StatusOK, voter)
}

func (p *VoteApi) GetPoll(c *gin.Context, voteID uint) *schema.Poll {
	var poll schema.Poll
	vote := p.GetVote(c, voteID)
	urlRequest := fmt.Sprintf("%s/polls/%d", p.pollAPIURL, vote.PollID)
	resp, err := http.Get(urlRequest)
	if err != nil {
		log.Println("Error getting poll from vote: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return nil
	}

	err = json.Unmarshal(body, &poll)
	if err != nil {
		log.Println("Error decoding response body to json: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return nil
	}

	return &poll
}

func (p *VoteApi) GetVoter(c *gin.Context, voteID uint) *schema.Voter {
	var voter schema.Voter
	vote := p.GetVote(c, voteID)
	urlRequest := fmt.Sprintf("%s/voters/%d", p.voterAPIURL, vote.PollID)
	resp, err := http.Get(urlRequest)
	if err != nil {
		log.Println("Error getting voter from vote: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return nil
	}

	err = json.Unmarshal(body, &voter)
	if err != nil {
		log.Println("Error decoding response body to json: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return nil
	}

	return &voter
}

func (p *VoteApi) ListVote(c *gin.Context) {
	id := getIDFromContext(c, "id")
	vote := p.GetVote(c, id)
	vote.Links = make(map[string]string)

	vote.Links["poll"] = fmt.Sprintf("votes/%d/poll", vote.VoteID)
	vote.Links["voter"] = fmt.Sprintf("votes/%d/voter", vote.VoteID)

	c.JSON(http.StatusOK, vote)
}

func (p *VoteApi) AddVoter(c *gin.Context) {
	var newVoter schema.Voter
	var returnedVoter schema.Voter
	id := getIDFromContext(c, "id")
	vote := p.GetVote(c, id)

	if err := c.ShouldBindJSON(&newVoter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if vote.VoterID != newVoter.VoterID {
		log.Println("Error adding voter: voter IDs do not match")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	urlRequest := fmt.Sprintf("%s/voters", p.voterAPIURL)
	data, errMarshal := json.Marshal(&newVoter)
	if errMarshal != nil {
		log.Println("Error encoding voter to json: ", errMarshal)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	resp, errPost := http.Post(urlRequest, "application/json", bytes.NewReader(data))
	defer resp.Body.Close()
	if errPost != nil {
		log.Println("Error sending new voter http request: ", errPost)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		log.Println("Error reading body of http response: ", errBody)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	errUnmarshal := json.Unmarshal(body, &returnedVoter)
	if errUnmarshal != nil {
		log.Println("Error decoding body of http response: ", errUnmarshal)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, &returnedVoter)
}

func (p *VoteApi) AddPoll(c *gin.Context) {
	var newPoll schema.Poll
	var returnedPoll schema.Poll
	id := getIDFromContext(c, "id")
	vote := p.GetVote(c, id)

	if err := c.ShouldBindJSON(&newPoll); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if vote.PollID != newPoll.PollID {
		log.Println("Error adding poll: poll IDs do not match")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	urlRequest := fmt.Sprintf("%s/polls", p.pollAPIURL)
	data, errMarshal := json.Marshal(&newPoll)
	if errMarshal != nil {
		log.Println("Error encoding poll to json: ", errMarshal)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	resp, errPost := http.Post(urlRequest, "application/json", bytes.NewReader(data))
	defer resp.Body.Close()
	if errPost != nil {
		log.Println("Error sending new poll http request: ", errPost)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		log.Println("Error reading body of http response: ", errBody)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	errUnmarshal := json.Unmarshal(body, &returnedPoll)
	if errUnmarshal != nil {
		log.Println("Error decoding body of http response: ", errUnmarshal)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, &returnedPoll)
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
