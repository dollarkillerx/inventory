build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o inventory -ldflags "-s -w" cmd/main.go
	upx inventory