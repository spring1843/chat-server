GOFLAGS ?= $(GOFLAGS:)

all: install test

build:
	@go build $(GOFLAGS) .

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

longtest:
	@LONGTESTS=1 go test -race ./longtests/...

mergecheck : beautify audit longtest

run: build
	@./src -config config.json

commit: beautify audit
	@git add -p .

virtuous:
	@for i in {1..100}; do make beautify audit; sleep 1;done

bench: install
	@go test -run=NONE -bench=. $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i ./...

certs:
	@openssl req -x509 -nodes -newkey rsa:2048 -keyout tls.key -out tls.crt -days 9650

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
	@go list -json | jq ".Deps" | grep -E -o "github.com+/[^\/]+/[^\/]+" | sed 's|[\",[:space:]]||g' | sort -u

fix:
	@go fix $(GOFLAGS) ./...
