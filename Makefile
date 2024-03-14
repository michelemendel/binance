build_client:
	@go build -o bin/client cmd/client/main.go

client_prod: build_client
	@ENV=prod ./bin/client

client_test: build_client
	@ENV=test ./bin/client

test_all:
	go test -v ./... -count=1

clean:
	rm -rf bin


	