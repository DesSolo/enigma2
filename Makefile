PROJECT_NAME = $(shell basename "$(PWD)")
BINARIES_DIRECTORY = bin
VERSION = $(shell cat VERSION)
LDFLAGS = "-w -s"

.DEFAULT_GOAL := help

_static:
	mkdir ${BINARIES_DIRECTORY}
	cp -r templates ${BINARIES_DIRECTORY}/templates

## clean: Clean binaries directory
clean:
	rm -rf ${BINARIES_DIRECTORY}

## test: Run all tests
test:
	go test ./...

## vet: Run go vet
vet:
	go vet ./...

## build-server: Build enigma server
build-server: clean _static
	go build -ldflags ${LDFLAGS} -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_server_linux_x64 cmd/server/main.go

## build-server-tar: Build server and compress to tar.gz
build-server-tar: build-server
	for filename in ${BINARIES_DIRECTORY}/enigma_server* ; do \
		echo "  > start compress $$filename ..." ; \
		tar -zcvf $$filename.tar.gz $$filename ${BINARIES_DIRECTORY}/templates ; \
		rm $$filename ;\
		echo "  > ... done" ; \
	done

## build-client: Build enigma client
build-client: clean
	GOOS=linux GOARCH=amd64 go build -ldflags ${LDFLAGS} -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_client_linux_x64 cmd/client/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags ${LDFLAGS} -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_client_windows_x64.exe cmd/client/main.go

## build-all: Build client and server
build-all: build-server build-client

## build-docker: Building a binary file to run in a container
build-docker: clean _static
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ${BINARIES_DIRECTORY}/${PROJECT_NAME}_docker cmd/server/main.go

## tag: Create and push git tag
tag:
ifeq (${VERSION},$(shell git describe --abbrev=0))
	$(error Last tag and curent version ${VERSION} conflict)
else
	git tag -a ${VERSION} -m 'Bump Version ${VERSION}'
	git push --tags
endif

## upload-github: Add binary to github releases
upload-github: build-server-tar build-client
	@bash scripts/upload-github-release-asset.sh github_api_token=${GITHUB_API_TOKEN} owner=DesSolo repo=enigma2 tag=${VERSION} files=./${BINARIES_DIRECTORY}/*

## help: Show this message and exit
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo