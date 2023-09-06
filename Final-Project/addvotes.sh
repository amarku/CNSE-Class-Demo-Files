#!/bin/bash
curl -d '{ "voteID": 1, "voterID": 1, "pollID": 1, "voteValue": 1 }' -H "Content-Type: application/json" -X POST http://localhost:3080/votes
