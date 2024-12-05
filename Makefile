build:
	go build -o server ./main.go 

run: build
	./server

lint:
	golangci-lint run --fix --config .golangci.yml

fix:
	golangci-lint run --fix
