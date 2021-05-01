build-container:
	docker build -t antham/doc-hunt:$(v) -t antham/doc-hunt:latest .

fmt:
	find ! -path "./vendor/*" -name "*.go" -exec go fmt {} \;

compile:
	gox -output "build/{{.Dir}}_{{.OS}}_{{.Arch}}"

lint:
	golangci-lint run

doc-hunt:
	doc-hunt check -e

gommit:
	gommit check range $(FROM) $(TO)

test-unit:
	./test.sh

test-package:
	go test -race -cover -coverprofile=/tmp/doc-hunt github.com/antham/doc-hunt/$(pkg)
	go tool cover -html=/tmp/doc-hunt
