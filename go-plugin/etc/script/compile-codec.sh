#!/bin/bash

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOPRIVATE=gitlab.alipay-inc.com,code.alipay.com

build_opts=""
if [[ -n ${PLUGIN_OS} && -n ${PLUGIN_ARCH} ]]; then
  build_opts="GOOS=${PLUGIN_OS} GOARCH=${PLUGIN_ARCH}"
  echo "compiling codec ${PLUGIN_TARGET} for ${PLUGIN_OS} ${PLUGIN_ARCH} ..."
else
  echo "compiling codec ${PLUGIN_TARGET} for linux $(dpkg --print-architecture) ..."
fi

export BUILD_OPTS=${build_opts}
make compile-codec

# build stream filter
if [[ -n "${PLUGIN_STREAM_FILTER}" ]]; then
  bash /go/src/"${PLUGIN_PROJECT_NAME}"/etc/script/compile-filter.sh
fi

# copy stream filter
if [[ -n "${PLUGIN_STREAM_FILTER}" ]]; then
  filters=(${PLUGIN_STREAM_FILTER//,/ })
  rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/codecs/${PLUGIN_TARGET}/stream_filters
  meta=
  for name in "${filters[@]}"; do
    mkdir -p /go/src/${PLUGIN_PROJECT_NAME}/build/codecs/${PLUGIN_TARGET}/stream_filters/${name}
    echo "cp  /go/src/${PLUGIN_PROJECT_NAME}/build/stream_filters/${name} /go/src/${PLUGIN_PROJECT_NAME}/build/codecs/${PLUGIN_TARGET}/stream_filters/"
    cp -r /go/src/${PLUGIN_PROJECT_NAME}/build/stream_filters/${name} \
      /go/src/${PLUGIN_PROJECT_NAME}/build/codecs/${PLUGIN_TARGET}/stream_filters/
    # append ,
    if [[ -n ${meta} ]]; then
      meta+=","
    fi
    meta+="\"${name}\""
  done

  # write metadata
  echo "{\"plugins\":[${meta}]}" >/go/src/${PLUGIN_PROJECT_NAME}/build/codecs/${PLUGIN_TARGET}/stream_filters/plugin-meta.json
fi

# build transcoder
if [[ -n "${PLUGIN_TRANSCODER}" ]]; then
  bash /go/src/"${PLUGIN_PROJECT_NAME}"/etc/script/compile-transcoder.sh
fi

# copy transcoder
if [[ ${PLUGIN_TRANSCODER} != "" ]]; then
  coders=(${PLUGIN_TRANSCODER//,/ })
  rm -rf /go/src/${PLUGIN_PROJECT_NAME}/build/codecs/${PLUGIN_TARGET}/transcoders
  for name in "${coders[@]}"; do
    mkdir -p /go/src/${PLUGIN_PROJECT_NAME}/build/codecs/${PLUGIN_TARGET}/transcoders/${name}/
    echo "cp  /go/src/${PLUGIN_PROJECT_NAME}/build/transcoders/${name} /go/src/${PLUGIN_PROJECT_NAME}/build/codecs/${PLUGIN_TARGET}/transcoders/"
    cp -r /go/src/${PLUGIN_PROJECT_NAME}/build/transcoders/${name} \
      /go/src/${PLUGIN_PROJECT_NAME}/build/codecs/${PLUGIN_TARGET}/transcoders/
  done
fi
