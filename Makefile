# Paths to packages
GO=$(shell which go)
ARCH=$(shell which arch)
UNAME=$(shell which uname)

# Modules
SERVER_MODULE = "github.com/djthorpe/go-rotel"

# Paths to locations, etc
BUILD_DIR := "build"
BUILD_ARCH := $(shell ${ARCH}  | tr A-Z a-z)
BUILD_PLATFORM := $(shell ${UNAME}  | tr A-Z a-z)
BUILD_VERSION := $(shell git describe --tags  | sed 's/^v//')
BUILD_TAG := rotel-${BUILD_PLATFORM}-${BUILD_ARCH}:${BUILD_VERSION}
CMD_DIR := $(wildcard cmd/*)

# Add linker flags
BUILD_LD_FLAGS += -X $(SERVER_MODULE)/pkg/version.GitSource=${SERVER_MODULE}
BUILD_LD_FLAGS += -X $(SERVER_MODULE)/pkg/version.GitTag=${BUILD_VERSION}
BUILD_LD_FLAGS += -X $(SERVER_MODULE)/pkg/version.GitBranch=$(shell git name-rev HEAD --name-only --always)
BUILD_LD_FLAGS += -X $(SERVER_MODULE)/pkg/version.GitHash=$(shell git rev-parse HEAD)
BUILD_LD_FLAGS += -X $(SERVER_MODULE)/pkg/version.GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BUILD_FLAGS = -ldflags "-s -w $(BUILD_LD_FLAGS)" 

all: clean cmd

cmd: $(CMD_DIR)

$(CMD_DIR): dependencies mkdir
	@echo Build cmd $(notdir $@)
	@${GO} build ${BUILD_FLAGS} -o ${BUILD_DIR}/$(notdir $@) ./$@

test: dependencies
	@echo Running tests
	@${GO} test ./pkg/...

docker: cmd
	@echo Building docker image: ${BUILD_TAG}
	docker build \
		--tag ${BUILD_TAG} \
		--build-arg VERSION=${BUILD_VERSION} \
		--build-arg ARCH=${BUILD_ARCH} \
		--build-arg PLATFORM=${BUILD_PLATFORM} \
		-f etc/docker/Dockerfile .
FORCE:

dependencies:
	@test -f "${GO}" && test -x "${GO}"  || (echo "Missing go binary" && exit 1)
	@test -f "${ARCH}" && test -x "${ARCH}"  || (echo "Missing arch binary" && exit 1)
	@test -f "${UNAME}" && test -x "${UNAME}"  || (echo "Missing uname binary" && exit 1)

mkdir:
	@echo Mkdir ${BUILD_DIR}
	@install -d ${BUILD_DIR}

clean:
	@echo Clean
	@rm -fr $(BUILD_DIR)
	@${GO} mod tidy
	@${GO} clean

