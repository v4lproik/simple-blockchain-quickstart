BINARY=simple-blockchain-quickstart

SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

default: build

build:
	go build -o $(BINARY)

clean:
	go clean -i -v -r
	rm -f $(BINARY)

test:
	@if [ -f ${TEST} ] ; then ./${TEST} ; fi

fmt:
	@gofmt -l -w $(SRC)