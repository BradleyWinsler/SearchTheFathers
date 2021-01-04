TAG=$(shell git describe --always)
DOCKER_REGISTRY=localhost:5000/search-the-fathers

.PHONY: sure
sure:
	go mod tidy

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -v -o ./app -ldflags="-s -w" ./main.go
	docker build -t $(DOCKER_REGISTRY):$(TAG) .
	docker-compose up --build --no-deps
