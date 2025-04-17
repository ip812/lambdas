run-local:
	@docker build -f Dockerfile.local -t ecr-push-notifier:test .
	@docker run -p 9000:8080 ecr-push-notifier:test

invoke-local:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" \
	-H "Content-Type: application/json" \
	-d '{ \
		"version": "0", \
		"id": "abcd-1234", \
		"detail-type": "ECR Image Action", \
		"source": "aws.ecr", \
		"account": "123456789012", \
		"time": "2025-04-17T12:00:00Z", \
		"region": "us-east-1", \
		"resources": [], \
		"detail": { \
			"action-type": "PUSH", \
			"repository-name": "ip812/ecr-push-notifier", \
			"image-tag": "69.1.0", \
			"image-digest": "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", \
			"artifact-media-type": "application/vnd.docker.container.image.v1+json", \
			"manifest-media-type": "application/vnd.docker.distribution.manifest.v2+json", \
			"result": "SUCCESS" \
		} \
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

