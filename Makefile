all: build

build_linux:
	@env GOOS=linux GOARCH=amd64 go build -o ./bin/chat-server.linux.amd64 github.com/spring1843/chat-server/src/

build_osx:
	@env GOOS=darwin GOARCH=amd64 go build -o ./bin/chat-server.osx.amd64 github.com/spring1843/chat-server/src/

build_windows:
	@env GOOS=windows GOARCH=amd64 go build -o ./bin/chat-server.windows.amd64.exe github.com/spring1843/chat-server/src/

remove_old_package:
	@rm -f ./bin/latest_package.zip

add_static_files: remove_old_package
	@cd ./bin && zip -vr ./latest_package.zip  ./static-web && cd ..

add_to_zip_package: add_static_files
	@zip -vj ./bin/latest_package.zip Dockerfile ./bin/Dockerrun.aws.json ./bin/chat-server.linux.amd64 ./bin/config.json

build_all_targets: build_linux build_osx build_windows

build: build_all_targets add_to_zip_package
