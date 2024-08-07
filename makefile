default:
	go build main.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go

run:
	go run main.go player
