VERSION = v0.1.0

BUILD = $(shell git rev-parse --short HEAD)
ROOTPKGNAME = github.com/f4hrenh9it/ro-chess

GO ?= go
GOPATH ?= $(shell ls -d ~/go)

TAGS = nocgo

TARGET_NAME = rch-server
GOCMD = cmd
DOCPORT = 6060
SRVPORT = 3563
DOCSRVPROC = $(shell lsof -t -i:${DOCPORT})
SRVPROC = $(shell lsof -t -i:${SRVPORT})

.PHONY: build generate test cover doc run

generate: build
	${GO} generate ${ROOTPKGNAME}/src/server/game/internal

lint: build
	${GO} get -u github.com/mgechev/revive
	revive -config lint.toml -formatter stylish ./...

test: build generate
	@echo "+ $@"
	${GO} test -race ./...

cover: build generate
	${GO} get -u github.com/dave/courtney
	courtney
	${GO} tool cover -html=coverage.out

doc: build generate
	@echo "+ $@"
	kill ${DOCSRVPROC} || true
	godoc -http :${DOCPORT} &
	sleep 2 # hate to refresh it when it's not ready
	open http://localhost:${DOCPORT}/pkg/${ROOTPKGNAME}

build:
	@echo "+ $@"
	${GO} build -tags "$(TAGS)" -o ${TARGET_NAME} ${GOCMD}/main.go

run: build generate test lint
	@echo "+ $@"
	./${TARGET_NAME}

all: build generate test lint
	kill ${SRVPROC} || true
	./${TARGET_NAME} &
	cd src/front && npm run dev

run-front:
	kill ${SRVPROC} || true
	./${TARGET_NAME} &
	cd src/front && npm run dev