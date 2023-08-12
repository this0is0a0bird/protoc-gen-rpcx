#/bin/bash


APPLICATION=protoc-gen-rpcx
CMD=main.go rpcx.go

# These are the values we want to pass for Version and BuildTime
GITTAG=1.0.0
BUILD_TIME=`date +%Y%m%d%H%M%S`
# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags="-X main.Version=${GITTAG} -X main.Build_Time=${BUILD_TIME} -s -w"
# GCFLAGS=-gcflags="-d=ssa/check_bce/debug=1"
GCFLAGS=-gcflags="all=-N -l" ## enable debug

default: macos

release: macos

macos:
	GO111MODULE=on	go build -v ${LDFLAGS} ${GCFLAGS} -o build/bin/${APPLICATION} ${CMD}
	-cp -f build/bin/protoc-gen-rpcx ../../bin/protoc-gen-rpcx
	-protoc -I . --go_out=paths=source_relative:. --rpcx_out=. --rpcx_opt=paths=source_relative helloworld.proto
	-go mod tidy
	-go mod vendor

linux:
	export GOPROXY=https://goproxy.io
	export GOPRIVATE=github.com/cctip
	GO111MODULE=on GOOS=linux  go   build ${LDFLAGS} ${GCFLAGS} -o build/bin/${APPLICATION} ${CMD}
	-cp -r cmd/etc build/bin
	-tar -czvf build/${APPLICATION}.tar.gz build/bin
	-rm -r build/bin

rpc:
	protoc -I . --go_out=paths=source_relative:. --rpcx_out=. --rpcx_opt=paths=source_relative helloworld.proto

clean:
	-rm -rf build

