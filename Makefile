.PHONY: all build build_frontend build_backend run clean

all: build

build: build_frontend build_backend

build_frontend:
	@echo "Building frontend..."
	@(cd web && npm run build)

build_backend:
	@echo "Building backend..."
	@echo "Preparing embedded web assets..."
	@rm -rf cmd/server/embedded_web_dist 
	@mkdir -p cmd/server/embedded_web_dist
	@cp -r web/dist/* cmd/server/embedded_web_dist/
	@go build -ldflags="-s -w" -o gandalf-budget ./cmd/server

run: build
	@echo "Running gandalf-budget..."
	@./gandalf-budget

clean:
	@echo "Cleaning up..."
	@rm -f gandalf-budget
	@rm -rf cmd/server/embedded_web_dist
	@(cd web && rm -rf dist)
	@echo "Done."

# Target to remind user about npm install, not run automatically
install_frontend_deps:
	@echo "Make sure to run 'npm install' in the 'web/' directory if you haven't already."
