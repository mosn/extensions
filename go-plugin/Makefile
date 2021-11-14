SHELL			= /bin/bash
TARGET			= ${plugin}
OUTPUT			= ${TARGET}.so

CODEC_PREFIX	= codec
CODEC_OUTPUT	= ${CODEC_PREFIX}-${OUTPUT}

SIDECAR_PROJECT_NAME 	= mosn.io/mosn
PROJECT_NAME			= github.com/mosn/wasm-sdk/go-plugin

SIDECAR_BUILD_IMAGE     = mosn-builder
BUILD_IMAGE				= golang:1.14.13

codec:
	docker run --rm \
	-e PLUGIN_TARGET=${TARGET} \
	-e PLUGIN_CODEC_OUTPUT=${CODEC_OUTPUT} \
	-e PLUGIN_CODEC_PREFIX=${CODEC_PREFIX} \
	-e PLUGIN_PROJECT_NAME=${PROJECT_NAME} \
	-v $(shell go env GOPATH):/go \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} \
	make compile-codec

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

sidecar:
	docker run --rm \
	--env PROJECT_NAME=${PROJECT_NAME} \
	--env SIDECAR_PROJECT_NAME=${SIDECAR_PROJECT_NAME} \
	-v $(shell go env GOPATH):/go \
	-v $(shell pwd):/go/src/${PROJECT_NAME} \
	-v $(shell go env GOPATH)/src/${SIDECAR_PROJECT_NAME}:/go/src/${SIDECAR_PROJECT_NAME} \
	-w /go/src/${SIDECAR_PROJECT_NAME} ${BUILD_IMAGE} \
	sh /go/src/${PROJECT_NAME}/etc/script/compile.sh

clean:
	@rm -rf build/*

.PHONY: codec clean