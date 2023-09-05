package schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

type Vote struct {
	VoteID    uint
	VoterID   uint
	PollID    uint
	VoteValue uint
	Links     map[string]string
}

func NewVote(pid, vid, vtrid, vval uint) *Vote {
	return &Vote{
		VoteID:    vid,
		VoterID:   vtrid,
		PollID:    pid,
		VoteValue: vval,
	}
}

func NewSampleVote() *Vote {
	return &Vote{
		VoteID:    1,
		PollID:    1,
		VoterID:   1,
		VoteValue: 1,
	}
}

func (p *Vote) ToJson() string {
	b, _ := json.Marshal(p)
	return string(b)
}

const (
	RedisNilError        = "redis: nil"
	RedisKeyPrefix       = "vote:"
	RedisDefaultLocation = "0.0.0.0:6379"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

type VoteList struct {
	cache
}

func NewVoteList() *VoteList {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	voteList := NewWithCacheInstance(redisUrl)

	return voteList
}

func NewWithCacheInstance(location string) *VoteList {
	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	ctx := context.Background()

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error())
		os.Exit(1)
	}

	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	return &VoteList{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}
}

func (vl *VoteList) AddVote(newVote Vote) error {
	redisKey := redisKeyFromId(newVote.VoteID)
	var existingVote Vote
	if err := vl.getVoteFromRedis(redisKey, &existingVote); err == nil {
		return errors.New("vote with that ID already exists")
	}
	if _, err := vl.jsonHelper.JSONSet(redisKey, ".", newVote); err != nil {
		return err
	}
	return nil
}

func (vl *VoteList) GetVote(id uint) (*Vote, error) {
	var vote Vote
	pattern := redisKeyFromId(id)
	err := vl.getVoteFromRedis(pattern, &vote)
	if err != nil {
		return &Vote{}, err
	}
	return &vote, nil
}

func (vl *VoteList) GetAllVotes() ([]Vote, error) {
	var voteList []Vote
	var vote Vote

	pattern := RedisKeyPrefix + "*"
	ks, err := vl.cacheClient.Keys(vl.context, pattern).Result()
	if err != nil {
		return nil, err
	}
	for _, key := range ks {
		log.Println(key)
		err = vl.getVoteFromRedis(key, &vote)
		if err != nil {
			return nil, err
		}
		voteList = append(voteList, vote)
	}

	sort.Slice(voteList, func(i, j int) bool {
		return voteList[i].VoteID < voteList[j].VoteID
	})

	return voteList, nil
}

func redisKeyFromId(id uint) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

func (vl *VoteList) getVoteFromRedis(key string, Vote *Vote) error {
	voteObject, err := vl.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(voteObject.([]byte), Vote)
	if err != nil {
		return err
	}
	return nil
}
