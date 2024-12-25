build:
	@ go build -o bin/gobank

run: build
	@ ./bin/gobank

migrate: 
	@  go run ./migrations/migrator/main.go

test:
	@ go test -v	