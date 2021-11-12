SHELL			= /bin/bash
TARGET			= ${plugin}
OUTPUT			= ${TARGET}.so

CODEC_PREFIX	= codec
CODEC_OUTPUT	= ${CODEC_PREFIX}-${OUTPUT}

PROJECT_NAME	= github.com/mosn/wasm-sdk/go-plugin

clean:
	@rm -rf build/*

codec:
	@rm -rf build/codecs/${TARGET}
	go build -o ${CODEC_OUTPUT} -mod=mod --buildmode=plugin -gcflags "all=-N -l" ${PROJECT_NAME}/plugins/codecs/${TARGET}/main
	mkdir -p build/codecs/${TARGET}
	mv ${CODEC_OUTPUT} build/codecs/${TARGET}
	@cd build/codecs/${TARGET} && $(shell which md5sum) -b ${CODEC_OUTPUT} | cut -d' ' -f1  > ${CODEC_PREFIX}-${TARGET}.md5
	cp configs/codecs/${TARGET}/*.json build/codecs/${TARGET}

.PHONY: codec clean