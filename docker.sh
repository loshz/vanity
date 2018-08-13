#!/usr/bin/env bash

function error {
	echo "error: $1"
	exit 1
}

if [ -z "$1" ]; then
	error "no docker image specified, try ./docker.sh danbondd/vanity"
fi

DOCKER_IMAGE=$1

echo -e "\033[1mChecking dependencies...\033[0m"
command -v docker >/dev/null 2>&1 || fail "docker not installed"
echo "ok"

echo -e "\033[1mBuilding Docker image...\033[0m"
docker build --build-arg DOCKER_IMAGE=${DOCKER_IMAGE} --tag $DOCKER_IMAGE . || error "building Docker image failed"
