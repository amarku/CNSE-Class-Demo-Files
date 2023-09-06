# Final Project

This project demonstrates the use of 3 APIs working together, including hypermedia to help find information.

## Build Docker Containers
Run the script `buildcontainers.sh` to build all 3 Docker containers.

## Run containers

Call `docker compose up`

## Run Test Scripts

`addvotes.sh` adds votes to the votes database in the votes-api.

`addpolls.sh` adds a poll to the polls database through the votes-api, which then calls the polls-api.

`addvoters.sh` adds voters to the votes database through the votes-api, which then calls the voters-api.

`addduplicates.sh` attempts to add items to each API with IDs that are already used.
This demonstrates how the APIs handle this error.

## URLs for Getting Added Data
#### The number 1 is used here for ID Numbers. Change the ID number to see other entries.
all votes from votes api: http://localhost:3080/votes

single vote from votes api: http://localhost:3080/votes/1

all polls from polls api: http://localhost:2080/polls

single poll from polls api: http://localhost:2080/polls/1

all voters from voters api: http://localhost:1080/voters

single voter from voters api: http://localhost:1080/voters/1

### URLs to Get Polls and Voters Data Through Votes API
single poll from votes api: http://localhost:3080/votes/1/poll

single voter from votes api: http://localhost:3080/votes/1/voter
