.PHONY: compose-up compose-down test cover mockgen swag

compose-up:
	docker-compose up --build -d && docker-compose logs -f

compose-down:
	docker-compose down --remove-orphans

test:
	go test -v ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out

mockgen:
	mockgen -source=internal/service/service.go -destination=internal/mocks/servicemocks/service.go -package=servicemocks
	mockgen -source=internal/repo/repo.go       -destination=internal/mocks/repomocks/repo.go       -package=repomocks

swag:
	swag init -g cmd/app/main.go --parseInternal --parseDependency
