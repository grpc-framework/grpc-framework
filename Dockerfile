FROM golang
LABEL org.opencontainers.image.authors="lpabon@purestorage.com"

ENV GOPATH=/go
RUN apt update

##
## grpc-framework additions
##
# Install tools
COPY ./gfw/tools/grpcfw* /usr/local/bin/
# Add protofiles
RUN mkdir -p /go/src/github.com/grpc-framework/grpc-framework/v2
COPY . /go/src/github.com/grpc-framework/grpc-framework/v2

##
## Install software specific to this arch
##
RUN bash /go/src/github.com/grpc-framework/grpc-framework/v2/hack/docker-build.sh

##
## Set working directory
##
RUN mkdir -p /go/src/code
WORKDIR /go/src/code
