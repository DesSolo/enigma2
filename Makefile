PROJECT_NAME = $(shell basename "$(PWD)")
BINARIES_DIRECTORY = bin
VERSION = $(shell cat VERSION)
LDFLAGS = "-w -s"

.DEFAULT_GOAL := build-server

_static:
	mkdir ${BINARIES_DIRECTORY}
	cp -r templates ${BINARIES_DIRECTORY}/templates

clean:
	rm -rf ${BINARIES_DIRECTORY}

build-server: clean _static
	go build -ldflags ${LDFLAGS} -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_server_linux_x64 cmd/server/main.go

build-client: clean
	GOOS=linux GOARCH=amd64 go build -ldflags ${LDFLAGS} -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_client_linux_x64 cmd/client/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags ${LDFLAGS} -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_client_windows_x64.exe cmd/client/main.go

build-all: build-server build-client

build-docker: clean _static
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_docker cmd/server/main.go

tag:
ifeq (${VERSION},$(shell git describe --abbrev=0))
	$(error Last tag and curent version ${VERSION} conflict)
else
	git tag -a ${VERSION} -m 'Bump Version ${VERSION}'
	git push --tags
endif