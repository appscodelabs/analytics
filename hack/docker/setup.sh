#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

LIB_ROOT=$(dirname "${BASH_SOURCE}")/..
source "$LIB_ROOT/libbuild/common/lib.sh"
source "$LIB_ROOT/libbuild/common/public_image.sh"

GOPATH=$(go env GOPATH)
SRC=$GOPATH/src
BIN=$GOPATH/bin
ROOT=$GOPATH

APPSCODE_ENV=${APPSCODE_ENV:-dev}

IMG=client-ip
# TAG=0.1.0
if [ -f "$GOPATH/src/github.com/appscode/client-ip/dist/.tag" ]; then
	export $(cat $GOPATH/src/github.com/appscode/client-ip/dist/.tag | xargs)
fi

build_binary() {
	pushd $GOPATH/src/github.com/appscode/client-ip
	./hack/builddeps.sh
    ./hack/make.py build
	detect_tag $GOPATH/src/github.com/appscode/client-ip/dist/.tag
	popd
}

build_docker() {
	pushd $GOPATH/src/github.com/appscode/client-ip/hack/docker
	cp $GOPATH/src/github.com/appscode/client-ip/dist/client-ip/client-ip-linux-amd64 client-ip
	chmod 755 client-ip

	cat >Dockerfile <<EOL
FROM alpine

COPY client-ip /client-ip

USER nobody:nobody
ENTRYPOINT ["/client-ip"]
EOL
	local cmd="docker build -t appscode/$IMG:$TAG ."
	echo $cmd; $cmd

	rm client-ip Dockerfile
	popd
}

build() {
	build_binary
	build_docker
}

source_repo $@
