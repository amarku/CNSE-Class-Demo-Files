#!/bin/bash
cd poll-api
./build.sh
cd ../voter-api
./build.sh
cd ../votes-api
./build.sh
cd ..
