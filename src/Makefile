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

get_make_deps:
	@go get -u golang.org/x/tools/cmd/goimports
	@go get -u github.com/golang/lint/golint

test: install
	@go test $(GOFLAGS) ./...

long_test:
	@LONGTESTS=1 go test -race ./longtests/...

no_net_race:
	@SKIP_NETWORK=1 go test -race ./...

travis : beautify vet no_net_race lint

run: build
	@./src -config config.json

commit: beautify audit
	@git add -p .

race_loop:
	@for i in {1..100}; do make beautify audit; sleep 1;done

bench: install
	@go test -run=NONE -bench=. $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i ./...

certs:
	@openssl req -x509 -nodes -newkey rsa:2048 -keyout tls.key -out tls.crt -days 9650

serve: build install
	@chat-server -c config.json

build_linux:
	@env GOOS=linux GOARCH=amd64 go build -o ../bin/chat-server.linux.amd64 github.com/spring1843/chat-server/src/

build_osx:
	@env GOOS=darwin GOARCH=amd64 go build -o ../bin/chat-server.osx.amd64 github.com/spring1843/chat-server/src/

build_windows:
	@env GOOS=windows GOARCH=amd64 go build -o ../bin/chat-server.linux.amd64.exe github.com/spring1843/chat-server/src/

build_bin: build_linux build_osx build_windows

deps:
	@go list -json | jq ".Deps" | grep -E -o "github.com+/[^\/]+/[^\/]+" | sed 's|[\",[:space:]]||g' | sort -u

fix:
	@go fix $(GOFLAGS) ./...
