BRANCH        := $(shell git rev-parse --abbrev-ref HEAD)
SHORT_COMMIT  := $(shell git rev-parse --short HEAD)
VERSION       := $(BRANCH) $(SHORT_COMMIT)

OUTPUT_DIR := output

GOOS ?= windows
GOARCH ?= amd64

CMDS := squ-apiserver squ-agent squctl
BINS := $(addprefix $(OUTPUT_DIR)/, $(CMDS))

.PHONY: clean all $(CMDS)

all: $(BINS)

$(CMDS): %: $(OUTPUT_DIR)/%

$(OUTPUT_DIR)/%: cmd/%/*.go
	@mkdir -p $(OUTPUT_DIR)
	go build \
	-ldflags '-X "squirrel-dev/cmd/$*/app.Version=$(VERSION)"' \
	-o $@ ./cmd/$*

clean:
	rm -rf $(OUTPUT_DIR)