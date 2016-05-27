SOURCES := $(filter-out $(wildcard *_test.go),$(wildcard *.go))
TARGET := cwimport

GO_VERSION := 1.6.2 # for docker image only

$(TARGET): $(SOURCES)
	docker run \
		--rm \
		-v $(PWD):/go/src/github.com/trayio/$(TARGET) \
		-w /go/src/github.com/trayio/$(TARGET) \
		-e CGO_ENABLED=0 \
		golang:$(GO_VERSION) go build --ldflags '-extldflags "-static"' -o $(TARGET)

clean:
	rm -f $(TARGET)
