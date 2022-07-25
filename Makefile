BINARY=simple-blockchain-quickstart
HOT_RELOAD_BIN=air

SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

#variables specific to node1
SBQ_GENESIS_FILE_PATH_1 = ${SBQ_GENESIS_FILE_PATH_NODE_1}
SBQ_TRANSACTIONS_FILE_PATH_1 = ${SBQ_TRANSACTIONS_FILE_PATH_NODE_1}
SBQ_KEYSTORE_DIR_PATH_1 = ${SBQ_KEYSTORE_DIR_PATH_NODE_1}
SBQ_USERS_FILE_PATH_1 = ${SBQ_USERS_FILE_PATH_NODE_1}
SBQ_NODES_FILE_PATH_1 = ${SBQ_NODES_FILE_PATH_NODE_1}
SBQ_NODE_MINER_ADDRESS_1 = ${SBQ_NODE_MINER_ADDRESS_NODE_1}

#variables specific to node2
SBQ_GENESIS_FILE_PATH_2 = ${SBQ_GENESIS_FILE_PATH_NODE_2}
SBQ_TRANSACTIONS_FILE_PATH_2 = ${SBQ_TRANSACTIONS_FILE_PATH_NODE_2}
SBQ_KEYSTORE_DIR_PATH_2 = ${SBQ_KEYSTORE_DIR_PATH_NODE_2}
SBQ_USERS_FILE_PATH_2 = ${SBQ_USERS_FILE_PATH_NODE_2}
SBQ_NODES_FILE_PATH_2 = ${SBQ_NODES_FILE_PATH_NODE_2}
SBQ_NODE_MINER_ADDRESS_2 = ${SBQ_NODE_MINER_ADDRESS_NODE_2}

#useful in dockerfile
print-%  : ; @echo $* = $($*)

default: build

proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto

server:
	./bin/${BINARY} -g ${SBQ_GENESIS_FILE_PATH_1} \
	-d ${SBQ_TRANSACTIONS_FILE_PATH_1} \
	-k ${SBQ_KEYSTORE_DIR_PATH_1} \
	-u ${SBQ_USERS_FILE_PATH_1} \
	-n ${SBQ_NODES_FILE_PATH_1} \
	-m ${SBQ_NODE_MINER_ADDRESS_1} \
	-r

server_hot_reload_1:
	${HOT_RELOAD_BIN} -- -g ${SBQ_GENESIS_FILE_PATH_1} \
	-d ${SBQ_TRANSACTIONS_FILE_PATH_1} \
	-k ${SBQ_KEYSTORE_DIR_PATH_1} \
	-u ${SBQ_USERS_FILE_PATH_1} \
	-n ${SBQ_NODES_FILE_PATH_1} \
	-m ${SBQ_NODE_MINER_ADDRESS_1} \
	-r

server_hot_reload_2:
	${HOT_RELOAD_BIN} -- -g ${SBQ_GENESIS_FILE_PATH_2} \
	-d ${SBQ_TRANSACTIONS_FILE_PATH_2} \
	-k ${SBQ_KEYSTORE_DIR_PATH_2} \
	-u ${SBQ_USERS_FILE_PATH_2} \
	-n ${SBQ_NODES_FILE_PATH_2} \
	-m ${SBQ_NODE_MINER_ADDRESS_2} \
	-r

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
	mkdir -p coverage
	go test -coverprofile=coverage/coverage.out -covermode=count ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

test-coverage-install:
	go get github.com/dave/courtney
	go install github.com/dave/courtney

test-coverage:
	mkdir -p coverage
	courtney -v -o coverage/coverage.out ./...

vet:
	go vet

fmt:
	@gofmt -l -w $(SRC)

format:
	gofumpt -l -w .