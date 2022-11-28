SOURCE_DIR        := $(shell pwd)
BUILD_DIR         := ${SOURCE_DIR}/build
BUILD_PACKAGE     := github.com/tmarcu/breeji-offloader
BINARY_NAME       := breeji-offloader
VERSION           := 0.0.1

LDFLAGS           := -ldflags " -X main.VersionBuild=$(VERSION) -w -extldflags \"-static\" "

GOLANGCI_VERSION := 1.50.1

.PHONY: build-cross all

build: build-cross

build-cross: linux windows darwin

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 ${BUILD_PACKAGE}

windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -a ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe ${BUILD_PACKAGE}

darwin:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -a ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64 ${BUILD_PACKAGE}

lint:
	${GOPATH}/bin/golangci-lint run ./...

test:
	go test -cover ./... -coverprofile=coverage.out

coverage: test
	go tool cover -html=coverage.out -o coverage.html

clean:
	go clean ./...
	rm -rf $(BUILD_DIR) *.upx upx coverage.out

release:
	git archive --format=tar.gz --verbose -o ${BINARY_NAME}-${VERSION}.tar.gz HEAD --prefix=${BINARY_NAME}-${VERSION}/
	#-git tag v${VERSION} # ignore error as it likely means the tag already exists

prepare-lint:
	curl -L -O https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_VERSION}/golangci-lint-${GOLANGCI_VERSION}-darwin-arm64.tar.gz && \
	tar -xf golangci-lint-${GOLANGCI_VERSION}-darwin-arm64.tar.gz && \
	cp golangci-lint-${GOLANGCI_VERSION}-darwin-arm64/golangci-lint ${GOPATH}/bin/ && \
	rm -r golangci-lint-${GOLANGCI_VERSION}-*
