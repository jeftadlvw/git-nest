.PHONY: build

UNAME := $(shell uname)
ifeq ($(OS), WINDOWS_NT)
	RM := del /f
	TIME := Get-Date -UFormat %s
	DETECTED_OS := Windows
else
	RM := rm
	TIME := date +%s
	DETECTED_OS := $(UNAME)
endif

# define general build variables
APP_NAME := git-nest
BUILD_DIR := _build
ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
CURRENT_TIMESTAMP := $(shell $(TIME))

# define values to inject at compile time
INJECT_BASE := "github.com/jeftadlvw/git-nest/internal"
INJECT_VERSION_KEY := $(INJECT_BASE)/constants.version
INJECT_VERSION_VALUE := ${GIT_NEST_BUILD_VERSION}

INJECT_COMMIT_KEY := $(INJECT_BASE)/constants.ref
INJECT_COMMIT_VALUE := ${GIT_NEST_BUILD_COMMIT_SHA}

INJECT_COMPILE_TIME_KEY := $(INJECT_BASE)/constants.compilationTimestampStr
INJECT_COMPILE_TIME_VALUE := $(or ${GIT_NEST_COMPILE_TIME}, $(CURRENT_TIMESTAMP))

INJECT_EPHEMERAL_BUILD_KEY := $(INJECT_BASE)/constants.ephemeralBuildStr
INJECT_EPHEMERAL_BUILD_VALUE := false

# build for local OS and architecture
build: clean test
	@go build \
		-o $(BUILD_DIR)/$(APP_NAME) \
		-ldflags " \
			-X $(INJECT_VERSION_KEY)=$(INJECT_VERSION_VALUE) \
			-X $(INJECT_COMMIT_KEY)=$(INJECT_COMMIT_VALUE) \
			-X $(INJECT_COMPILE_TIME_KEY)=$(INJECT_COMPILE_TIME_VALUE) \
			-X $(INJECT_EPHEMERAL_BUILD_KEY)=$(INJECT_EPHEMERAL_BUILD_VALUE) \
			"

clean:
	@$(RM) -rf _build

test:
	@go test ./.../tests

git-test: build
	@PATH="$$PATH:$(ROOT_DIR)/build" git nest

test-env:
	@docker build -t git-nest/test-env _docker/test_env
	@CONTAINER_NAME=git-nest_testenv; \
		if [[ "$$(docker ps -aqf name=$$CONTAINER_NAME)" ]]; then \
			docker start -ai $$CONTAINER_NAME; \
		else \
			docker run -it --name $$CONTAINER_NAME -v ./:/app/src:ro git-nest/test-env; \
		fi

debug:
	@echo "OS:\t\t\t$(DETECTED_OS)"
	@echo "ROOT_DIR:\t\t$(ROOT_DIR)"
	@echo "BUILD_DIR:\t\t$(BUILD_DIR)"
	@echo
	@echo "INJECT_VERSION:\t\t$(INJECT_VERSION_KEY) -> $(INJECT_VERSION_VALUE)"
	@echo "INJECT_COMMIT:\t\t$(INJECT_COMMIT_KEY) -> $(INJECT_COMMIT_VALUE)"
	@echo "INJECT_TIME:\t\t$(INJECT_COMPILE_TIME_KEY) -> $(INJECT_COMPILE_TIME_VALUE)"
	@echo "INJECT_EPHEMERAL:\t$(INJECT_EPHEMERAL_BUILD_KEY) -> $(INJECT_EPHEMERAL_BUILD_VALUE)"
