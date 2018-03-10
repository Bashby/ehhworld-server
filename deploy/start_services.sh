#!/bin/bash
set -Eeuo pipefail

# Colors
GREEN='\033[0;32m'
RESET='\033[0m' # No Color

# Create our network, if not already exists
printf "${GREEN}Checking for docker network 'ehh-world-dev-network' ...${RESET}\n"
if docker network inspect ehh-world-dev-network > /dev/null; then
    printf "${GREEN}Network already exists, joining ...${RESET}\n"
else
    printf "${GREEN}Network doesn't exist, creating ...${RESET}\n"
    docker network create ehh-world-dev-network
fi

# Clean-up existing containers
printf "${GREEN}Cleaning up currently running services ...${RESET}\n"
docker-compose down --remove-orphans

# Build our services
printf "${GREEN}Building services ...${RESET}\n"
docker-compose build

# Start our services
printf "${GREEN}Starting services ...${RESET}\n"
docker-compose up

# Clean-up
printf "${GREEN}Finished executing ...${RESET}\n"
