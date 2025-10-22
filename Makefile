APP_NAME := jebi

.PHONY: all test install

all: test

test:
	@echo "🧪 Running tests..."
	go test ./... -v

install:
	@echo "📦 Installing $(APP_NAME)..."
	go install ./...
