SOURCES := $(filter-out $(wildcard *_test.go),$(wildcard *.go))
TARGET := $(shell basename `pwd`)

$(TARGET): $(SOURCES)
	CGO_ENABLED=0 go build --ldflags '-extldflags "-static"'

clean:
	rm -f $(TARGET)
