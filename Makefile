PROG_NAME := "yum-get-repomd"
VERSION = 0.1.$(shell date +%Y%m%d.%H%M)
FLAGS := "-s -w -X main.version=${VERSION}"


build:
	#go mod tidy
	#go mod vendor
	CGO_ENABLED=0 go build -ldflags=${FLAGS} -o ${PROG_NAME} main.go repomd.go filelib.go loadKeys.go 
