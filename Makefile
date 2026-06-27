# Binary output name
BIN ?= minesweeper-cli

# Package name
PKG := github.com/nicolito128/minesweeper-go

# Architecture
ARCH ?= $(shell go env GOOS)-$(shell go env GOARCH)

# Program version
VERSION ?= main

# Output directory
OUTPUT_DIR ?= bin

# Go environment
platform = $(subst -, ,$(ARCH))
GOOS = $(word 1, $(platform))
GOARCH = $(word 2, $(platform))
GOPROXY ?= "https://proxy.golang.org,direct"

ifeq ($(DEBUG),1)
    GCFLAGS := all=-N -l
else
    GCFLAGS := 
endif

LDFLAGS := -s -w

all:
	@$(MAKE) build

build: $(OUTPUT_DIR)/$(GOOS)/$(GOARCH)/$(BIN)

$(OUTPUT_DIR)/$(GOOS)/$(GOARCH)/$(BIN): make-dirs
	@echo "building: $@"
	GOPROXY=$(GOPROXY) GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-o "$(OUTPUT_DIR)/$(GOOS)/$(GOARCH)/$(BIN)" \
		-gcflags="$(GCFLAGS)" \
		-installsuffix "static" \
		-ldflags="$(LDFLAGS)" \
		"$(PKG)/cmd/$(BIN)"

make-dirs:
	@mkdir -p $(OUTPUT_DIR)/$(GOOS)/$(GOARCH)

clean:
	find $(OUTPUT_DIR)/ -mindepth 1 ! -name "*.md" -exec rm -rf {} +

tidy:
	go mod tidy

run: build
	./$(OUTPUT_DIR)/$(GOOS)/$(GOARCH)/$(BIN) $(ARGS) \

.PHONY: all make-dirs clean tidy run
