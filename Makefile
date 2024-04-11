COMMIT:=$(shell git rev-parse --short HEAD)
VERSION:=$(shell date -u '+%Y/%m/%dT%H:%M:%S')
SERVER:=sintep@192.168.254.107:/home/sintep/iRo/

.DEFAULT_GOAL:=all

.PHONY: all
all: check build install

.PHONY: build
build:
	@echo -e "\033[94m# ${VERSION}V${COMMIT} build"
	@go build -v -mod vendor -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT}" ./cmd/iRo
	@/usr/bin/size iRo
	@echo -e "#\033[39m"

.PHONY: test
test:
	@echo -e "\033[96m# ${VERSION}V${COMMIT} test"
	@echo -e "#\033[39m"
#	go test ./...

.PHONY: check
check:
	@echo -e "\033[31m# ${VERSION}V${COMMIT} check"
	@golangci-lint run ./cmd/iRo
	@echo -e "#\033[39m"

.PHONY: install
install:
	@echo -e "\033[92m# ${VERSION}V${COMMIT} install"
	@scp iRo $(SERVER)
	@scp -r config web script $(SERVER)
	@echo -e "#\033[39m"

.PHONY: mod
mod:
	@echo -e "\033[35m# ${VERSION}V${COMMIT} mod"
	@go mod vendor
	@echo -e "#\033[39m"
