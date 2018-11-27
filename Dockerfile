FROM golang:latest
LABEL maintainer="celian.garcia1@gmail.com"

# Some arguments used for labelling
ARG BUILD_DATE
ARG VCS_REF
ARG BUILD_VERSION
ARG SSH_RSA_KEY

# Define some environment variables used in docker build
ENV GITLAB_URL="gitlab.milobella.com"
ENV PROJECT_NAME="milobella/oratio"

# Labels.
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date=$BUILD_DATE
LABEL org.label-schema.name=$PROJECT_NAME
LABEL org.label-schema.description="Http server which is the main and only entry point of Milobella."
LABEL org.label-schema.url="https://www.milobella.com/"
LABEL org.label-schema.vcs-url="https://$GITLAB_URL/$PROJECT_NAME"
LABEL org.label-schema.vcs-ref=$VCS_REF
LABEL org.label-schema.version=$BUILD_VERSION
LABEL org.label-schema.docker.cmd="docker run -it $GITLAB_URL/$PROJECT_NAME:$BUILD_VERSION"

# Install dependency management tool
RUN go get -u github.com/golang/dep/cmd/dep

# Push the current repository into the srcs and define it as working dir
ENV GOPATH_SOURCES="$GOPATH/src"
ENV APPLICATION_SOURCES="$GOPATH_SOURCES/$GITLAB_URL/$PROJECT_NAME"
ADD . $APPLICATION_SOURCES
WORKDIR $APPLICATION_SOURCES

# gitlab.milobella.com security (necessary for dep ensure)
RUN echo "[url \"git@$GITLAB_URL:\"]\n\tinsteadOf = https://$GITLAB_URL/" >> /root/.gitconfig
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config
ADD $SSH_RSA_KEY /root/.ssh
RUN chmod 700 /root/.ssh/id_rsa

# Install dependencies using the private RSA key
RUN dep ensure
RUN go build -o /bin/main cmd/oratio/main.go

# Clean step (necessary for security and size considerations)
RUN rm -r /root/.ssh && rm /root/.gitconfig && rm -r $GOPATH_SOURCES 

#  Build the main command 
CMD /bin/main
