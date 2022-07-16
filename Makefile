BINARY=simple-blockchain-quickstart
HOT_RELOAD_BIN=air

SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
TEST=deployment_script/test.sh

#useful in dockerfile
print-%  : ; @echo $* = $($*)

default: build

proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto

server:
	./bin/${BINARY} -g ./testdata/node1/genesis.json -d ./testdata/node1/blocks.db -k ./testdata/node1/keystore/ -u ./testdata/node1/users.toml -n ./testdata/node1/network_nodes.toml -r

server_hot_reload_1:
	${HOT_RELOAD_BIN} -- -g ./testdata/node1/genesis.json -d ./testdata/node1/blocks.db -k ./testdata/node1/keystore/ -u ./testdata/node1/users.toml -n ./testdata/node1/network_nodes.toml -r

server_hot_reload_2:
	${HOT_RELOAD_BIN} -- -g ./testdata/node2/genesis.json -d ./testdata/node2/blocks.db -k ./testdata/node2/keystore/ -u ./testdata/node2/users.toml -n ./testdata/node2/network_nodes.toml -r

dep:
	go mod download

build:
	go build -o ./bin/${BINARY}

run:
	./bin/${BINARY}

clean:
	go clean -i -v -r
	rm ./bin/${BINARY}

test:
	@if [ -f ${TEST} ] ; then ./${TEST} ; fi

vet:
	go vet

fmt:
	@gofmt -l -w $(SRC)