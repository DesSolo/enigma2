BINARIES_DIRECTORY = "bin"
PROJECT_NAME = $(shell basename "$(PWD)")

_static:
	@echo "  > Create bin directory ${BINARIES_DIRECTORY} ..."
	@mkdir ${BINARIES_DIRECTORY}
	@echo "  > Copy static files ..."
	@cp -r templates ${BINARIES_DIRECTORY}/templates

clean:
	@echo "  > Remove bin directory ${BINARIES_DIRECTORY} ..."
	@rm -rf ${BINARIES_DIRECTORY}

build-server: clean _static
	@echo "  > Build server ${PROJECT_NAME} ..."
	@go build -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}
	@echo "...done"

build-client: clean
	@echo "  > Build clients ${PROJECT_NAME} ..."
	@echo "  > linux x64 ..."
	@GOOS=linux GOARCH=amd64 go build -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_cl_linux_x64 client/client.go
	@echo "  > windows x64 ..."
	@GOOS=windows GOARCH=amd64 go build -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_cl_windows_x64.exe client/client.go
	@echo "...done"

build-docker: clean _static
	@echo "  > Build project for docker ${PROJECT_NAME} ..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_docker
	@echo "...done"