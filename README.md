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
sh config/local.env
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
make build && ./bin/simple-blockchain-quickstart -d ./testdata/node1/blocks.db -g ./testdata/node1/genesis.json -k ./testdata/node1/keystore/ -u ./testdata/node1/users.toml -n ./testdata/node1/network_nodes.toml -r 19:51:41
go build -o ./bin/simple-blockchain-quickstart
1.657907504219829e+09	info	Transactions file: ./testdata/node1/blocks.db
1.6579075042198799e+09	info	Genesis file: ./testdata/node1/genesis.json
1.657907504219883e+09	info	Users file: ./testdata/node1/users.toml
1.657907504219885e+09	info	Nodes file: ./testdata/node1/network_nodes.toml
1.657907504219887e+09	info	Keystore dir: ./testdata/node1/keystore/
1.657907504219889e+09	info	Output: console
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
- using env:	export GIN_MODE=release
- using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /api/auth/login           --> github.com/v4lproik/simple-blockchain-quickstart/domains/auth.AuthEnv.Login-fm (5 handlers)
[GIN-debug] GET    /api/auth/.well-known/jwks.json --> github.com/v4lproik/gin-jwks-rsa.Jkws.func1 (5 handlers)
[GIN-debug] POST   /api/balances/            --> github.com/v4lproik/simple-blockchain-quickstart/domains/balances.(*BalancesEnv).ListBalances-fm (6 handlers)
[GIN-debug] GET    /api/healthz              --> github.com/v4lproik/simple-blockchain-quickstart/domains/healthz.RunDomain.func1 (5 handlers)
[GIN-debug] GET    /api/nodes/status         --> github.com/v4lproik/simple-blockchain-quickstart/domains/nodes.NodesEnv.NodeStatus-fm (5 handlers)
[GIN-debug] POST   /api/nodes/blocks         --> github.com/v4lproik/simple-blockchain-quickstart/domains/nodes.NodesEnv.NodeListBlocks-fm (5 handlers)
[GIN-debug] PUT    /api/transactions/        --> github.com/v4lproik/simple-blockchain-quickstart/domains/transactions.TransactionsEnv.AddTransaction-fm (6 handlers)
[GIN-debug] PUT    /api/wallets/             --> github.com/v4lproik/simple-blockchain-quickstart/domains/wallets.(*WalletsEnv).CreateWallet-fm (6 handlers)
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
export SBQ_JWT_KEY_PATH="./testdata/node1/private.pem";
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
The users are declared in ```./testdata/node1/users.toml```. See the Test data section for the test accounts.
```
curl localhost:8080/api/balances/ -X POST                                                                                                                          15:03:11
{"error":{"code":401,"status":"Unauthorized","message":"authentication token cannot be found","context":[]}}

> curl localhost:8080/api/auth/login -X POST -d '{"username": "v4lproik", "password":"P@assword-to-access-api1"}' -H 'Content-type: application/json'
{"access_token":"eyJhbGciOiJSUzI1NiIsImtpZCI6InNicS1hdXRoLWtleS1pZCIsInR5cCI6IkpXVCJ9.eyJkYXQiOnsiTmFtZSI6InY0bHByb2lrIiwiSGFzaCI6IiRhcmdvbjJpZCR2PTE5JG09NjU1MzYsdD0zLHA9MiRGdVNVWlEwbXJUTTl1SXBQOHFwTlV3JG5kMjFqdGdVWmpKanowNzhqZGxTREt4cWFqdjVwYWl4bG9HR05nVE1KSXcifSwiZXhwIjoxNjU2NTk0MDE1LCJpYXQiOjE2NTY1MDc2MTUsIm5iZiI6MTY1NjUwNzYxNX0.4VBqD9Cg2KH96CioyRtSIlM2edGneXxZLrxG46Qub4Pol-NWOXI9_PAmIL_DmQEvF95x44m9Vl8VF2RZdO42B03cxKZPKIzjZjalHqyEl3YPyz27kP7d_YCCMjSzKMbx8Np7u9orWjlC5MayCB2rtgefag3DkKGJWUAIH5OfDPy6B-XLsgL8caWN0aM4TCelC-geo2bC488Xk79YffhfNLJPuvgKuuUeWaWLz-YHcALbguqRP_ehqDvn5vzBBWAS_aCYN3W9-dsOHttfSRKaxmxQm-hxcp01T7ezXgNO3gnJmfuWff-96UKZVb0QPzG1ltPWInqheKRypviuAEIHUg"}

> curl localhost:8080/api/balances/ -X POST -H "X-API-TOKEN: eyJhbGciOiJSUzI1NiIsImtpZCI6InNicS1hdXRoLWtleS1pZCIsInR5cCI6IkpXVCJ9.eyJkYXQiOnsiTmFtZSI6InY0bHByb2lrIiwiSGFzaCI6IiRhcmdvbjJpZCR2PTE5JG09NjU1MzYsdD0zLHA9MiRGdVNVWlEwbXJUTTl1SXBQOHFwTlV3JG5kMjFqdGdVWmpKanowNzhqZGxTREt4cWFqdjVwYWl4bG9HR05nVE1KSXcifSwiZXhwIjoxNjU2NTk0MDE1LCJpYXQiOjE2NTY1MDc2MTUsIm5iZiI6MTY1NjUwNzYxNX0.4VBqD9Cg2KH96CioyRtSIlM2edGneXxZLrxG46Qub4Pol-NWOXI9_PAmIL_DmQEvF95x44m9Vl8VF2RZdO42B03cxKZPKIzjZjalHqyEl3YPyz27kP7d_YCCMjSzKMbx8Np7u9orWjlC5MayCB2rtgefag3DkKGJWUAIH5OfDPy6B-XLsgL8caWN0aM4TCelC-geo2bC488Xk79YffhfNLJPuvgKuuUeWaWLz-YHcALbguqRP_ehqDvn5vzBBWAS_aCYN3W9-dsOHttfSRKaxmxQm-hxcp01T7ezXgNO3gnJmfuWff-96UKZVb0QPzG1ltPWInqheKRypviuAEIHUg"
{"balances":[{"account":"0x7b65a12633dbe9a413b17db515732d69e684ebe2","value":998000},{"account":"0xa6aa1c9106f0c0d0895bb72f40cfc830180ebeaf","value":1003000}]}
```
### Test data
In the folder ```./testdata/node*/``` you can find some data that could be used to test the application.  
```
Username: v4lproik
Password: P@assword-to-access-api1
Hash    : $argon2id$v=19$m=65536,t=3,p=2$FuSUZQ0mrTM9uIpP8qpNUw$nd21jtgUZjJjz078jdlSDKxqajv5paixloGGNgTMJIw

Account : 0x7b65a12633dbe9a413b17db515732d69e684ebe2
Password: P@assword-to-access-keystore1
Keystore: testdata/node1/keystore/UTC--2022-06-26T13-49-16.552956900Z--7b65a12633dbe9a413b17db515732d69e684ebe2
```
```
Username: cloudvenger
Password: P@assword-to-access-api2
Hash    : $argon2id$v=19$m=65536,t=3,p=2$j2yd8FWqhApKrrqmkkLMQA$Lfh/7K+oP3IWdTrQSjURBS6PFttzlksmozz8kuGBCqk

Account : 0x7b65a12633dbe9a413b17db515732d69e684ebe2
Password: P@assword-to-access-keystore2
Keystore: testdata/node1/keystore/UTC--2022-06-26T13-50-53.976229800Z--a6aa1c9106f0c0d0895bb72f40cfc830180ebeaf
```
