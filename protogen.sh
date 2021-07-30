#!/usr/bin/env bash

set -eo pipefail

ibc_dir=$(go list -f '{{.Dir}}' -m github.com/cosmos/ibc-go)
fabibc_dir=$(go list -f '{{.Dir}}' -m github.com/hyperledger-labs/yui-fabric-ibc)

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  protoc \
  -I "proto" \
  -I "$ibc_dir/proto" \
  -I "$ibc_dir/third_party/proto" \
  -I "$fabibc_dir/proto" \
  -I "$fabibc_dir/third_party/proto" \
  --gocosmos_opt=Mproofs.proto=github.com/confio/ics23/go \
  --gocosmos_out=plugins=interfacetype+grpc,\
Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')
done

cp -r github.com/datachainlab/fabric-ibc-lightclientd/* .
rm -rf github.com/datachainlab/fabric-ibc-lightclientd/*
rmdir -p github.com/datachainlab/fabric-ibc-lightclientd
