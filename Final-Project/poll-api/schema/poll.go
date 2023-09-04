package schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
	"log"
	"os"
	"sort"
)

type pollOption struct {
	PollOptionID    uint   `json:"pollOptionID"`
	PollOptionValue string `json:"pollOptionValue"`
}

type Poll struct {
	PollID       uint         `json:"pollID"`
	PollTitle    string       `json:"pollTitle"`
	PollQuestion string       `json:"pollQuestion"`
	PollOptions  []pollOption `json:"pollOptions"`
}

const (
	RedisNilError        = "redis: nil"
	RedisKeyPrefix       = "poll:"
	RedisDefaultLocation = "0.0.0.0:6379"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

type PollList struct {
	cache
}

func NewPollList() *PollList {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	pollList := NewWithCacheInstance(redisUrl)

	return pollList
}

func NewWithCacheInstance(location string) *PollList {
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

	return &PollList{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}
}

func NewPoll(id uint, title, question string) *Poll {
	return &Poll{
		PollID:       id,
		PollTitle:    title,
		PollQuestion: question,
		PollOptions:  []pollOption{},
	}
}

func NewSamplePoll() *Poll {
	return &Poll{
		PollID:       1,
		PollTitle:    "Favorite Pet",
		PollQuestion: "What type of pet do you like best?",
		PollOptions: []pollOption{
			{PollOptionID: 1, PollOptionValue: "Dog"},
			{PollOptionID: 2, PollOptionValue: "Cat"},
			{PollOptionID: 3, PollOptionValue: "Fish"},
			{PollOptionID: 4, PollOptionValue: "Bird"},
			{PollOptionID: 5, PollOptionValue: "NONE"},
		},
	}
}

func (p *Poll) ToJson() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func (pl *PollList) AddPoll(newPoll Poll) error {
	redisKey := redisKeyFromId(newPoll.PollID)
	var existingPoll Poll
	if err := pl.getPollFromRedis(redisKey, &existingPoll); err == nil {
		return errors.New("poll with that ID already exists")
	}
	if _, err := pl.jsonHelper.JSONSet(redisKey, ".", newPoll); err != nil {
		return err
	}
	return nil
}

func (pl *PollList) GetPoll(id uint) (*Poll, error) {
	var poll Poll
	pattern := redisKeyFromId(id)
	err := pl.getPollFromRedis(pattern, &poll)
	if err != nil {
		return &Poll{}, err
	}
	return &poll, nil
}

func (pl *PollList) GetAllPolls() ([]Poll, error) {
	var pollList []Poll
	var poll Poll

	pattern := RedisKeyPrefix + "*"
	ks, err := pl.cacheClient.Keys(pl.context, pattern).Result()
	if err != nil {
		return nil, err
	}
	for _, key := range ks {
		log.Println(key)
		err = pl.getPollFromRedis(key, &poll)
		if err != nil {
			return nil, err
		}
		pollList = append(pollList, poll)
	}

	sort.Slice(pollList, func(i, j int) bool {
		return pollList[i].PollID < pollList[j].PollID
	})

	return pollList, nil
}

func redisKeyFromId(id uint) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

func (pl *PollList) getPollFromRedis(key string, Poll *Poll) error {
	pollObject, err := pl.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(pollObject.([]byte), Poll)
	if err != nil {
		return err
	}
	return nil
}
