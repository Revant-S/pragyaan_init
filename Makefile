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

gitHooks:
	bash scripts/setup_hooks.sh
initial_setup:
	go mod tidy
	make gitHooks
	make build
	swag init
	make watch