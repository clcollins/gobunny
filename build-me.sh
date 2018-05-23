#!/usr/bin/env bash

# Generic Go app image builder
# Gets its appname from the directory name

set -o errexit

appname="$(basename $(pwd))"
bin_prefix="github.com/clcollins"
bindir="/go/src/${bin_prefix}/${appname}"

builder="${appname}_builder"
builder_image="${appname}:builder"
builder_dockerfile="Dockerfile-builder"

app_image="${appname}:latest"
app_dockerfile="Dockerfile"

function cleanup() {
  echo "Cleaning up..."
  docker rm $builder
}

# Always cleanup on exit
trap cleanup EXIT

docker build --tag ${builder_image} \
  --file ${builder_dockerfile} .

docker create --name ${builder} ${builder_image}

mkdir -p ./pkg
docker cp ${builder}:${bindir}/${appname} ./pkg

docker build --tag ${app_image} \
  --file ${app_dockerfile} .

