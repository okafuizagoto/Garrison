#!/bin/bash

source ./scripts/functions.sh

check_docker_active

printf "${GREEN}Starting Gold Gym services...${NC}\n\n"

# Copy app.env if missing
if [ ! -f "app.env" ]; then
    echo "[setup] Copying .env.example -> app.env"
    cp .env.example app.env
fi

docker compose up -d "$@"

printf "\n${GREEN}Services started. Use ./logs.sh <service> to tail logs.${NC}\n"
printf "  Available: users, members, db, redis, rabbitmq, reverse-proxy, mailcatcher\n\n"
