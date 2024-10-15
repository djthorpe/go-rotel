# Paths to packages
GO=$(shell which go)
DOCKER=$(shell which docker)
GIT=$(shell which git)

# Set OS and Architecture
ARCH ?= $(shell arch | tr A-Z a-z | sed 's/x86_64/amd64/' | sed 's/i386/amd64/' | sed 's/armv7l/arm/' | sed 's/aarch64/arm64/')
OS ?= $(shell uname | tr A-Z a-z)
VERSION ?= $(shell git describe --tags --always | sed 's/^v//')
DOCKER_REGISTRY ?= ghcr.io/mutablelogic

# Set docker tag, etc
BUILD_TAG := ${DOCKER_REGISTRY}/go-rotel-${OS}-${ARCH}:${VERSION}
ROOT_PATH := $(CURDIR)
BUILD_DIR := build
CMD_DIR := $(wildcard cmd/*)

# Build flags
BUILD_MODULE := $(shell cat go.mod | head -1 | cut -d ' ' -f 2)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitSource=${BUILD_MODULE}
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitTag=$(shell git describe --tags --always)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitBranch=$(shell git name-rev HEAD --name-only --always)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GitHash=$(shell git rev-parse HEAD)
BUILD_LD_FLAGS += -X $(BUILD_MODULE)/pkg/version.GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BUILD_FLAGS = -ldflags "-s -w $(BUILD_LD_FLAGS)" 

all: clean cmd

cmd: mkdir go-dep go-tidy $(CMD_DIR)

$(CMD_DIR): mkdir
	@echo Build cmd $(notdir $@)
	@${GO} build ${BUILD_FLAGS} -o ${BUILD_DIR}/$(notdir $@) ./$@

test: dependencies
	@echo Running tests
	@${GO} test ./pkg/...

docker: docker-dep
	@echo build docker image: ${BUILD_TAG} for ${OS}/${ARCH}
	@${DOCKER} build \
		--tag ${BUILD_TAG} \
		--build-arg ARCH=${ARCH} \
		--build-arg OS=${OS} \
		--build-arg SOURCE=${BUILD_MODULE} \
		--build-arg VERSION=${VERSION} \
		-f etc/docker/Dockerfile .

# Push docker container
docker-push: docker 
	@echo push docker image: ${BUILD_TAG}
	@${DOCKER} push ${BUILD_TAG}

# Check for docker
docker-dep:
	@test -f "${DOCKER}" && test -x "${DOCKER}"  || (echo "Missing docker binary" && exit 1)

# Check for git
git-dep:
	@test -f "${GIT}" && test -x "${GIT}"  || (echo "Missing git binary" && exit 1)

# Check for go
go-dep:
	@test -f "${GO}" && test -x "${GO}"  || (echo "Missing go binary" && exit 1)

# go mod tidy
go-tidy: go-dep
	@echo Tidy
	@go mod tidy

mkdir:
	@echo Mkdir ${BUILD_DIR}
	@install -d ${BUILD_DIR}

clean: go-tidy
	@echo "Cleaning"
	@rm -fr $(BUILD_DIR)

