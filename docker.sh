#!/bin/bash

function error {
	echo "error: $1"
	exit 1
}

if [ -z "$1" ]; then
	error "no Docker image specified, try ./docker.sh danbondd/vanity"
fi

DOCKER_IMAGE=$1

echo -e "\033[1mChecking dependencies...\033[0m"
command -v docker >/dev/null 2>&1 || fail "Docker not installed"
command -v go >/dev/null 2>&1 || fail "Go not installed"
echo "Success"

echo -e "\033[1mBuilding binary...\033[0m"
GOOS=linux go build -o vanity . || error "building binary failed"
echo "Success"

echo -e "\033[1mBuilding Docker image...\033[0m"
docker build -t $DOCKER_IMAGE . || error "building Docker image failed"

echo -e "\n\033[1;32mSuccessfully built docker image $DOCKER_IMAGE\033[0m"

