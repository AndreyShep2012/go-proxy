build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-proxy ./main.go

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o go-proxy ./main.go

