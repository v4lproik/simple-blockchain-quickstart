BINARY=simple-blockchain-quickstart

SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
TEST=deployment_script/test.sh

default: build

dep:
	go mod download

build:
	go build -o ${BINARY}

run:
	./${BINARY}

clean:
	go clean -i -v -r
	rm ${BINARY}

test:
	@if [ -f ${TEST} ] ; then ./${TEST} ; fi

vet:
	go vet

fmt:
	@gofmt -l -w $(SRC)