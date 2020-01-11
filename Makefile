NAME=staple

# Set the build dir, where built cross-compiled binaries will be output
BUILDDIR := build

# List the GOOS and GOARCH to build
GO_LDFLAGS_STATIC=-ldflags "-s -w $(CTIMEVAR) -extldflags -static"

.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -ldflags="-s -w" -i -o ${BUILDDIR}/${NAME} cmd/root.go

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	go clean -i

.PHONY: run
run:
	go run cmd/root.go

.PHONY: start-https
start-https:
	go run cmd/root.go --server-key-path ./certs/key.pem --server-crt-path ./certs/cert.pem
