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