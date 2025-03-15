# Define the binary name
BINARY_NAME := ed

# Define the build directory (optional, for cleaner builds)
BUILD_DIR := build

# Define the installation directory
INSTALL_DIR := $(HOME)/bin

# All targets
all: build install

# Build the binary
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME)

# Install the binary
install: build
	@install -m 0755 $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR) # corrected line

# Clean up build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Uninstall the binary
uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
