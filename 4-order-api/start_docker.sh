#!/bin/bash

sudo systemctl start docker
docker stop firebird

cd ~/NesK/purple_go/6_app/

docker-compose up -d

docker ps --format "{{.Names}}"
