.PHONY: all build build_frontend build_backend run clean

all: build

build: build_frontend build_backend

build_frontend:
	@echo "Building frontend..."
	@(cd web && npm run build)

build_backend:
	@echo "Building backend..."
	@go build -ldflags="-s -w" -o gandalf-budget ./cmd/server

run: build
	@echo "Running gandalf-budget..."
	@./gandalf-budget

clean:
	@echo "Cleaning up..."
	@rm -f gandalf-budget
	@(cd web && rm -rf dist)
	@echo "Done."

# Target to remind user about npm install, not run automatically
install_frontend_deps:
	@echo "Make sure to run 'npm install' in the 'web/' directory if you haven't already."
