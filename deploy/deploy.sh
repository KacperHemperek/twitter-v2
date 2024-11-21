#! /usr/bin/bash

docker compose -p twitter stop

docker pull kacperhemperek/tw-api:latest
docker pull kacperhemperek/tw-web:latest

docker compose -p twitter up -d

echo "Application Started"
