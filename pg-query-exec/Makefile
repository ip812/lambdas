run-local:
	@docker build -f Dockerfile.local -t pg-query-exec:test .
	@docker run -p 9000:8080 pg-query-exec:test

invoke-local:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" \
	-H "Content-Type: application/json" \
	-d '{ \
		"database_name": "postgres", \
		"query": "SELECT datname FROM pg_database WHERE datistemplate = false;" \
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

