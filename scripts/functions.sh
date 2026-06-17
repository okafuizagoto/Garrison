#!/bin/bash
# Shared utility functions for Gold Gym scripts

# Color variables
RED='\033[0;31m'
GREEN='\033[0;32m'
BROWN='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m'

SCRIPTPATH=$(pwd -P)

CONTAINERS='goldgym-db goldgym-redis goldgym-rabbitmq goldgym-mailcatcher goldgym-reverse-proxy goldgym-users goldgym-members'

# Pause and wait for user input
function pause() {
    read -p "$*"
}

function join_by { local IFS="$1"; shift; echo "$*"; }

# Ask for confirmation before destructive actions
function confirmation() {
    printf "${RED}$* ${NC}\n"
    printf "${BROWN}Type the number of your choice:${NC}\n"
    select yn in "Yes, do it" "No, cancel"; do
        case $yn in
            "Yes, do it")
                return 1
                break
                ;;
            "No, cancel")
                printf "${GREEN}Cancelled. Goodbye.${NC}\n"
                exit
                ;;
            *) printf "${PURPLE}Not a valid choice.${NC}\n";;
        esac
    done
}

# Create directory if it does not exist
function create_dir() {
    if [[ ! -e $1 ]]; then
        printf "${GREEN}[create_dir]${NC} $1\n"
        mkdir -p "$1"
    elif [[ ! -d $1 ]]; then
        printf "$1 ${RED}exists but is not a directory${NC}\n" 1>&2
    fi
}

# Returns 0 if directory is empty, 1 if not
function is_empty_dir() {
    if [ "$(ls -A $1)" ]; then
        return 1
    else
        return 0
    fi
}

# Check Docker is running
function check_docker_active() {
    docker_info=$(docker info 2>&1)
    if [[ $docker_info == *"connect to the Docker daemon"* ]]; then
        printf "${RED}Docker is not running. Please start Docker first.${NC}\n"
        exit 1
    fi
    printf "${BROWN}[Check] Docker is active${NC}\n"
}

# Check if a named container is running
function check_container_running() {
    CONTAINER=$1
    RUNNING=$(docker inspect --format="{{ .State.Running }}" $CONTAINER 2>/dev/null)
    if [ "$RUNNING" == "true" ]; then
        return 1
    else
        return 0
    fi
}

printf "${PURPLE}\n"
echo "   ____       _     _      ____"
echo "  / ___| ___ | | __| |    / ___| _   _ _ __ ___"
echo " | |  _ / _ \| |/ _\` |   | |  _ | | | | '_ \` _ \\"
echo " | |_| | (_) | | (_| |   | |_| || |_| | | | | | |"
echo "  \____|\___/|_|\__,_|    \____| \__, |_| |_| |_|"
echo "                                  |___/"
echo ""
echo "   Gold Gym Backend v2 - Development Box"
printf "${NC}\n"
