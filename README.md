# Simple-blockchain-quickstart [![CircleCI](https://dl.circleci.com/status-badge/img/gh/v4lproik/simple-blockchain-quickstart/tree/master.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/v4lproik/simple-blockchain-quickstart/tree/master) [![codecov](https://codecov.io/gh/v4lproik/simple-blockchain-quickstart/branch/master/graph/badge.svg?token=LBUG7Y80Q9)](https://codecov.io/gh/v4lproik/simple-blockchain-quickstart) [![Go Report Card](https://goreportcard.com/badge/github.com/v4lproik/simple-blockchain-quickstart)](https://goreportcard.com/report/github.com/v4lproik/simple-blockchain-quickstart) [![api doc](https://badges.aleen42.com/src/apiary.svg)](https://simpleblockchainquickstart.docs.apiary.io/) [![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](.github/ISSUE_TEMPLATE/code_of_conduct.md) [![Docker Image Size](https://badgen.net/docker/size/v4lproik/simple-blockchain-quickstart?icon=docker&label=image%20size)](https://hub.docker.com/r/v4lproik/simple-blockchain-quickstart/)

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
cat config/local.env
export SBQ_ENV="local"
export SBQ_SERVER_ADDRESS="localhost"
export SBQ_SERVER_PORT="8080"
export SBQ_SERVER_IS_SSL="false"
export SBQ_SERVER_CERT_FILE=""
export SBQ_SERVER_KEY_FILE=""
export SBQ_SERVER_HTTP_CORS_ALLOWED_ORIGINS="http://localhost:8080"
export SBQ_SERVER_HTTP_CORS_ALLOWED_METHODS="GET,POST"
export SBQ_SERVER_HTTP_CORS_ALLOWED_HEADERS=""
export SBQ_IS_AUTHENTICATION_ACTIVATED="true";
export SBQ_IS_JKMS_ACTIVATED="true";
export SBQ_DOMAINS_TO_START="AUTH,BALANCES,HEALTHZ,NODES,TRANSACTIONS,WALLETS"
export SBQ_JWT_KEY_PATH="./testdata/private.pem";
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
make build && ./simple-blockchain-quickstart -d ./testdata/blocks.db -g ./testdata/genesis.json -k ./testdata/keystore/ -u ./testdata/users.toml transaction list
go build -o simple-blockchain-quickstart
1.656521937350836e+09	info	Transactions file: ./testdata/blocks.db
1.656521937350888e+09	info	Genesis file: ./testdata/genesis.json
1.656521937350891e+09	info	Users file: ./testdata/users.toml
1.656521937350893e+09	info	Keystore dir: ./testdata/keystore/
1.6565219373508952e+09	info	Output: console
1.6565219373573241e+09	info	#####################
1.6565219373573499e+09	info	# Accounts balances #
1.656521937357353e+09	info	#####################
1.656521937357364e+09	info	State: a03d8c9088049b01b25d468919f827a393772f4bcecaf8795f454338c75b6bb2
1.656521937357368e+09	info	Height: 8
1.6565219373573701e+09	info	---------------------
1.656521937357375e+09	info	0x7b65a12633dbe9a413b17db515732d69e684ebe2: 998000
1.656521937357379e+09	info	0xa6aa1c9106f0c0d0895bb72f40cfc830180ebeaf: 1003000
1.6565219373573818e+09	info	---------------------
```
### Run as node
```
./simple-blockchain-quickstart -g ./testdata/genesis.json -d ./testdata/blocks.db -r
1.6558218035227594e+09  info    Transactions file: ./testdata/blocks.db
1.6558218035228403e+09  info    Genesis file: ./testdata/genesis.json
1.6558218035228572e+09  info    Output: console
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /api/healthz              --> github.com/v4lproik/simple-blockchain-quickstart/domains/healthz.RunDomain.func1 (5 handlers)
1.6558218035261767e+09  info    start server without tls
[GIN-debug] Listening and serving HTTP on 127.0.0.1:8080
```
### Run in container
The docker image has been built so the mandatory options are passed in an env file. The extra options are passed through the variable ```cmd```.
To sum up ```cmd``` is responsible for switching from running the app as a client or as a node. The options related to the app itself are stored in ```config/<env>.conf```.
```
#first remove all images locally if you rebuild from this folder
docker-compose down --rmi all
#eg. run as a node for fish shell users
set -l cmd -r; docker-compose -f docker-compose-local.yml up
#eg. run as cli for bash users
cmd="transaction list" docker-compose -f docker-compose-local.yml up
```
## Testing
```
chmod +x ./deployment_script/test.sh
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
## Authentication
JWK authentication is optional, you can activate a verification for each endpoint if necessary. You need to activate it through the environment variable.
```
export SBQ_IS_AUTHENTICATION_ACTIVATED="true";
export SBQ_JWT_KEY_PATH="./testdata/private.pem";
export SBQ_JWT_KEY_ID="sbq-auth-key-id";
export SBQ_JWT_EXPIRES_IN_HOURS="24";
export SBQ_JWT_DOMAIN="localhost";
export SBQ_JWT_AUDIENCE="localhost:8080";
export SBQ_JWT_ISSUER="sbq-local";
export SBQ_JWT_ALGO="HS256";
```
JWT tokens are being signed with the private key passed as an environment variable. Also, you need to provide the url containing the JKMS derived from the private key.
You can use this project JKMS service exposing the public key parameters needed to verify a JWT token.
```
export SBQ_IS_JKMS_ACTIVATED="true";
export SBQ_JWT_JKMS_REFRESH_CACHE_INTERVAL_IN_MIN="1";
export SBQ_JWT_JKMS_REFRESH_CACHE_RATE_LIMIT_IN_MIN="1000";
export SBQ_JWT_JKMS_REFRESH_CACHE_TIMEOUT_IN_SEC="1";
```
The users are declared in ```./testdata/users.toml```. See the Test data section for the test accounts.
```
curl localhost:8080/api/balances/ -X POST                                                                                                                          15:03:11
{"error":{"code":401,"status":"Unauthorized","message":"authentication token cannot be found","context":[]}}

> curl localhost:8080/api/auth/login -X POST -d '{"username": "v4lproik", "password":"P@assword-to-access-api1"}' -H 'Content-type: application/json'
{"access_token":"eyJhbGciOiJSUzI1NiIsImtpZCI6InNicS1hdXRoLWtleS1pZCIsInR5cCI6IkpXVCJ9.eyJkYXQiOnsiTmFtZSI6InY0bHByb2lrIiwiSGFzaCI6IiRhcmdvbjJpZCR2PTE5JG09NjU1MzYsdD0zLHA9MiRGdVNVWlEwbXJUTTl1SXBQOHFwTlV3JG5kMjFqdGdVWmpKanowNzhqZGxTREt4cWFqdjVwYWl4bG9HR05nVE1KSXcifSwiZXhwIjoxNjU2NTk0MDE1LCJpYXQiOjE2NTY1MDc2MTUsIm5iZiI6MTY1NjUwNzYxNX0.4VBqD9Cg2KH96CioyRtSIlM2edGneXxZLrxG46Qub4Pol-NWOXI9_PAmIL_DmQEvF95x44m9Vl8VF2RZdO42B03cxKZPKIzjZjalHqyEl3YPyz27kP7d_YCCMjSzKMbx8Np7u9orWjlC5MayCB2rtgefag3DkKGJWUAIH5OfDPy6B-XLsgL8caWN0aM4TCelC-geo2bC488Xk79YffhfNLJPuvgKuuUeWaWLz-YHcALbguqRP_ehqDvn5vzBBWAS_aCYN3W9-dsOHttfSRKaxmxQm-hxcp01T7ezXgNO3gnJmfuWff-96UKZVb0QPzG1ltPWInqheKRypviuAEIHUg"}

> curl localhost:8080/api/balances/ -X POST -H "X-API-TOKEN: eyJhbGciOiJSUzI1NiIsImtpZCI6InNicS1hdXRoLWtleS1pZCIsInR5cCI6IkpXVCJ9.eyJkYXQiOnsiTmFtZSI6InY0bHByb2lrIiwiSGFzaCI6IiRhcmdvbjJpZCR2PTE5JG09NjU1MzYsdD0zLHA9MiRGdVNVWlEwbXJUTTl1SXBQOHFwTlV3JG5kMjFqdGdVWmpKanowNzhqZGxTREt4cWFqdjVwYWl4bG9HR05nVE1KSXcifSwiZXhwIjoxNjU2NTk0MDE1LCJpYXQiOjE2NTY1MDc2MTUsIm5iZiI6MTY1NjUwNzYxNX0.4VBqD9Cg2KH96CioyRtSIlM2edGneXxZLrxG46Qub4Pol-NWOXI9_PAmIL_DmQEvF95x44m9Vl8VF2RZdO42B03cxKZPKIzjZjalHqyEl3YPyz27kP7d_YCCMjSzKMbx8Np7u9orWjlC5MayCB2rtgefag3DkKGJWUAIH5OfDPy6B-XLsgL8caWN0aM4TCelC-geo2bC488Xk79YffhfNLJPuvgKuuUeWaWLz-YHcALbguqRP_ehqDvn5vzBBWAS_aCYN3W9-dsOHttfSRKaxmxQm-hxcp01T7ezXgNO3gnJmfuWff-96UKZVb0QPzG1ltPWInqheKRypviuAEIHUg"
{"balances":[{"account":"0x7b65a12633dbe9a413b17db515732d69e684ebe2","value":998000},{"account":"0xa6aa1c9106f0c0d0895bb72f40cfc830180ebeaf","value":1003000}]}
```
### Test data
In the folder ./databases you can find some data that could be used to test the application.  
```
Username: v4lproik
Password: P@assword-to-access-api1
Hash    : $argon2id$v=19$m=65536,t=3,p=2$FuSUZQ0mrTM9uIpP8qpNUw$nd21jtgUZjJjz078jdlSDKxqajv5paixloGGNgTMJIw

Account : 0x7b65a12633dbe9a413b17db515732d69e684ebe2
Password: P@assword-to-access-keystore1
Keystore: testdata/keystore/UTC--2022-06-26T13-49-16.552956900Z--7b65a12633dbe9a413b17db515732d69e684ebe2
```
```
Username: cloudvenger
Password: P@assword-to-access-api2
Hash    : $argon2id$v=19$m=65536,t=3,p=2$j2yd8FWqhApKrrqmkkLMQA$Lfh/7K+oP3IWdTrQSjURBS6PFttzlksmozz8kuGBCqk

Account : 0x7b65a12633dbe9a413b17db515732d69e684ebe2
Password: P@assword-to-access-keystore2
Keystore: testdata/keystore/UTC--2022-06-26T13-50-53.976229800Z--a6aa1c9106f0c0d0895bb72f40cfc830180ebeaf
```
