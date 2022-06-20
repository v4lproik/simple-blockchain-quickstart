BINARY=simple-blockchain-quickstart

SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

default: build

build:
	 GOARCH=amd64 GOOS=darwin go build -o ${BINARY}-darwin main.go
     GOARCH=amd64 GOOS=linux go build -o ${BINARY}-linux main.go
     GOARCH=amd64 GOOS=window go build -o ${BINARY}-windows main.go

run:
	./${BINARY}

clean:
	go clean -i -v -r
	rm ${BINARY}-darwin
	rm ${BINARY}-linux
	rm ${BINARY}-windows

test:
	@if [ -f ${TEST} ] ; then ./${TEST} ; fi

fmt:
	@gofmt -l -w $(SRC)