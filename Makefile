build:
	go build -v ./main.go
random:
	go run ./main.go random
dump:
	go run ./main.go dump -n 5

DEFAULT_GOAL := build