#!/bin/bash
set -Eeuo pipefail

# Colors
GREEN='\033[0;32m'
REDBG_WHITEFG='\033[41m'
RESET='\033[0m' # No Color
PRIVATE_KEY=${REPO_PRIVATE_KEY:-}

# Check that we have a private key for pulling our repo
[ -z "${PRIVATE_KEY}" ] && { printf "${REDBG_WHITEFG}SSH Private Key Missing!\nYou MUST set REPO_PRIVATE_KEY to a working SSL private key to pull from the EhhWorld Server bitbucket git repository!${RESET}"; exit 1; }

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
docker-compose down

# Build our services
printf "${GREEN}Building services ...${RESET}\n"
docker-compose build

# Wait for dependencies
printf "${GREEN}Waiting on dependencies to start ...${RESET}\n"
docker-compose run --rm start_dependencies

# Start our services
printf "${GREEN}Starting services ...${RESET}\n"
docker-compose up ehhio_server

# Clean-up
printf "${GREEN}Finished executing ...${RESET}\n"
