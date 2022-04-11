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

SIDECAR_PROJECT_NAME 	= mosn.io/mosn
PROJECT_NAME			= mosn.io/extensions/go-plugin

# only for local test
SIDECAR_GITLAB_PROJECT_NAME = gitlab.alipay-inc.com/ant-mesh/mosn

SIDECAR_BUILD_IMAGE     = mosn-builder
BUILD_IMAGE				= golang:1.14.13
BASE_IMAGE				= zonghaishang/delve:v1.6.1

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

ant:
	make ant.arch

ant.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make ant.arch

ant.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make ant.arch

pkg-codec:
    ifeq ("$(os_arch)","Darwin_arm64")
		DOCKER_BUILD_OPTS="--platform=linux/amd64" make pkg-codec.arch
    else
		make pkg-codec.arch
    endif

pkg-codec.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make pkg-codec.arch

pkg-codec.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make pkg-codec.arch

pkg-filter:
    ifeq ("$(os_arch)","Darwin_arm64")
		DOCKER_BUILD_OPTS="--platform=linux/amd64" make pkg-filter.arch
    else
		make pkg-filter.arch
    endif

pkg-filter.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make pkg-filter.arch

pkg-filter.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make pkg-filter.arch

pkg-trans:
    ifeq ("$(os_arch)","Darwin_arm64")
		DOCKER_BUILD_OPTS="--platform=linux/amd64" make pkg-trans.arch
    else
		make pkg-trans.arch
    endif

pkg-trans.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make pkg-trans.arch

pkg-trans.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make pkg-trans.arch

pkg-ant:
	make pkg-ant.arch

pkg-ant.amd64:
	GOOS=linux GOARCH=amd64 DOCKER_BUILD_OPTS="--platform=linux/amd64" make pkg-ant.arch

pkg-ant.arm64:
	GOOS=linux GOARCH=arm64 DOCKER_BUILD_OPTS="--platform=linux/arm64" make pkg-ant.arch

start:
	bash $(shell go env GOPATH)/src/${PROJECT_NAME}/etc/ant/start.sh

debug:
	DLV_DEBUG=true bash $(shell go env GOPATH)/src/${PROJECT_NAME}/etc/ant/start.sh

stop:
	bash $(shell go env GOPATH)/src/${PROJECT_NAME}/etc/ant/stop.sh

clean:
	@rm -rf build/codecs build/stream_filters build/transcoders build/target

codec.arch:
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
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-v $(shell go env GOPATH)/pkg/mod/:/go/pkg/mod/ \
	-w /go/src/${PROJECT_NAME} -it ${BUILD_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/compile-codec.sh

stream-filter.arch:
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_STREAM_FILTER=${TARGET} \
	-e PLUGIN_STEAM_FILTER_PREFIX=${STEAM_FILTER_PREFIX} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
    -e PLUGIN_ARCH=${PLUGIN_GOARCH} \
    -v $(shell go env GOPATH)/pkg/mod/:/go/pkg/mod/ \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/compile-filter.sh

transcoder.arch:
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_TRANSCODER=${TARGET} \
	-e PLUGIN_TRANSCODER_PREFIX=${TRANSCODER_PREFIX} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
    -e PLUGIN_ARCH=${PLUGIN_GOARCH} \
    -v $(shell go env GOPATH)/pkg/mod/:/go/pkg/mod/ \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/compile-transcoder.sh

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

ant.arch:
	bash $(shell go env GOPATH)/src/${PROJECT_NAME}/etc/ant/stop.sh
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PROJECT_NAME=${PROJECT_NAME} \
	-e SIDECAR_PROJECT_NAME=${SIDECAR_GITLAB_PROJECT_NAME} \
	-v $(shell go env GOPATH)/pkg/mod/:/go/pkg/mod/ \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-v $(shell go env GOPATH)/src/${SIDECAR_GITLAB_PROJECT_NAME}:/go/src/${SIDECAR_GITLAB_PROJECT_NAME} \
	-w /go/src/${SIDECAR_GITLAB_PROJECT_NAME} -it ${BUILD_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/compile.sh

pkg-codec.arch:
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_CODEC_PREFIX=${CODEC_PREFIX} \
	-e PLUGIN_CODEC=${TARGET} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PLUGIN_BUILD_PLATFORM=$(shell uname) \
	-e PLUGIN_BUILD_PLATFORM_ARCH=$(shell arch) \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BASE_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/package-codec.sh

pkg-filter.arch:
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_STREAM_FILTER=${TARGET} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-e PLUGIN_STEAM_FILTER_PREFIX=${STEAM_FILTER_PREFIX} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PLUGIN_BUILD_PLATFORM=$(shell uname) \
	-e PLUGIN_BUILD_PLATFORM_ARCH=$(shell arch) \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BASE_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/package-filter.sh

pkg-trans.arch:
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_TRANSCODER=${TARGET} \
	-e PLUGIN_TRANSCODER_PREFIX=${TRANSCODER_PREFIX} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PLUGIN_BUILD_PLATFORM=$(shell uname) \
	-e PLUGIN_BUILD_PLATFORM_ARCH=$(shell arch) \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BASE_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/package-transcoder.sh

pkg-ant.arch:
	docker run --rm ${DOCKER_BUILD_OPTS} \
	-e PLUGIN_OS=${PLUGIN_GOOS} \
	-e PLUGIN_ARCH=${PLUGIN_GOARCH} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BASE_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/package-ant.sh

.PHONY: codec clean start stop debug
