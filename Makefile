.PHONY: build test run clean

build:
	go build -o worker ./cmd/worker

test:
	go test ./...

vet:
	go vet ./...

run:
	go run ./cmd/worker

clean:
	rm -f worker worker.exe

tidy:
	go mod tidy

help:
	@echo "Available targets:"
	@echo "  build  - Build the worker binary"
	@echo "  test   - Run tests"
	@echo "  vet    - Run go vet"
	@echo "  run    - Run worker locally"
	@echo "  clean  - Remove built binaries"
	@echo "  tidy   - Clean up go.mod"
