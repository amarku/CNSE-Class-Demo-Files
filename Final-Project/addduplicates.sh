#!/bin/bash
curl -d '{ "voteID": 2, "voterID": 5, "pollID": 1, "voteValue": 4}' -H "Content-Type: application/json" -X POST http://localhost:3080/votes
curl -d '{ "pollID": 1, "pollTitle": "Favorite Pet", "pollQuestion": "What type of pet do you like best?", "pollOptions": [
{"pollOptionID": 1, "pollOptionValue": "Dog"},
{"pollOptionID": 2, "pollOptionValue": "Cat"},
{"pollOptionID": 3, "pollOptionValue": "Fish"},
{"pollOptionID": 4, "pollOptionValue": "Bird"},
{"pollOptionID": 5, "pollOptionValue": "NONE"}
]}' -H "Content-Type: application/json" -X POST http://localhost:3080/votes/1/poll
curl -d '{ "voterID": 1, "firstName": "John", "lastName": "Doe" }' -H "Content-Type: application/json" -X POST http://localhost:3080/votes/1/voter