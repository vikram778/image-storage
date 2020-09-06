#!/usr/bin/env bash
# exit on error

set -eu

cd ${APP_DIR}/docs

# remove old swagger files if exists inside the definitions
if [ -f ./definitions/swagger.json ]; then
    echo ">>>>>>>> Remove the swagger json in definition..."
    rm ./definitions/swagger.json
fi

echo ">>>>>>>> Building the api definition..."
go run swagjson.go --input=./definitions -output=./definitions

# remove old swagger files if exists inside docs root
if [ -f ./swagger.json ]; then
    echo ">>>>>>>> Remove the swagger json in docs root..."
    rm ./swagger.json
fi

echo ">>>>>>>> Copying the api definition to docs root..."
#  on successful generation, move the swagger.json to docs root
mv ./definitions/swagger.json swagger.json

# remove existing docs root
rm -Rf ${SWAGGER_PATH}

# creates docs root
mkdir -p ${SWAGGER_PATH}

#move the docs to home docs
cp -a ./* ${SWAGGER_PATH}






