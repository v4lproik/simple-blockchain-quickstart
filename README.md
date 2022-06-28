# Simple-blockchain-quickstart [![CircleCI](https://dl.circleci.com/status-badge/img/gh/v4lproik/simple-blockchain-quickstart/tree/master.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/v4lproik/simple-blockchain-quickstart/tree/master) [![codecov](https://codecov.io/gh/v4lproik/simple-blockchain-quickstart/branch/master/graph/badge.svg?token=LBUG7Y80Q9)](https://codecov.io/gh/v4lproik/simple-blockchain-quickstart) [![Go Report Card](https://goreportcard.com/badge/github.com/v4lproik/simple-blockchain-quickstart)](https://goreportcard.com/report/github.com/v4lproik/simple-blockchain-quickstart) [![api doc](https://badges.aleen42.com/src/apiary.svg)](https://simpleblockchainquickstart.docs.apiary.io/) [![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](.github/ISSUE_TEMPLATE/code_of_conduct.md)

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
export SBQ_IS_JKMS_ACTIVATED="true";
export SBQ_JWT_KEY_PATH="./database/private.pem";
export SBQ_JWT_KEY_ID="sbq-auth-key-id";
export SBQ_JWT_EXPIRES_IN_HOURS="24";
export SBQ_JWT_DOMAIN="localhost";
export SBQ_JWT_AUDIENCE="localhost:8080";
export SBQ_JWT_ISSUER="sbq-local";
export SBQ_JWT_ALGO="HS256";
export SBQ_JWT_JKMS_URL="http://localhost:8080/api/auth/.well-known/jwks.json";
export SBQ_JWT_JKMS_REFRESH_CACHE_INTERVAL_IN_MIN="1";
export SBQ_JWT_JKMS_REFRESH_CACHE_RATE_LIMIT_IN_MIN="1000";
export SBQ_JWT_JKMS_REFRESH_CACHE_TIMEOUT_IN_SEC="1";
```
### Building  
```
make build
```
## Hot reload
1. Install air
```
go install github.com/cosmtrek/air@latest
```
2. (Optional) Export the GOPATH set by goenv if you are using goenv
```
export PATH=$(go env GOPATH)/bin:$PATH
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
1.6557351363071659e+09  info    0xa6aa1c9106f0c0d0895bb72f40cfc830180ebeaf: 1003000
1.6557351363071706e+09  info    0x7b65a12633dbe9a413b17db515732d69e684ebe2: 998000
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
make test -B
```
## Generate doc
1. Install [swag](https://github.com/swaggo/swag)
2. (Optional) Export the GOPATH set by goenv if you are using goenv
```
export PATH=$(go env GOPATH)/bin:$PATH
```
3. Run the swagger
```
swag init
```
### Test data
In the folder ./databases you can find some data that could be used to test the application.  
```
Username: v4lproik  
Account: 0x7b65a12633dbe9a413b17db515732d69e684ebe2  
Password: P@assword-to-access-keystore1
Keystore: databases/keystore/UTC--2022-06-26T13-49-16.552956900Z--7b65a12633dbe9a413b17db515732d69e684ebe2
```
```
Username: cloudvenger  
Account: 0x7b65a12633dbe9a413b17db515732d69e684ebe2  
Password: P@assword-to-access-keystore2  
Keystore: databases/keystore/UTC--2022-06-26T13-50-53.976229800Z--a6aa1c9106f0c0d0895bb72f40cfc830180ebeaf
```