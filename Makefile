SOURCES := $(filter-out $(wildcard *_test.go),$(wildcard *.go))
TARGET := $(shell basename `pwd`)

GO := $(shell command -v go)
GO_VERSION := 1.6.0 # for docker image only

ifndef GO
	GO := docker run --rm -v $(PWD):/go/src/github.com/trayio/$(TARGET) -w /go/src/github.com/trayio/$(TARGET) -e CGO_ENABLED=0 golang:$(GO_VERSION) go
endif

$(TARGET): $(SOURCES)
	CGO_ENABLED=0 $(GO) build --ldflags '-extldflags "-static"'

clean:
	rm -f $(TARGET)
