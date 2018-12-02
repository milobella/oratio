#!/usr/bin/env bash

echo "== Variables initialization :"

BUILD_VERSION=$(cat VERSION.txt)
echo "BUILD_VERSION=\"$BUILD_VERSION\""

GITLAB_HOST="gitlab.milobella.com"
echo "GITLAB_HOST=\"$GITLAB_HOST\""

PROJECT_NAME="milobella"
echo "PROJECT_NAME=\"$PROJECT_NAME\""

MODULE_NAME="oratio"
echo "MODULE_NAME=\"$MODULE_NAME\""

IMAGE="$GITLAB_HOST/$PROJECT_NAME/$MODULE_NAME:$BUILD_VERSION"
echo "IMAGE=\"$IMAGE\""

echo "== Running docker image [$IMAGE]"
docker run -it -e "ORATIO_SERVER_LOG_LEVEL=<root>=DEBUG" ${IMAGE} || exit 1
echo "== Docker image successfully run [$IMAGE]"