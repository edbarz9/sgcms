build:
	go build main.go

pi:
	GOARCH=arm GOOS=linux go build main.go

run:
	go run main.go
