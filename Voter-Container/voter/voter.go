package voter

import (
	"encoding/json"
	"time"
)

type voterPoll struct {
	PollID   uint
	VoteDate time.Time
}

type Voter struct {
	VoterID     uint        `json:"voterID"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	VoteHistory []voterPoll `json:"voteHistory"`
}
type VoterList struct {
	Voters map[uint]Voter //A map of VoterIDs as keys and Voter structs as values
}

// constructor for VoterList struct
func NewVoterList() VoterList {
	return VoterList{
		make(map[uint]Voter),
	}
}

func NewVoter(id uint, fn, ln string) *Voter {
	return &Voter{
		FirstName:   fn,
		LastName:    ln,
		VoteHistory: []voterPoll{},
	}
}

func NewSampleVoter() *Voter {
	return &Voter{
		VoterID:   1,
		FirstName: "John",
		LastName:  "Doe",
		VoteHistory: []voterPoll{
			{PollID: 1, VoteDate: time.Now()},
		},
	}
}

func (v *Voter) AddPoll(pollID uint) {
	v.VoteHistory = append(v.VoteHistory, voterPoll{PollID: pollID, VoteDate: time.Now()})
}

func (v *Voter) AddPollWithTimeDetails(pollID uint, timeOfPoll time.Time) {
	v.VoteHistory = append(v.VoteHistory, voterPoll{PollID: pollID, VoteDate: timeOfPoll})
}

func (v *Voter) ToJson() string {
	b, _ := json.Marshal(v)
	return string(b)
}

func (v *Voter) GetVoteHistory() []time.Time {
	var pollDates []time.Time
	for _, history := range v.VoteHistory {
		pollDates = append(pollDates, history.VoteDate)
	}

	return pollDates
}

func (v *Voter) GetPollById(pollId uint) time.Time {
	return v.VoteHistory[pollId].VoteDate
}
