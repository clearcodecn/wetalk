#!/bin/bash

set -ex

install_protoc () {
  OS=`go env GOOS`
  Version="3.9.1"
  if [[ "$OS" = "darwin" ]]; then
    OS="osx"
  fi
  URL=" https://github.com/protocolbuffers/protobuf/releases/download/v${Version}/protoc-${Version}-${OS}-x86_64.zip"
  wget $URL -O protoc.zip
  mkdir -p .tmp/protoc
  unzip protoc.zip -d ./.tmp/protoc
  mkdir bin
  mv .tmp/protoc/bin/protoc bin/protoc
  mv .tmp/protoc/include bin/include
  rm -rf .tmp

}

if [[ ! -f "bin/protoc" ]] ; then
  install_protoc
fi

if [[ ! -f "bin/protoc-gen-go" ]] ; then

  export GOBIN=`pwd`/bin
  echo $GOBIN
  export GO111MODULE=on
  export GOPROXY=https://goproxy.io

  go install -v github.com/golang/protobuf/protoc-gen-go
  go install -v github.com/gogo/protobuf/protoc-gen-gogofast
  go install -v github.com/gogo/protobuf/protoc-gen-gogoslick
  go install -v github.com/golang/mock/mockgen
fi

imports=(
  "proto"
  "./third_party/protobuf"
  "./bin/include"
)

protoc="./bin/protoc"
protocarg=""
for i in "${imports[@]}"
do
  protocarg+="--proto_path=$i "
done

mappings=(
  "google/api/annotations.proto=istio.io/gogo-genproto/googleapis/google/api"
  "google/api/http.proto=istio.io/gogo-genproto/googleapis/google/api"
  "google/rpc/code.proto=istio.io/gogo-genproto/googleapis/google/rpc"
  "google/rpc/error_details.proto=istio.io/gogo-genproto/googleapis/google/rpc"
  "google/rpc/status.proto=istio.io/gogo-genproto/googleapis/google/rpc"
  "google/protobuf/any.proto=github.com/gogo/protobuf/types"
  "google/protobuf/duration.proto=github.com/gogo/protobuf/types"
  "google/protobuf/empty.proto=github.com/gogo/protobuf/types"
  "google/protobuf/struct.proto=github.com/gogo/protobuf/types"
  "google/protobuf/timestamp.proto=github.com/gogo/protobuf/types"
  "google/protobuf/wrappers.proto=github.com/gogo/protobuf/types"
  "gogoproto/gogo.proto=github.com/gogo/protobuf/gogoproto"
)

gogoarg="plugins=grpc"
# assign importmap for canonical protos
for mapping in "${mappings[@]}"
do
  gogoarg+=",M$mapping"
done

$protoc ${protocarg} ./proto/*.proto \
      --plugin=protoc-gen-go=./bin/protoc-gen-go \
      --plugin=protoc-gen-gogoslick=./bin/protoc-gen-gogoslick --gogoslick_out=${gogoarg}:./proto