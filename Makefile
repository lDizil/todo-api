run:
	docker-compose up

docker-rebuild:
	docker-compose build up

.PHONY: test-db-up test-db-down test-integration

test-db-up:
	docker-compose -f docker-compose.test.yml up -d

test-db-down:
	docker-compose -f docker-compose.test.yml down -v

test:
	go test ./... -cover
	go test ./...
	go test -v ./internal/tests/integration/...
