export TOP=$(shell pwd)
export GFW=$(TOP)/gfw
export GRPC_FRAMEWORK_TAG=latest
export GRPC_FRAMEWORK_CONTAINER=quay.io/grpc-framework/grpc-framework:$(GRPC_FRAMEWORK_TAG)
TAG := dev
HAS_ERRCHECK := $(shell command -v errcheck 2> /dev/null)
PKGS := $(shell go list ./... | grep -v vendor | grep -v examples)

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

all: docker-build docker-verify

.PHONY: docker-build
docker-build:
	$(DOCKERCMD) make build

SUBDIRS = apis example
.PHONY: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: build
build: gobuild $(SUBDIRS)

.PHONY: gobuild
gobuild:
	@echo ">>> go build"
	go build $(PKGS)

.PHONY: fmt
fmt:
	@echo ">>> go fmt"
	go fmt $(PKGS) | wc -l | xargs | grep "^0"

.PHONY: vet
vet:
	@echo ">>> go vet"
	@go vet $(PKGS)

.PHONY: errcheck
errcheck:
ifndef HAS_ERRCHECK
	go install github.com/kisielk/errcheck@latest
endif
	@echo ">>> errcheck"
	errcheck $(PKGS)

.PHONY: test
test:
	@echo ">>> go test"
	go test $(PKGS)

.PHONY: verify
verify: vet fmt test
	$(MAKE) -C example verify

.PHONY: docker-verify
docker-verify:
	$(DOCKERCMD) make verify

.PHONY: travis-verify
travis-verify: all pr-verify docker-verify

.PHONY: pr-verify
pr-verify:
	git-validation -run DCO,short-subject
	go mod vendor && git grep -rw GPL vendor | grep LICENSE | egrep -v "yaml.v2" | wc -l | grep "^0"

# Run this after creating and pushing a release tag into the repo
go-mod-publish:
	GOPROXY=proxy.golang.org go list -m github.com/grpc-framework/grpc-framework/v2@$(shell git describe --tags)

proto:
	$(MAKE) -C pkg proto

clean:
	$(MAKE) clean -C test/app

container:
	docker build -t quay.io/grpc-framework/grpc-framework:$(TAG) .

container-buildx-install:
	@echo "Setting up multiarch emulation"
	docker run --privileged --rm tonistiigi/binfmt --install all
	@echo "Setting up multiarch builder"
	docker buildx create --name gfwbuilder --driver docker-container --bootstrap
	docker buildx use gfwbuilder

# Run: make container-buildx-install first to install the emulation
container-release:
	@echo "This will automatically push. Must be logged in to quay.io"
	docker buildx build \
		--push \
		--platform linux/amd64,linux/arm64  \
		--tag quay.io/grpc-framework/grpc-framework:$(TAG) .

container-buildx-uninstall:
	docker buildx stop gfwbuilder
	docker buildx rm gfwbuilder

./venv:
	python3 -m venv venv
	bash -c "source venv/bin/activate && \
		pip3 install --upgrade pip && \
		pip3 install -r requirements.txt"
	@echo "Type: 'source venv/bin/activate' to get access to mkdocs"

doc-env: ./venv

doc-build: doc-env
	bash -c "source venv/bin/activate && \
		cd website && \
		mkdocs build"

doc-serve: doc-env
	bash -c "source venv/bin/activate && \
		cd website && \
		mkdocs serve"

.PHONY: clean proto go-mod-publish travis-verify verify \
	testapp test pr-verify errcheck vet fmt build \
	doc-env doc-build doc-serve container testapp-verify \
	container-buildx-install container-release container-buildx-uninstall

