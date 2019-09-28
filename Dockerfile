FROM golang:1.13.1
LABEL maintainer="celian.garcia1@gmail.com"

# Some arguments used for labelling
ARG BUILD_DATE
ARG VCS_REF
ARG BUILD_VERSION
ARG GITLAB_TOKEN
ARG PROJECT_NAME
ARG MODULE_NAME
ARG MODULE_DESCRIPTION
ARG DOCKER_IMAGE

# Labels.
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date=$BUILD_DATE
LABEL org.label-schema.name="$PROJECT_NAME::$MODULE_NAME"
LABEL org.label-schema.description=$MODULE_DESCRIPTION
LABEL org.label-schema.url="https://www.$PROJECT_NAME.com/"
LABEL org.label-schema.vcs-url="https://milobella.com/gitlab/$PROJECT_NAME/$MODULE_NAME"
LABEL org.label-schema.vcs-ref=$VCS_REF
LABEL org.label-schema.version=$BUILD_VERSION
LABEL org.label-schema.docker.cmd="docker run -it $DOCKER_IMAGE:$BUILD_VERSION"

# Push the current repository into the srcs and define it as working dir
ENV GOPATH_SOURCES="$GOPATH/src"
ENV GOPRIVATE="milobella.com"
ENV APPLICATION_SOURCES="$GOPATH_SOURCES/milobella.com/gitlab/$PROJECT_NAME/$MODULE_NAME"
COPY . $APPLICATION_SOURCES
WORKDIR $APPLICATION_SOURCES

# milobella.com security (necessary for go mod dependencies)
RUN git config --global url."https://oauth2:${GITLAB_TOKEN}@milobella.com/gitlab".insteadOf "https://milobella.com/gitlab"

# Build the ability
RUN go build -o /bin/main cmd/ability/main.go

# Remove milobella token
RUN git config --global --remove-section url."https://oauth2:${GITLAB_TOKEN}@milobella.com/gitlab"

# Build the main command
CMD /bin/main
