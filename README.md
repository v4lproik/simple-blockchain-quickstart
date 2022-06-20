# Simple-blockchain-quickstart [![CircleCI](https://dl.circleci.com/status-badge/img/gh/v4lproik/simple-blockchain-quickstart/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/v4lproik/simple-blockchain-quickstart/tree/master)
This is merely a skeleton that helps you quickly set up a simplified version of a blockchain app written in Golang.
## Getting started
### Install Golang & Useful packages
1. Install goenv
1. Install make
2. Install golang => 1.18.3
3. Download dependencies
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
## Building  
```
make build
```
## Testing
```
make test
```