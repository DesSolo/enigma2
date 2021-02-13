BINARIES_DIRECTORY = bin
PROJECT_NAME = $(shell basename "$(PWD)")

.DEFAULT_GOAL := build-server

_static:
	mkdir ${BINARIES_DIRECTORY}
	cp -r templates ${BINARIES_DIRECTORY}/templates

clean:
	rm -rf ${BINARIES_DIRECTORY}

build-server: clean _static
	go build -o ${BINARIES_DIRECTORY}/${PROJECT_NAME} cmd/server/main.go

build-client: clean
	GOOS=linux GOARCH=amd64 go build -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_cl_linux_x64 cmd/client/main.go
	GOOS=windows GOARCH=amd64 go build -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_cl_windows_x64.exe cmd/client/main.go

build-docker: clean _static
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_docker cmd/server/main.go
