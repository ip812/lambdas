run-local:
	@docker build -f Dockerfile.local -t hello:test .
	@docker run -p 9000:8080 hello:test

invoke-local:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" \
	-H "Content-Type: application/json" \
	-d '{ \
		"key": "foo", \
		"value": "bar" \
	}'

fmt:
	@go fmt $(shell go list ./... | grep -v /vendor/)
	@find . -path ./vendor -prune -o -name '*.go' -exec goimports -l -w {} +

update-deps:
	@go mod tidy
	@go get -u
	@go mod vendor

vuln-check:
	@govulncheck -show verbose ./...

