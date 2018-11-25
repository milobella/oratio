FROM golang:latest
LABEL maintainer="celian.garcia1@gmail.com"

# Some arguments used for labelling
ARG BUILD_DATE
ARG VCS_REF
ARG BUILD_VERSION

# Labels.
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date=$BUILD_DATE
LABEL org.label-schema.name="milobella/oratio"
LABEL org.label-schema.description="Http server which is the main and only entry point of Milobella."
LABEL org.label-schema.url="https://www.milobella.com/"
LABEL org.label-schema.vcs-url="https://gitlab.milobella.com/milobella/oratio"
LABEL org.label-schema.vcs-ref=$VCS_REF
LABEL org.label-schema.version=$BUILD_VERSION
LABEL org.label-schema.docker.cmd="docker run -it -v ~/.ssh/id_rsa:/root/.ssh/id_rsa milobella/oratio:$BUILD_VERSION"

# Install dependency management tool
RUN go get -u github.com/golang/dep/cmd/dep

# Push the current repository into the srcs and define it as working dir
ENV APPLICATION_SOURCES="$GOPATH/src/gitlab.milobella.com/milobella/oratio"
ADD . $APPLICATION_SOURCES
WORKDIR $APPLICATION_SOURCES

# gitlab.milobella.com security (necessary for dep ensure)
RUN echo "[url \"git@gitlab.milobella.com:\"]\n\tinsteadOf = https://gitlab.milobella.com/" >> /root/.gitconfig
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config

###  Build the main command ############################################################################################
# We put the dep ensure here because it depends on two things which should be provided at the run (more clean)
# - The rsa key because we use ssh instead of https (dep ensure do some git clone):
#           -v ~/.ssh/id_rsa:/root/.ssh/id_rsa
# - The host alias for gitlab.milobella.com (only if developer machine is in the same network than gitlab)
#           --add-host gitlab.milobella.com:<gitlab-internal-host>
CMD dep ensure \
    && go build -o /bin/main cmd/oratio/main.go \
    && /bin/main
