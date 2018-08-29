build:
	go get github.com/aws/aws-sdk-go/aws/...
	go get github.com/satori/go.uuid
	go get github.com/Squwid/bytegolf/bgaws
	go get github.com/Squwid/bytegolf/runner

	env GOOS=linux go build -ldflags="-s -w" -o bin/bytegolf *.go