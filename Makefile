build:
	go build ./cmd/pets

run:
	go run ./cmd/pets

mock-all: mock-repository mock-service
mock-service:
	mockgen -source=internal/service/service.go -destination=mocks/service/mockService.go

mock-repository:
	mockgen -source=internal/repository/repository.go -destination=mocks/repository/mockRepository.go

test: mock-all
	go test ./...

cover: mock-all
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out