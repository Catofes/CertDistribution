all:
		mkdir -p build
		env GO111MODULE=on CGO_ENABLED=0 go build -o build/web -ldflags "-X 'main._version_=$(shell git log --pretty=format:"%h" -1)'" .


