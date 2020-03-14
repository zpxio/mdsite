.DEFAULT_GOAL := build

# Aliases
GOCMD=go
GO_BUILD=$(GOCMD) build
GO_CLEAN=$(GOCMD) clean
GO_TEST=$(GOCMD) test
GO_GET=$(GOCMD) get

BUILD_DIR := out
BUILD_EXE := mdsite
BUILD_TARGET := ${BUILD_DIR}/${BUILD_EXE}

BASEDIR := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))

# Optional stuff for demos/samples
SITE_PORT := 9999
SITE_BASE := ${BASEDIR}/sample
SITE_PATH := ${SITE_BASE}/site
SITE_CONF := ${SITE_BASE}/config

build:
	@-mkdir -p ${BUILD_DIR}
	@-echo "BUILD: ${BUILD_TARGET}"
	$(GO_BUILD) -o $(BUILD_TARGET) -v cmd/mdsite/main.go

run: build
	./${BUILD_TARGET}

sample: build
	@echo "Base Path: ${BASEDIR}"
	./${BUILD_TARGET} --port ${SITE_PORT} --site ${SITE_PATH} --config ${SITE_CONF}

test:
	$(GO_TEST) -v ./...

clean:
	$(GO_CLEAN)
	rm -f $(BUILD_TARGET)
