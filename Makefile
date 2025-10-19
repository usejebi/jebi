run:
	@echo "Running the application..."
	go run main.go

build:
	@echo "Building the application..."
	go build -o gfs main.go

test:
	@echo "Running tests..."
	go test ./...

install:
	@echo "Installing the application..."
	go install main.go