PROJECT_NAME=i3br
BIN_DIR=bin

run:
	@echo "Running the project..."
	@go run .

build:
	@echo "Building the project..."
	@go build -o $(BIN_DIR)/$(PROJECT_NAME) . 
	@echo "Binary has been placed in $(BIN_DIR)/$(PROJECT_NAME)"

install: build
	@echo "Adding $(BIN_DIR) to PATH..."
	@echo 'export PATH=$$PATH:$(shell pwd)/$(BIN_DIR)' >> ~/.bashrc
	@echo "$(BIN_DIR) has been added to PATH"

clean:
	@echo "Cleaning the binary directory..."
	@rm -rf $(BIN_DIR)

.PHONY: run build install clean
