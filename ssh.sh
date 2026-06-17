#!/bin/bash

source ./scripts/functions.sh

check_docker_active

if [ -z $1 ]; then
    printf "${RED}Please provide the service name (e.g. users, members, db)${NC}\n\n"
    exit 1
fi

printf "${GREEN}Connecting shell to goldgym-$1 ...${NC}\n\n"

check_container_running "goldgym-$1"

if [ $? == 1 ]; then
    docker exec -it "goldgym-$1" sh
    exit
fi

printf "${RED}[Error] Container 'goldgym-$1' is not running.${NC}\n\n"
