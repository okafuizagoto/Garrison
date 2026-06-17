#!/bin/bash

source ./scripts/functions.sh

check_docker_active

if [ -z $1 ]; then
    printf "${RED}Please provide the service name (e.g. users, members, db, redis)${NC}\n\n"
    exit 1
fi

printf "${GREEN}Connecting to container goldgym-$1 ...${NC}\n\n"

check_container_running "goldgym-$1"

if [ $? == 1 ]; then
    docker logs --tail 100 -f "goldgym-$1"
    exit
fi

printf "${RED}[Error] Container 'goldgym-$1' does not exist or is not running.${NC}\n\n"
