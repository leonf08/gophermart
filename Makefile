build:
	go build -o ./bin/ ./cmd/...

serve-swagger:
	swagger serve -F=swagger ./docs/swagger.yaml

clean:
	rm -rf ./bin