ROOT=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
BIN_PATH=$(ROOT)/bin
GOCMD=go
GOBUILD=$(GOCMD) build
GOBUILD_PLUGIN=$(GOBUILD) -buildmode=plugin
GOCLEAN=$(GOCMD) clean
BINARY_NAME=$(BIN_PATH)/crankshaft

PLUGINS=$(notdir $(wildcard plugins/*))

server: plugins
	$(GOBUILD) -o $(BINARY_NAME) main.go
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

plugins: $(PLUGINS)

$(PLUGINS):
	cd plugins/$@ && $(GOBUILD_PLUGIN) -o $(BIN_PATH)/modules/$@.so