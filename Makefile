build:
	go build -o server ./main.go 

run: build
	./server

watch:
	reflex -s -r '\.go$$' make run

lint:
	golangci-lint run --fix --config .golangci.yml

fix:
	golangci-lint run --fix