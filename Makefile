run:
	go run cmd/todo/main.go TODO

build:
	go build cmd/todo/main.go -o build/todo

build-win:
	go build cmd/todo/main.go -o build/todo.exe