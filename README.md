# Simple-blockchain-quickstart [![CircleCI](https://dl.circleci.com/status-badge/img/gh/v4lproik/simple-blockchain-quickstart/tree/master.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/v4lproik/simple-blockchain-quickstart/tree/master) [![codecov](https://codecov.io/gh/v4lproik/simple-blockchain-quickstart/branch/master/graph/badge.svg?token=LBUG7Y80Q9)](https://codecov.io/gh/v4lproik/simple-blockchain-quickstart) [![Go Report Card](https://goreportcard.com/badge/github.com/v4lproik/simple-blockchain-quickstart)](https://goreportcard.com/report/github.com/v4lproik/simple-blockchain-quickstart)
This is merely a skeleton that helps you quickly set up a simplified version of a blockchain app written in Golang.
## Getting started
### Install Golang & Useful packages
1. Install goenv
2. Install make
3. Install [protoc](https://grpc.io/docs/protoc-installation/)
4. Install golang => 1.18.3
5. Download dependencies
```
make dep
```
### Set env variables  
```
export SBQ_ENV="local"
export SBQ_SERVER_ADDRESS="localhost"
export SBQ_SERVER_PORT="8080"
export SBQ_SERVER_IS_SSL="false"
export SBQ_SERVER_CERT_FILE=""
export SBQ_SERVER_KEY_FILE=""
export SBQ_SERVER_HTTP_CORS_ALLOWED_ORIGINS="http://localhost:8080"
export SBQ_SERVER_HTTP_CORS_ALLOWED_METHODS="GET,POST"
export SBQ_SERVER_HTTP_CORS_ALLOWED_HEADERS=""
```
### Building  
```
make build
```
### Run as client
```
./simple-blockchain-quickstart -g ./databases/genesis.json -d ./databases/blocks.db transaction list
go build -o simple-blockchain-quickstart
1.65573513630384e+09    info    Transactions file: ./databases/blocks.db
1.6557351363038728e+09  info    Genesis file: ./databases/genesis.json
1.6557351363038917e+09  info    Output: console
1.6557351363070958e+09  info    #####################
1.6557351363071187e+09  info    # Accounts balances #
1.6557351363071227e+09  info    #####################
1.6557351363071418e+09  info    State: 87977917793e5fb015311393023ee3ebad19accd1a1c8d7907d58cb686c5ac0a
1.655735136307148e+09   info    ---------------------
1.6557351363071659e+09  info    cloudvenger: 1003000
1.6557351363071706e+09  info    v4lproik: 998000
1.6557351363071866e+09  info    ---------------------
```
### Run as node
```
./simple-blockchain-quickstart -g ./databases/genesis.json -d ./databases/blocks.db -r
1.6558218035227594e+09  info    Transactions file: ./databases/blocks.db
1.6558218035228403e+09  info    Genesis file: ./databases/genesis.json
1.6558218035228572e+09  info    Output: console
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /api/healthz              --> github.com/v4lproik/simple-blockchain-quickstart/domains/healthz.RunDomain.func1 (5 handlers)
1.6558218035261767e+09  info    start server without tls
[GIN-debug] Listening and serving HTTP on 127.0.0.1:8080
```
## Testing
```
make test
```
## TODO
 - Add genesis and transaction files integrity check before launching the app
 - Extend cli commands via grpc calls  
 - ~~Break common components into a common package~~
 - ~~Add tests~~
 - Create database flavour
 - ~~Add custom api errors~~
 - Add error mapping between the commands package error to api error (right now we assume they are all unknown error)
 - Add an env variable which enumerates the functional domains that need to start
 - Add a switch for public/private access through RBAC
 - Add the gas component
 - Add a test seeder and remove the blocks_test.db and genesis_test.db
 - ~~Add better error message when parsing a JSON payload through BindJson~~
 - Add error encapsulation through the different layer in order to deliver better http error code
