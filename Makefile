default: build

build:
	@go build

test:
	@go test

cover:
	@go test -tags test -coverprofile=c
	@go tool cover -html=c

clean:
	@rm -f up c

sure:
	@go test -race -cpu=4
	@go fmt
	@go vet
	@golint ./...
	@go build -gcflags=-s
	@make -s clean

notes:
	@egrep "(NOTE|TODO|FIXME)" . -R|grep -v Makefile|grep -v .git

.PHONY: test build
