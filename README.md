# Course code for [Building a Secrets Sharing Application](https://www.manning.com/liveproject/build-a-secrets-sharing-web-application) from Manning press.

$ go test -v ./pkg/handlers/...
$ go build -o secrets-app cmd/secrets-server/main.go
$ PASSWORD=one SALT=two DATA_FILE_PATH=/tmp/data.json go run .
