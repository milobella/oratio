#!/usr/bin/env bash

echo "== Variables initialization :"

GITLAB_HOST="gitlab.milobella.com"
echo "GITLAB_HOST=\"$GITLAB_HOST\""

PROJECT_NAME="milobella"
echo "PROJECT_NAME=\"$PROJECT_NAME\""

BUILD_VERSION=$(cat VERSION.txt)
echo "BUILD_VERSION=\"$BUILD_VERSION\""

MODULE_NAME="oratio"
echo "MODULE_NAME=\"$MODULE_NAME\""

IMAGE="$GITLAB_HOST/$PROJECT_NAME/$MODULE_NAME:$BUILD_VERSION"
echo "IMAGE=\"$IMAGE\""

echo "== Pushing docker image [$IMAGE]"
docker push ${IMAGE} || exit 1
echo "Docker image pushed ! [$IMAGE]"