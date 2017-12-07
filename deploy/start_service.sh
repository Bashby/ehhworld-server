#!/bin/sh
set -Eeuxo pipefail

docker network create ehh-world-network

docker-compose build

docker-compose up