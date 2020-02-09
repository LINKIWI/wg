APP_NAME = wg
BIN_DIR = bin

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

all: $(APP_NAME)

$(APP_NAME):
	go build \
	    -ldflags "-w -s -X wg/internal/meta.Version=$(VERSION)" \
	    -o $(BIN_DIR)/$(APP_NAME)-$(GOOS)-$(GOARCH) \
	    cmd/wg/main.go

lint:
	.ci/lint.sh

clean:
	rm -f $(BIN_DIR)/*

.PHONY: lint clean
