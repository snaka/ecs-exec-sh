.PHONY: build tidy fmt install

build:
	go build -o ecs-exec-sh ./cmd/ecs-exec-sh

tidy:
	go mod tidy

fmt:
	go fmt ./...

install:
	go install github.com/snaka/ecs-exec-sh/cmd/ecs-exec-sh
