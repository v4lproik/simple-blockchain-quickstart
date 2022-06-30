BINARY=simple-blockchain-quickstart
HOT_RELOAD_BIN=air

SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
TEST=deployment_script/test.sh

default: build

proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto

server:
	./${BINARY} -g ./databases/genesis.json -d ./databases/blocks.db -k ./databases/keystore/ -u ./databases/users.toml -n ./databases/network_nodes.toml -r

server_hot_reload:
	${HOT_RELOAD_BIN} -- -g ./databases/genesis.json -d ./databases/blocks.db -k ./databases/keystore/ -u ./databases/users.toml -n ./databases/network_nodes.toml -r

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