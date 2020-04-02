
.PHONY: build
build:
	go build -o build/main ./cmd/goblast/main.go

.PHONY: test
test:
	go test ./test/...
