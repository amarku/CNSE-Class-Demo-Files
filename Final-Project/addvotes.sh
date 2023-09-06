#!/bin/bash
curl -d '{ "voteID": 1, "voterID": 1, "pollID": 1, "voteValue": 1 }' -H "Content-Type: application/json" -X POST http://localhost:3080/votes
curl -d '{ "voteID": 2, "voterID": 2, "pollID": 1, "voteValue": 2 }' -H "Content-Type: application/json" -X POST http://localhost:3080/votes
curl -d '{ "voteID": 3, "voterID": 3, "pollID": 1, "voteValue": 4 }' -H "Content-Type: application/json" -X POST http://localhost:3080/votes
