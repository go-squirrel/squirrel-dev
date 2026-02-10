BRANCH        := $(shell git rev-parse --abbrev-ref HEAD)
SHORT_COMMIT  := $(shell git rev-parse --short HEAD)
VERSION       := $(BRANCH) $(SHORT_COMMIT)

OUTPUT_DIR := squirrel
TAR_NAME := squirrel-$(BRANCH).tar.gz

GOOS ?= linux
GOARCH ?= amd64

# 支持的架构列表
SUPPORTED_ARCHS := amd64 arm64

# 多架构构建使用的输出目录
MULTIARCH_OUTPUT_DIR := $(OUTPUT_DIR)/multiarch

CMDS := squ-apiserver squ-agent squctl
BINS := $(addprefix $(OUTPUT_DIR)/, $(CMDS))
CONFIGS := $(addprefix $(OUTPUT_DIR)/config/, agent.yaml apiserver.yaml squctl.yaml)

.PHONY: clean all package $(CMDS) image frontend all-arch package-all-arch

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

# 多架构构建
all-arch:
	@for arch in $(SUPPORTED_ARCHS); do \
		echo "Building for linux/$$arch..."; \
		$(MAKE) GOOS=linux GOARCH=$$arch OUTPUT_DIR=$(MULTIARCH_OUTPUT_DIR)/linux-$$arch build-bins; \
		$(MAKE) OUTPUT_DIR=$(MULTIARCH_OUTPUT_DIR)/linux-$$arch build-configs; \
	done

# 单独构建二进制文件（不依赖前端）
build-bins: $(BINS)

# 单独构建配置文件
build-configs: $(CONFIGS)

# 构建所有架构的压缩包
package-all-arch: frontend all-arch
	@for arch in $(SUPPORTED_ARCHS); do \
		echo "Packaging for linux-$$arch..."; \
		tar -czf $(OUTPUT_DIR)/squirrel-$(BRANCH)-linux-$$arch.tar.gz -C $(MULTIARCH_OUTPUT_DIR)/linux-$$arch .; \
	done
	echo "Done.";

image: all
	docker build -f dockerfiles/Dockerfile-apiserver -t gosquirrel/squ-apiserver .
	docker build -f dockerfiles/Dockerfile-agent -t gosquirrel/squ-agent .