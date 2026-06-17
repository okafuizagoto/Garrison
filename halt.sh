#!/bin/bash

source ./scripts/functions.sh

check_docker_active

printf "${BROWN}Stopping Gold Gym services...${NC}\n\n"
docker compose stop "$@"
printf "${GREEN}Done.${NC}\n\n"
