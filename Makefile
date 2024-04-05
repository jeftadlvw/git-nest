ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
.PHONY: build

build:
	@rm -rf build
	@go build -o build/

git-test: build
	@PATH="$$PATH:$(ROOT_DIR)/build" git nest

root-dir:
	@echo $(ROOT_DIR)
