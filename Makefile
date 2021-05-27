GOCMD ?= go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
BINARY_NAME = auto-allow-web

OUTPUT_PATH = build

TARGET = $(OUTPUT_PATH)/$(BINARY_NAME)
SOURCES = ./cmd/server

INSTALL_PATH = /usr/local/sbin

TAG ?= $(shell git describe --tags --abbrev=0 2>/dev/null)

default: all

all: test $(TARGET)

build:
	mkdir -p $(OUTPUT_PATH) || 0
	$(GOBUILD) -o $(TARGET) -v -ldflags="-X main.VERSION=$(TAG)" $(SOURCES)

$(TARGET): build

test:
#	$(GOTEST) -v

clean:
	$(GOCLEAN)
	rm -f $(TARGET)

run: $(TARGET)
	$(TARGET)

install: $(TARGET)
	install -g root -o root -p -s $(TARGET) $(INSTALL_PATH)

.PHONY: build clean run install all default test