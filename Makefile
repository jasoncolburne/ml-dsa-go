tidy:
	go mod tidy

test: tidy
	go test ./...

test-neon: tidy
	go test -tags=neon ./...

benchmarks: tidy
	go test -bench=.

format:
	gofmt -w .
