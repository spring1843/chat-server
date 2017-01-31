GOFLAGS ?= $(GOFLAGS:)

all: install test

build:
	@go build $(GOFLAGS) ./...

format:
	@gofmt -l -s -w .

lint:
	@golint ./...

optimize_imports:
	@goimports -l -w .

beautify: format optimize_imports

vet:
	@go vet ./...

race:
	@go test -race $(GOFLAGS) ./...

audit: vet race lint

install:
	@go get $(GOFLAGS) ./...

test: install
	@go test $(GOFLAGS) ./...

commit : beautify audit
	@git add -p .

bench: install
	@go test -run=NONE -bench=. $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i ./...

serve: build install
	@chat-server -c config.json

buildlinux:
	@env GOOS=linux GOARCH=amd64 go build -o ../bin/chat-server.linux.amd64 github.com/spring1843/chat-server/src/

buildosx:
	@env GOOS=darwin GOARCH=amd64 go build -o ../bin/chat-server.osx.amd64 github.com/spring1843/chat-server/src/

buildwindows:
	@env GOOS=windows GOARCH=amd64 go build -o ../bin/chat-server.linux.amd64.exe github.com/spring1843/chat-server/src/

buildbin: buildlinux buildosx buildwindows

deps:
	@go list -json | jq ".Deps" | grep -E -o "github.com+/[^\/]+" | sed 's|[\",[:space:]]||g' | sort -u

fix:
	@go fix $(GOFLAGS) ./...
