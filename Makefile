BRANCH        := $(shell git rev-parse --abbrev-ref HEAD)
SHORT_COMMIT  := $(shell git rev-parse --short HEAD)
VERSION       := $(BRANCH) $(SHORT_COMMIT)

OUTPUT_DIR := squirrel
TAR_NAME := squirrel-$(BRANCH).tar.gz

GOOS ?= linux
GOARCH ?= amd64

CMDS := squ-apiserver squ-agent squctl
BINS := $(addprefix $(OUTPUT_DIR)/, $(CMDS))
CONFIGS := $(addprefix $(OUTPUT_DIR)/config/, agent.yaml apiserver.yaml squctl.yaml)

.PHONY: clean all package $(CMDS) image frontend

all: frontend $(BINS) $(CONFIGS)

frontend:
	@echo "Building frontend..."
	@cd front && npm install&& npm run build
	@rm -rf internal/squ-apiserver/server/dist
	@mv front/dist internal/squ-apiserver/server
	@echo "Frontend built and moved to internal/squ-apiserver/server/dist"

package: all
	tar -czf $(TAR_NAME) -C $(OUTPUT_DIR) .
	mv $(TAR_NAME) $(OUTPUT_DIR)

$(CMDS): %: $(OUTPUT_DIR)/%

$(OUTPUT_DIR)/%: cmd/%/*.go
	@mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 \
	go build \
	-ldflags '-X "squirrel-dev/cmd/$*/app.Version=$(VERSION)"' \
	-o $@ ./cmd/$*

$(OUTPUT_DIR)/config/%.yaml: config/%.yaml
	@mkdir -p $(OUTPUT_DIR)/config
	cp -f $< $@

clean:
	rm -rf $(OUTPUT_DIR)

image: all
	docker build -f dockerfiles/Dockerfile-apiserver -t gosquirrel/squ-apiserver .
	docker build -f dockerfiles/Dockerfile-agent -t gosquirrel/squ-agent .