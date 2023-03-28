SHELL			= /bin/bash
TARGET			= ${plugin}
OUTPUT			= ${TARGET}.so

STREAM_FILTER = ${filter}
TRANSCODER = ${trans}

CODEC_PREFIX	= codec
CODEC_OUTPUT	= ${CODEC_PREFIX}-${OUTPUT}
CODEC_ZIP_OUTPUT = ${TARGET}.zip
TRANSCODER_ZIP_OUTPUT = ${TARGET}.zip

PLUGIN_GOOS = ${GOOS}
PLUGIN_GOARCH = ${GOARCH}

TRANSCODER_PREFIX = transcoder
TRANSCODER_OUTPUT = ${TRANSCODER_PREFIX}-${OUTPUT}

STEAM_FILTER_PREFIX = filter

TRACE_PREFIX = trace

GIT_VERSION 			= $(shell git log -1 --pretty=format:%H | head -c 8)
PROJECT_NAME			= $(shell head -1 go.mod | cut  -d' ' -f2)
FULL_PROJECT_NAME		= $(shell pwd)

# only for local test
SIDECAR_GITLAB_PROJECT_NAME = gitlab.alipay-inc.com/ant-mesh/mosn

SIDECAR_BUILD_IMAGE     = mosn-builder
BUILD_IMAGE				= golang:1.20.1
BASE_IMAGE				= zonghaishang/delve:v1.20.1

os_arch:=$(shell uname)_$(shell arch)

codec:
	make codec.arch

codec.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make codec.arch

codec.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make codec.arch

filter: # alias for stream-filter (more short)
	make stream-filter.arch

filter.amd64:
	make stream-filter.amd64

filter.arm64:
	make stream-filter.arm64

stream-filter.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make stream-filter.arch

stream-filter.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make stream-filter.arch

trans: # alias for transcoder (more short)
	make transcoder.arch

trans.amd64:
	make transcoder.amd64

trans.arm64:
	make transcoder.arm64

transcoder.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make transcoder.arch

transcoder.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make transcoder.arch

trace:
	make trace.arch

trace.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make trace.arch

trace.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make trace.arch

ant:
	make ant.arch

ant.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make ant.arch

ant.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make ant.arch

pkg-codec:
ifeq ("$(os_arch)","Darwin_arm64")
	make pkg-codec.amd64
else ifeq ("$(os_arch)","Darwin_i386")
	make pkg-codec.amd64
else
	make pkg-codec.arch
endif

pkg-codec.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make pkg-codec.arch

pkg-codec.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make pkg-codec.arch

pkg-filter:
ifeq ("$(os_arch)","Darwin_arm64")
	make pkg-filter.amd64
else ifeq ("$(os_arch)","Darwin_i386")
	make pkg-filter.amd64
else
	make pkg-filter.arch
endif

pkg-filter.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make pkg-filter.arch

pkg-filter.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make pkg-filter.arch

pkg-trans:
ifeq ("$(os_arch)","Darwin_arm64")
	make pkg-trans.amd64
else ifeq ("$(os_arch)","Darwin_i386")
	make pkg-trans.amd64
else
	make pkg-trans.arch
endif

pkg-trans.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make pkg-trans.arch

pkg-trans.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make pkg-trans.arch

##
pkg-trace:
ifeq ("$(os_arch)","Darwin_arm64")
	make pkg-trace.amd64
else ifeq ("$(os_arch)","Darwin_i386")
	make pkg-trace.amd64
else
	make pkg-trace.arch
endif

pkg-trace.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make pkg-trace.arch

pkg-trace.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make pkg-trace.arch

pkg-ant:
	make pkg-ant.arch

pkg-ant.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make pkg-ant.arch

pkg-ant.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make pkg-ant.arch

start:
    ifeq ("$(os_arch)","Darwin_arm64")
		make start.arm64
    else
		make start.amd64
    endif

start.amd64:
	FULL_PROJECT_NAME=${FULL_PROJECT_NAME} PROJECT_NAME=${PROJECT_NAME} \
	SIDECAR_GITLAB_PROJECT_NAME=${SIDECAR_GITLAB_PROJECT_NAME} bash ${FULL_PROJECT_NAME}/etc/ant/start.sh

start.arm64:
	FULL_PROJECT_NAME=${FULL_PROJECT_NAME} PROJECT_NAME=${PROJECT_NAME} DOCKER_BUILD_OPTS="--platform=linux/arm64" \
    SIDECAR_GITLAB_PROJECT_NAME=${SIDECAR_GITLAB_PROJECT_NAME} bash ${FULL_PROJECT_NAME}/etc/ant/start.sh

debug:
	DLV_DEBUG=true make start

stop:
	bash ${FULL_PROJECT_NAME}/etc/ant/stop.sh

clean:
	@rm -rf build/codecs build/stream_filters build/transcoders build/target

codec.arch:
	mkdir -p /tmp/go-build-cache
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_TARGET=${TARGET} \
	-e PLUGIN_CODEC_OUTPUT=${CODEC_OUTPUT} \
	-e PLUGIN_CODEC_PREFIX=${CODEC_PREFIX} \
    -e PLUGIN_STEAM_FILTER_PREFIX=${STEAM_FILTER_PREFIX} \
    -e PLUGIN_TRANSCODER_PREFIX=${TRANSCODER_PREFIX} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-e PLUGIN_STREAM_FILTER=${STREAM_FILTER} \
	-e PLUGIN_TRANSCODER=${TRANSCODER} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PLUGIN_GIT_VERSION=${GIT_VERSION} \
	-v ${FULL_PROJECT_NAME}:/go/src/${PROJECT_NAME} \
	-v $(shell go env GOPATH)/pkg/mod/:/go/pkg/mod/ \
	-v /tmp/go-build-cache:/root/.cache/go-build \
	-w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/compile-codec.sh

stream-filter.arch:
	mkdir -p /tmp/go-build-cache
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_STREAM_FILTER=${TARGET} \
	-e PLUGIN_STEAM_FILTER_PREFIX=${STEAM_FILTER_PREFIX} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PLUGIN_GIT_VERSION=${GIT_VERSION} \
	-v $(shell go env GOPATH)/pkg/mod/:/go/pkg/mod/ \
	-v /tmp/go-build-cache:/root/.cache/go-build \
	-v ${FULL_PROJECT_NAME}:/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/compile-filter.sh

transcoder.arch:
	mkdir -p /tmp/go-build-cache
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_TRANSCODER=${TARGET} \
	-e PLUGIN_TRANSCODER_PREFIX=${TRANSCODER_PREFIX} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PLUGIN_GIT_VERSION=${GIT_VERSION} \
	 -v $(shell go env GOPATH)/pkg/mod/:/go/pkg/mod/ \
	-v ${FULL_PROJECT_NAME}:/go/src/${PROJECT_NAME} \
	-v /tmp/go-build-cache:/root/.cache/go-build \
	-w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/compile-transcoder.sh

trace.arch:
	mkdir -p /tmp/go-build-cache
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_TRACE=${TARGET} \
	-e PLUGIN_TRACE_PREFIX=${TRACE_PREFIX} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PLUGIN_GIT_VERSION=${GIT_VERSION} \
	 -v $(shell go env GOPATH)/pkg/mod/:/go/pkg/mod/ \
	-v ${FULL_PROJECT_NAME}:/go/src/${PROJECT_NAME} \
	-v /tmp/go-build-cache:/root/.cache/go-build \
	-w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/compile-trace.sh

compile-codec:
	@rm -rf build/codecs/${PLUGIN_TARGET}
	GO111MODULE=on ${BUILD_OPTS} go build \
		--buildmode=plugin -mod=mod \
		-gcflags "all=-N -l" \
		-o ${PLUGIN_CODEC_OUTPUT} ${PLUGIN_PROJECT_NAME}/plugins/codecs/${PLUGIN_TARGET}/main
	mkdir -p build/codecs/${PLUGIN_TARGET}
	mv ${PLUGIN_CODEC_OUTPUT} build/codecs/${PLUGIN_TARGET}
	@cd build/codecs/${PLUGIN_TARGET} && $(shell which md5sum) -b ${PLUGIN_CODEC_OUTPUT} | cut -d' ' -f1  > ${PLUGIN_CODEC_PREFIX}-${PLUGIN_TARGET}.md5
	cp configs/codecs/${PLUGIN_TARGET}/*.json build/codecs/${PLUGIN_TARGET}
	@echo "compile codec ${PLUGIN_TARGET} complete"

compile-stream-filter:
	@rm -rf build/stream_filters/${PLUGIN_TARGET}
	GO111MODULE=on ${BUILD_OPTS} go build \
		--buildmode=plugin -mod=mod \
		-gcflags "all=-N -l" \
		-o ${PLUGIN_STEAM_FILTER_OUTPUT} ${PLUGIN_PROJECT_NAME}/plugins/stream_filters/${PLUGIN_TARGET}/main
	mkdir -p build/stream_filters/${PLUGIN_TARGET}
	mv ${PLUGIN_STEAM_FILTER_OUTPUT} build/stream_filters/${PLUGIN_TARGET}
	@cd build/stream_filters/${PLUGIN_TARGET} && $(shell which md5sum) -b ${PLUGIN_STEAM_FILTER_OUTPUT} | cut -d' ' -f1  > ${PLUGIN_STEAM_FILTER_PREFIX}-${PLUGIN_TARGET}.md5
	cp configs/stream_filters/${PLUGIN_TARGET}/*.json build/stream_filters/${PLUGIN_TARGET}
	@echo "compile filter ${PLUGIN_TARGET} complete"

compile-transcoder:
	@rm -rf build/transcoders/${PLUGIN_TARGET}
	GO111MODULE=on ${BUILD_OPTS} go build \
		--buildmode=plugin -mod=mod \
		-gcflags "all=-N -l" \
		-o ${PLUGIN_TRANSCODER_OUTPUT} ${PLUGIN_PROJECT_NAME}/plugins/transcoders/${PLUGIN_TARGET}/main
	mkdir -p build/transcoders/${PLUGIN_TARGET}
	mv ${PLUGIN_TRANSCODER_OUTPUT} build/transcoders/${PLUGIN_TARGET}
	@cd build/transcoders/${PLUGIN_TARGET} && $(shell which md5sum) -b ${PLUGIN_TRANSCODER_OUTPUT} | cut -d' ' -f1  > ${PLUGIN_TRANSCODER_PREFIX}-${PLUGIN_TARGET}.md5
	cp configs/transcoders/${PLUGIN_TARGET}/*.json build/transcoders/${PLUGIN_TARGET}
	@echo "compile transcoder ${PLUGIN_TARGET} complete"

compile-trace:
	@rm -rf build/traces/${PLUGIN_TARGET}
	@rm -rf build/codecs/bundle/support/traces/${PLUGIN_TARGET}
	GO111MODULE=on ${BUILD_OPTS} go build \
		--buildmode=plugin -mod=mod \
		-gcflags "all=-N -l" \
		-o ${PLUGIN_TRACE_OUTPUT} ${PLUGIN_PROJECT_NAME}/plugins/traces/${PLUGIN_TARGET}/main
	@mkdir -p build/traces/${PLUGIN_TARGET}
	mkdir -p build/codecs/bundle/support/traces/${PLUGIN_TARGET}
	@mv ${PLUGIN_TRACE_OUTPUT} build/traces/${PLUGIN_TARGET}
	@cd build/traces/${PLUGIN_TARGET} && $(shell which md5sum) -b ${PLUGIN_TRACE_OUTPUT} | cut -d' ' -f1  > ${PLUGIN_TRACE_PREFIX}-${PLUGIN_TARGET}.md5
	@cp configs/traces/${PLUGIN_TARGET}/*.json build/traces/${PLUGIN_TARGET}
	cp -r build/traces/${PLUGIN_TARGET} build/codecs/bundle/support/traces/
	@echo "compile trace ${PLUGIN_TARGET} complete"

ant.arch:
	bash ${FULL_PROJECT_NAME}/etc/ant/stop.sh
	@test ! -f ~/.netrc  && touch ~/.netrc || true
	docker run -it --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PROJECT_NAME=${PROJECT_NAME} \
	-e GIT_TERMINAL_PROMPT=1 \
	-e SIDECAR_PROJECT_NAME=${SIDECAR_GITLAB_PROJECT_NAME} \
	-v $(shell go env GOPATH)/pkg/mod/:/go/pkg/mod/ \
	-v ${FULL_PROJECT_NAME}:/go/src/${PROJECT_NAME} \
	-v $(shell go env GOCACHE):/root/.cache/go-build \
	-v ~/.netrc:/root/.netrc \
	-v $(shell go env GOPATH)/src/${SIDECAR_GITLAB_PROJECT_NAME}:/go/src/${SIDECAR_GITLAB_PROJECT_NAME} \
	-w /go/src/${SIDECAR_GITLAB_PROJECT_NAME} ${BUILD_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/compile.sh

pkg-codec.arch:
	mkdir -p /tmp/go-build-cache
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_CODEC_PREFIX=${CODEC_PREFIX} \
	-e PLUGIN_CODEC=${TARGET} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PLUGIN_BUILD_PLATFORM=$(shell uname) \
	-e PLUGIN_BUILD_PLATFORM_ARCH=$(shell arch) \
	-e PLUGIN_GIT_VERSION=${GIT_VERSION} \
	-v ${FULL_PROJECT_NAME}:/go/src/${PROJECT_NAME} \
	-v /tmp/go-build-cache:/root/.cache/go-build \
	-v $(shell go env GOPATH)/pkg/mod/:/go/pkg/mod/ \
	-w /go/src/${PROJECT_NAME} ${BASE_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/package-codec.sh
	@echo "package codec ${TARGET} complete"

pkg-filter.arch:
	mkdir -p /tmp/go-build-cache
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_STREAM_FILTER=${TARGET} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-e PLUGIN_STEAM_FILTER_PREFIX=${STEAM_FILTER_PREFIX} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PLUGIN_BUILD_PLATFORM=$(shell uname) \
	-e PLUGIN_BUILD_PLATFORM_ARCH=$(shell arch) \
	-e PLUGIN_GIT_VERSION=${GIT_VERSION} \
	-v ${FULL_PROJECT_NAME}:/go/src/${PROJECT_NAME} \
	-v /tmp/go-build-cache:/root/.cache/go-build \
	-v $(shell go env GOPATH)/pkg/mod/:/go/pkg/mod/ \
	-w /go/src/${PROJECT_NAME} ${BASE_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/package-filter.sh
	@echo "package filter ${TARGET} complete"

pkg-trans.arch:
	mkdir -p /tmp/go-build-cache
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_TRANSCODER=${TARGET} \
	-e PLUGIN_TRANSCODER_PREFIX=${TRANSCODER_PREFIX} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PLUGIN_BUILD_PLATFORM=$(shell uname) \
	-e PLUGIN_BUILD_PLATFORM_ARCH=$(shell arch) \
	-e PLUGIN_GIT_VERSION=${GIT_VERSION} \
	-v ${FULL_PROJECT_NAME}:/go/src/${PROJECT_NAME} \
	-v /tmp/go-build-cache:/root/.cache/go-build \
	-v $(shell go env GOPATH)/pkg/mod/:/go/pkg/mod/ \
	-w /go/src/${PROJECT_NAME} ${BASE_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/package-transcoder.sh
	@echo "package transcoder ${TARGET} complete"

pkg-trace.arch:
	mkdir -p /tmp/go-build-cache
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_TRACE=${TARGET} \
	-e PLUGIN_TRACE_PREFIX=${TRACE_PREFIX} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PLUGIN_BUILD_PLATFORM=$(shell uname) \
	-e PLUGIN_BUILD_PLATFORM_ARCH=$(shell arch) \
	-e PLUGIN_GIT_VERSION=${GIT_VERSION} \
	-v ${FULL_PROJECT_NAME}:/go/src/${PROJECT_NAME} \
	-v /tmp/go-build-cache:/root/.cache/go-build \
	-v $(shell go env GOPATH)/pkg/mod/:/go/pkg/mod/ \
	-w /go/src/${PROJECT_NAME} ${BASE_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/package-trace.sh
	@echo "package trace ${TARGET} complete"

pkg-ant.arch:
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-v ${FULL_PROJECT_NAME}:/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BASE_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/package-ant.sh

.PHONY: codec clean start stop debug
