build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/migrate handlers/migrate/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/scheduler handlers/scheduler/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/producer handlers/producer/main.go
