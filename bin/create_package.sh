#!/usr/bin/env bash

rm latest_package.zip
zip -vj ./latest_package.zip Dockerfile
zip -vj ./latest_package.zip Dockerrun.aws.json
zip -vj ./latest_package.zip chat-server.linux.amd64
zip -vj ./latest_package.zip config.json

# Go back one dir so the static
zip -rv  ./latest_package.zip ./static-web
