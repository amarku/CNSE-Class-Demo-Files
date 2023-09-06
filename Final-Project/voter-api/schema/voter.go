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

type Voter struct {
	VoterID   uint   `json:"voterID"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "voter:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

type VoterList struct {
	cache
}

// NewVoterList constructor for VoterList struct
func NewVoterList() *VoterList {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	voterList := NewWithCacheInstance(redisUrl)

	return voterList
}

func NewWithCacheInstance(location string) *VoterList {
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

	return &VoterList{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}
}

func (v *Voter) ToJson() string {
	b, _ := json.Marshal(v)
	return string(b)
}

func (vl *VoterList) AddVoter(newVoter Voter) error {
	redisKey := redisKeyFromId(newVoter.VoterID)
	var existingVoter Voter
	if err := vl.getVoterFromRedis(redisKey, &existingVoter); err == nil {
		return errors.New("voter with that ID already exists")
	}
	if _, err := vl.jsonHelper.JSONSet(redisKey, ".", newVoter); err != nil {
		return err
	}
	return nil
}

func (vl *VoterList) GetVoter(id uint) (*Voter, error) {
	var voter Voter
	pattern := redisKeyFromId(id)
	err := vl.getVoterFromRedis(pattern, &voter)
	if err != nil {
		return &Voter{}, err
	}
	return &voter, nil
}

func (vl *VoterList) GetAllVoters() ([]Voter, error) {
	var voterList []Voter
	var voter Voter

	pattern := RedisKeyPrefix + "*"
	ks, err := vl.cacheClient.Keys(vl.context, pattern).Result()
	if err != nil {
		return nil, err
	}
	for _, key := range ks {
		log.Println(key)
		err = vl.getVoterFromRedis(key, &voter)
		if err != nil {
			return nil, err
		}
		voterList = append(voterList, voter)
	}

	sort.Slice(voterList, func(i, j int) bool {
		return voterList[i].VoterID < voterList[j].VoterID
	})

	return voterList, nil
}

func redisKeyFromId(id uint) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

func (vl *VoterList) getVoterFromRedis(key string, voter *Voter) error {
	voterObject, err := vl.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(voterObject.([]byte), voter)
	if err != nil {
		return err
	}
	return nil
}
