# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=go-vault-session
CMD_PATH=./server/cmd/server/main.go

build-local-mac:
	GOOS=darwin GOARCH=amd64 ${GOBUILD} -o ./${BINARY_NAME} -v ${CMD_PATH}

build-local:
	${GOBUILD} -o ./${BINARY_NAME} -v ${CMD_PATH}

test:
	docker-compose up -d
	./run-tests.sh
	docker-compose down
