.PHONY: build run clean setup test

# Detect OS and set variables accordingly
ifeq ($(OS),Windows_NT)
    EXECUTABLE := tutorgpt.exe
    RM := del /Q
    PATH_SEP := \\
    COPY := copy
else
    EXECUTABLE := tutorgpt
    RM := rm -f
    PATH_SEP := /
    COPY := cp
endif

# Build the binary and move it to root
build:
	@echo "Building TutorGPT..."
	@cd backend && go build -o tutorgpt cmd/main.go
	@cd backend && set GOOS=windows&& go build -o tutorgpt.exe cmd/main.go
	@echo "Moving executables to root directory..."
ifeq ($(OS),Windows_NT)
	@if exist backend\tutorgpt.exe $(COPY) backend\tutorgpt.exe tutorgpt.exe
	@if exist backend\tutorgpt $(COPY) backend\tutorgpt tutorgpt
else
	@if [ -f backend/tutorgpt ]; then $(COPY) backend/tutorgpt ./tutorgpt; fi
	@if [ -f backend/tutorgpt.exe ]; then $(COPY) backend/tutorgpt.exe ./tutorgpt.exe; fi
endif

# Run the application
run: build
	@echo "Starting TutorGPT..."
	@$(EXECUTABLE)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@cd backend && $(RM) tutorgpt tutorgpt.exe
	@$(RM) tutorgpt tutorgpt.exe

# Setup development environment
setup:
	@echo "Setting up development environment..."
	@cd backend && go mod init tutorgpt && go mod tidy
ifeq ($(OS),Windows_NT)
	@if not exist .env (copy .env.example .env && echo Created .env file. Please update it with your actual values.)
else
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file. Please update it with your actual values."; \
	fi
endif

# Run tests
test:
	@echo "Running tests..."
	@cd backend && go test -v ./... 