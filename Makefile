BINARY=simple-blockchain-quickstart

SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
TEST=deployment_script/test.sh

default: build

proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto

server:
	./${BINARY} -g ./databases/genesis.json -d ./databases/blocks.db -k ./databases/keystore/ -r

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