PROJECT_NAME := "goexamples"
PKG := "github.com/william/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY:

lint:
  go fmt ${PKG_LIST}
  go vet ${PKG_LIST}
  golint -set_exit_status ${PKG_LIST}

test:
  @go test -msan -race -short -v ${PKG_LIST}

tidy:
  @go mod tidy

build:
  @go build -o goexamples ${PKG_LIST}

clean:
  @rm -f $(PROJECT_NAME)

image:
  @docker build -t goexample .