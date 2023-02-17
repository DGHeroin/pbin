NAME=pbs
pwd=$(shell pwd)
GO_BUILD=go build
RELEASE_DIR?=release

$(shell mkdir -p ${RELEASE_DIR})

PHONEY:

all: linux-amd64  darwin-amd64 windows-amd64

linux-amd64:
	cd $(pwd)/cmd/pbs/;\
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	$(GO_BUILD) -o $(pwd)/${RELEASE_DIR}/$(NAME)-linux-amd64 -v -trimpath

linux-arm64:
	cd $(pwd)/cmd/pbs/;\
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
	$(GO_BUILD) -o $(pwd)/${RELEASE_DIR}/$(NAME)-linux-arm64 -v -trimpath

darwin-amd64:
	cd $(pwd)/cmd/pbs/;\
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
	$(GO_BUILD) -o $(pwd)/${RELEASE_DIR}/$(NAME)-darwin-amd64 -v -trimpath

windows-amd64:
	cd $(pwd)/cmd/pbs/;\
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
	$(GO_BUILD) -o $(pwd)/${RELEASE_DIR}/$(NAME)-windows-amd64.exe -v -trimpath