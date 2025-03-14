TAG = dev
GRPC_FRAMEWORK_TAG=latest
GRPC_FRAMEWORK_CONTAINER=quay.io/grpc-framework/grpc-framework:$(GRPC_FRAMEWORK_TAG)
DOCKERCMD=docker run \
		--privileged --rm \
		-v $(shell pwd):/go/src/code \
		-e "LINT_OUTPUT=$(LINT_OUTPUT)" \
		-e "GOPATH=/go" \
		-e "DOCKER_PROTO=yes" \
		-e "PROTO_USER=$(shell id -u)" \
		-e "PROTO_GROUP=$(shell id -g)" \
		-e "PATH=/bin:/usr/bin:/usr/local/bin:/go/bin:/usr/local/go/bin" \
		$(GRPC_FRAMEWORK_CONTAINER)

all: apis server client

./bin:
	mkdir ./bin

.PHONY: server
server: ./bin
	go build -o ./bin/server cmd/server.go

.PHONY: client
client: ./bin
	go build -o ./bin/client cmd/client/client.go

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: verify
verify:
	$(TOP)/hack/client-server-test.sh

.PHONY: apis
ifndef DOCKER_PROTO
apis:
	$(DOCKERCMD) make apis
else
apis:
	make -C apis
endif
