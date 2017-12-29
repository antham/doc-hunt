build-container:
	docker build -t antham/doc-hunt:$(v) -t antham/doc-hunt:latest .
fmt:
	find ! -path "./vendor/*" -name "*.go" -exec go fmt {} \;

gometalinter:
	gometalinter -D gotype --vendor --deadline=240s --dupl-threshold=200 -e '_string' -j 5 ./...

doc-hunt:
	doc-hunt check -e

run-tests:
	./test.sh

test-all: gometalinter run-tests doc-hunt

test-package:
	go test -race -cover -coverprofile=/tmp/doc-hunt github.com/antham/doc-hunt/$(pkg)
	go tool cover -html=/tmp/doc-hunt
