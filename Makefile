BINARIES_DIRECTORY = "bin"
PROJECT_NAME = $(shell basename "$(PWD)")

_mkdir:
	@echo "  > Create bin directory ${BINARIES_DIRECTORY} ..."
	@mkdir ${BINARIES_DIRECTORY}

clean:
	@echo "  > Remove bin directory ${BINARIES_DIRECTORY} ..."
	@rm -rf ${BINARIES_DIRECTORY}

build: clean _mkdir
	@echo "  > Build project ${PROJECT_NAME} ..."
	@go build -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}
	@cp -r templates ${BINARIES_DIRECTORY}/templates
	@echo "...done"

build-docker: clean _mkdir
	@echo "  > Build project for docker ${PROJECT_NAME} ..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}
	@cp -r templates ${BINARIES_DIRECTORY}/templates
	@echo "...done"