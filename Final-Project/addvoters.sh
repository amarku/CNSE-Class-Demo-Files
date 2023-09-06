#!/bin/bash
curl -d '{ "voterID": 1, "firstName": "Andrew", "lastName": "Marku" }' -H "Content-Type: application/json" -X POST http://localhost:3080/votes/1/voter
curl -d '{ "voterID": 2, "firstName": "Jose", "lastName": "Guzman Rodriguez" }' -H "Content-Type: application/json" -X POST http://localhost:3080/votes/2/voter
curl -d '{ "voterID": 3, "firstName": "Steven", "lastName": "Marku" }' -H "Content-Type: application/json" -X POST http://localhost:3080/votes/3/voter
