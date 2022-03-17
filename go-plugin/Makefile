SHELL			= /bin/bash
TARGET			= ${plugin}
OUTPUT			= ${TARGET}.so

STREAM_FILTER = ${filter}
TRANSCODER = ${trans}

CODEC_PREFIX	= codec
CODEC_OUTPUT	= ${CODEC_PREFIX}-${OUTPUT}
CODEC_ZIP_OUTPUT = ${TARGET}.zip
TRANSCODER_ZIP_OUTPUT = ${TARGET}.zip

TRANSCODER_PREFIX = transcoder
TRANSCODER_OUTPUT = ${TRANSCODER_PREFIX}-${OUTPUT}

STEAM_FILTER_PREFIX = filter

SIDECAR_PROJECT_NAME 	= mosn.io/mosn
PROJECT_NAME			= mosn.io/extensions/go-plugin

# only for local test
SIDECAR_GITLAB_PROJECT_NAME = gitlab.alipay-inc.com/ant-mesh/mosn

SIDECAR_BUILD_IMAGE     = mosn-builder
BUILD_IMAGE				= golang:1.14.13
BASE_IMAGE				= zonghaishang/delve:1.6.1

codec:
	docker run --rm \
	-e PLUGIN_TARGET=${TARGET} \
	-e PLUGIN_CODEC_OUTPUT=${CODEC_OUTPUT} \
	-e PLUGIN_CODEC_PREFIX=${CODEC_PREFIX} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-e PLUGIN_STREAM_FILTER=${STREAM_FILTER} \
	-e PLUGIN_TRANSCODER=${TRANSCODER} \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} -it ${BUILD_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/compile-codec.sh

trans: # alias for transcoder (more short)
	make transcoder

transcoder:
	docker run --rm \
	-e PLUGIN_TRANSCODER=${TARGET} \
	-e PLUGIN_TRANSCODER_PREFIX=${TRANSCODER_PREFIX} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/compile-transcoder.sh

filter: # alias for stream-filter (more short)
	make stream-filter

stream-filter:
	docker run --rm \
	-e PLUGIN_STREAM_FILTER=${TARGET} \
	-e PLUGIN_STEAM_FILTER_PREFIX=${STEAM_FILTER_PREFIX} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/compile-filter.sh

compile-codec:
	@rm -rf build/codecs/${PLUGIN_TARGET}
	GO111MODULE=on go build \
		--buildmode=plugin -mod=mod \
		-gcflags "all=-N -l" \
		-o ${PLUGIN_CODEC_OUTPUT} ${PLUGIN_PROJECT_NAME}/plugins/codecs/${PLUGIN_TARGET}/main
	mkdir -p build/codecs/${PLUGIN_TARGET}
	mv ${PLUGIN_CODEC_OUTPUT} build/codecs/${PLUGIN_TARGET}
	@cd build/codecs/${PLUGIN_TARGET} && $(shell which md5sum) -b ${PLUGIN_CODEC_OUTPUT} | cut -d' ' -f1  > ${PLUGIN_CODEC_PREFIX}-${PLUGIN_TARGET}.md5
	cp configs/codecs/${PLUGIN_TARGET}/*.json build/codecs/${PLUGIN_TARGET}

compile-transcoder:
	@rm -rf build/transcoders/${PLUGIN_TARGET}
	GO111MODULE=on go build \
		--buildmode=plugin -mod=mod \
		-gcflags "all=-N -l" \
		-o ${PLUGIN_TRANSCODER_OUTPUT} ${PLUGIN_PROJECT_NAME}/plugins/transcoders/${PLUGIN_TARGET}/main
	mkdir -p build/transcoders/${PLUGIN_TARGET}
	mv ${PLUGIN_TRANSCODER_OUTPUT} build/transcoders/${PLUGIN_TARGET}
	@cd build/transcoders/${PLUGIN_TARGET} && $(shell which md5sum) -b ${PLUGIN_TRANSCODER_OUTPUT} | cut -d' ' -f1  > ${PLUGIN_TRANSCODER_PREFIX}-${PLUGIN_TARGET}.md5
	cp configs/transcoders/${PLUGIN_TARGET}/*.json build/transcoders/${PLUGIN_TARGET}

compile-stream-filter:
	@rm -rf build/stream_filters/${PLUGIN_TARGET}
	GO111MODULE=on go build \
		--buildmode=plugin -mod=mod \
		-gcflags "all=-N -l" \
		-o ${PLUGIN_STEAM_FILTER_OUTPUT} ${PLUGIN_PROJECT_NAME}/plugins/stream_filters/${PLUGIN_TARGET}/main
	mkdir -p build/stream_filters/${PLUGIN_TARGET}
	mv ${PLUGIN_STEAM_FILTER_OUTPUT} build/stream_filters/${PLUGIN_TARGET}
	@cd build/stream_filters/${PLUGIN_TARGET} && $(shell which md5sum) -b ${PLUGIN_STEAM_FILTER_OUTPUT} | cut -d' ' -f1  > ${PLUGIN_STEAM_FILTER_PREFIX}-${PLUGIN_TARGET}.md5
	cp configs/stream_filters/${PLUGIN_TARGET}/*.json build/stream_filters/${PLUGIN_TARGET}


pkg-codec:
	docker run --rm \
	-e PLUGIN_CODEC=${TARGET} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BASE_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/package-codec.sh

pkg-filter:
	docker run --rm \
	-e PLUGIN_STREAM_FILTER=${TARGET} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BASE_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/package-filter.sh

pkg-trans:
	docker run --rm \
	-e PLUGIN_TRANSCODER=${TARGET} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BASE_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/package-transcoder.sh

sidecar:
	bash $(shell go env GOPATH)/src/${PROJECT_NAME}/etc/script/stop.sh
	docker run --rm \
	--env PROJECT_NAME=${PROJECT_NAME} \
	--env SIDECAR_PROJECT_NAME=${SIDECAR_PROJECT_NAME} \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-v $(shell go env GOPATH)/src/${SIDECAR_PROJECT_NAME}:/go/src/${SIDECAR_PROJECT_NAME} \
	-w /go/src/${SIDECAR_PROJECT_NAME} ${BUILD_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/compile.sh

ant:
	bash $(shell go env GOPATH)/src/${PROJECT_NAME}/etc/ant/stop.sh
	docker run --rm \
	--env PROJECT_NAME=${PROJECT_NAME} \
	--env SIDECAR_PROJECT_NAME=${SIDECAR_GITLAB_PROJECT_NAME} \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-v $(shell go env GOPATH)/src/${SIDECAR_GITLAB_PROJECT_NAME}:/go/src/${SIDECAR_GITLAB_PROJECT_NAME} \
	-w /go/src/${SIDECAR_GITLAB_PROJECT_NAME} -it ${BUILD_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/compile.sh

pkg-ant:
	docker run --rm \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BASE_IMAGE} \
	bash /go/src/${PROJECT_NAME}/etc/script/package-ant.sh

start:
	bash $(shell go env GOPATH)/src/${PROJECT_NAME}/etc/ant/start.sh

stop:
	bash $(shell go env GOPATH)/src/${PROJECT_NAME}/etc/ant/stop.sh

debug:
	DLV_DEBUG=true bash $(shell go env GOPATH)/src/${PROJECT_NAME}/etc/ant/start.sh

clean:
	@rm -rf build/codecs build/stream_filters build/transcoders

.PHONY: codec clean start stop debug
