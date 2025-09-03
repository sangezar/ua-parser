BINARY := ua-parser
BIN_DIR := bin
PKG := .

LDFLAGS ?= -s -w

.PHONY: all linux-amd64 darwin-arm64 clean

all: linux-amd64 darwin-arm64

linux-amd64:
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY)-linux-amd64 $(PKG)

darwin-arm64:
	mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY)-darwin-arm64 $(PKG)

clean:
	rm -f $(BIN_DIR)/$(BINARY)-linux-amd64 $(BIN_DIR)/$(BINARY)-darwin-arm64


