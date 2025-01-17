#!/bin/bash

swaggerCmd="docker run --rm -it --user $(id -u):$(id -g) -e GOPATH=$HOME/go:/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger:v0.27.0"
apiRepo="ssh://git-codecommit.eu-central-1.amazonaws.com/v1/repos/swagger-specs"
apiDir=api

echo "Cleaning existing API"
rm -rf $apiDir/models/*

echo "Cloning latest swagger specs"
git clone $apiRepo $apiDir/swagger-specs

echo "Generating server models"
$swaggerCmd generate model -t $apiDir -f $apiDir/swagger-specs/slots/api.yaml
rm -rf $apiDir/swagger-specs

scripts/strip-api-models.sh
rm $apiDir/models/messages_response.go
