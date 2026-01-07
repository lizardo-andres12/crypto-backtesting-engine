# Wrapper Makefile for CMake
# This allows you to run "make run" or "make test" from the root
# without worrying about cmake commands.

BUILD_DIR := build
EXECUTABLE := $(BUILD_DIR)/execution_engine.out

TYPE ?= Release

.PHONY: all build run clean

all: build

build:
	@mkdir -p $(BUILD_DIR)
	@cmake -S . -B $(BUILD_DIR) \
                -DCMAKE_BUILD_TYPE=$(TYPE)
	@cmake --build $(BUILD_DIR)
	@cp $(BUILD_DIR)/compile_commands.json .  # So clangd LSP stops complaining about `#include` paths

run: $(EXECUTABLE)
	./$(EXECUTABLE)

clean:
	@rm -rf $(BUILD_DIR)/*

