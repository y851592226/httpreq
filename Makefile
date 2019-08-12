test:
	go test -v
covertest:
	gotestcover -coverprofile=cover.out ./
	go tool cover -func=cover.out
	go tool cover -html=cover.out -o=cover.html
format:
	goimports -w ./*.go
	goimports -w ./binding/*.go
golangci-lint:
	golangci-lint run
