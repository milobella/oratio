#!/usr/bin/env bash

echo "== Variables initialization :"

# User can give a custom id_rsa, by default, we take the private/id_rsa path.
LOCAL_PRIVATE_RSA_KEY=private/id_rsa
PARAM_PRIVATE_RSA_KEY=${1:-${LOCAL_PRIVATE_RSA_KEY}}
# If user gave custom path, we copy it to the local path
if [[ ${PARAM_PRIVATE_RSA_KEY} != ${LOCAL_PRIVATE_RSA_KEY} ]]; then
    cp ${PARAM_PRIVATE_RSA_KEY} ${LOCAL_PRIVATE_RSA_KEY}
    echo "Copied your ssh rsa key locally into $LOCAL_PRIVATE_RSA_KEY ."
fi
echo "LOCAL_PRIVATE_RSA_KEY=\"$LOCAL_PRIVATE_RSA_KEY\""

BUILD_VERSION=$(cat VERSION.txt)
echo "BUILD_VERSION=\"$BUILD_VERSION\""

BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
echo "BUILD_DATE=\"$BUILD_DATE\""

GITLAB_HOST="gitlab.milobella.com"
echo "GITLAB_HOST=\"$GITLAB_HOST\""

PROJECT_NAME="milobella"
echo "PROJECT_NAME=\"$PROJECT_NAME\""

MODULE_NAME="oratio"
echo "MODULE_NAME=\"$MODULE_NAME\""

MODULE_DESCRIPTION="Http server which is the main and only entry point of Milobella."
echo "MODULE_DESCRIPTION=\"$MODULE_DESCRIPTION\""

IMAGE="$GITLAB_HOST/$PROJECT_NAME/$MODULE_NAME:$BUILD_VERSION"
echo "IMAGE=\"$IMAGE\""


echo "== Building docker image [$IMAGE]"
docker build -t ${IMAGE} \
    --add-host ${GITLAB_HOST}:192.168.1.15 \
    --build-arg BUILD_VERSION=${BUILD_VERSION} \
    --build-arg BUILD_DATE=${BUILD_DATE} \
    --build-arg GITLAB_HOST=${GITLAB_HOST} \
    --build-arg PROJECT_NAME=${PROJECT_NAME} \
    --build-arg MODULE_NAME=${MODULE_NAME} \
    --build-arg MODULE_DESCRIPTION="${MODULE_DESCRIPTION}" \
    --build-arg VCS_REF=master \
    --build-arg SSH_RSA_KEY=${LOCAL_PRIVATE_RSA_KEY} \
    . || exit 1
echo "Docker image built ! [$IMAGE]"