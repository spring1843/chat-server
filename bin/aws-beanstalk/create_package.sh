#!/usr/bin/env bash

rm latest_package.zip
zip -vj ./latest_package.zip ../../Dockerfile
zip -vj ./latest_package.zip ../../Dockerrun.aws.json
zip -vj ./latest_package.zip ../chat-server.linux.amd64
zip -vj ./latest_package.zip ../config.json
zip -vj ./latest_package.zip ../../src/tls.crt
zip -vj ./latest_package.zip ../../src/tls.key
cd ../../ && zip -rv  ./bin/aws-beanstalk/latest_package.zip ./static-web && cd ./bin/aws-beanstalk
