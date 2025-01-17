#!/bin/bash

swaggerCmd="docker run --rm -it --user $(id -u):$(id -g) -e GOPATH=$HOME/go:/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger:v0.28.0"
clientDir=dstore
apiRepo="ssh://git-codecommit.eu-central-1.amazonaws.com/v1/repos/swagger-specs"

rm -rf $clientDir/client
rm -rf $clientDir/models
rm -rf $clientDir/swagger-specs
git clone $apiRepo $clientDir/swagger-specs
$swaggerCmd generate client -f $clientDir/swagger-specs/d-store/api.yaml -t $clientDir
rm -rf $clientDir/swagger-specs
