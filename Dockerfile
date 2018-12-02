FROM golang:latest
LABEL maintainer="celian.garcia1@gmail.com"

# Some arguments used for labelling
ARG BUILD_DATE
ARG VCS_REF
ARG BUILD_VERSION
ARG SSH_RSA_KEY
ARG GITLAB_HOST
ARG PROJECT_NAME
ARG MODULE_NAME
ARG MODULE_DESCRIPTION

# Labels.
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date=$BUILD_DATE
LABEL org.label-schema.name="$PROJECT_NAME::$MODULE_NAME"
LABEL org.label-schema.description=$MODULE_DESCRIPTION
LABEL org.label-schema.url="https://www.$PROJECT_NAME.com/"
LABEL org.label-schema.vcs-url="https://$GITLAB_HOST/$PROJECT_NAME/$MODULE_NAME"
LABEL org.label-schema.vcs-ref=$VCS_REF
LABEL org.label-schema.version=$BUILD_VERSION
LABEL org.label-schema.docker.cmd="docker run -it $GITLAB_HOST/$PROJECT_NAME/$MODULE_NAME:$BUILD_VERSION"

# Install dependency management tool
RUN go get -u github.com/golang/dep/cmd/dep

# Push the current repository into the srcs and define it as working dir
ENV GOPATH_SOURCES="$GOPATH/src"
ENV APPLICATION_SOURCES="$GOPATH_SOURCES/$GITLAB_HOST/$PROJECT_NAME/$MODULE_NAME"
COPY . $APPLICATION_SOURCES
WORKDIR $APPLICATION_SOURCES

# gitlab.milobella.com security (necessary for dep ensure)
RUN echo "[url \"git@$GITLAB_HOST:\"]\n\tinsteadOf = https://$GITLAB_HOST/" >> /root/.gitconfig
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config
ADD $SSH_RSA_KEY /root/.ssh
RUN chmod 700 /root/.ssh/id_rsa

# Install dependencies using the private RSA key
RUN dep ensure
RUN go build -o /bin/main cmd/$MODULE_NAME/main.go
ENV CONFIGURATION_PATH=/etc/$MODULE_NAME.toml
RUN cp config/$MODULE_NAME.toml ${CONFIGURATION_PATH}

# Clean step (necessary for security considerations)
WORKDIR /
RUN rm -rf /root/.ssh
RUN rm /root/.gitconfig

# Build the main command
CMD /bin/main --configfile ${CONFIGURATION_PATH}
