#! /usr/bin/bash

cd twitter-v2
docker compose -p twitter stop

docker pull kacperhemperek/tw-api:latest
docker pull kacperhemperek/tw-api:latest

docker compose up -d

echo "Application Started"
