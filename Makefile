compile:
	git stash -u
	gox -output "build/{{.Dir}}_{{.OS}}_{{.Arch}}"

version:
	git stash -u
	sed -i "s/[[:digit:]]\+\.[[:digit:]]\+\.[[:digit:]]\+/$(v)/g" cmd/version.go
	git add -A
	git commit -m "feat(version) : "$(v)
	git tag v$(v) master

fmt:
	find ! -path "./vendor/*" -name "*.go" -exec go fmt {} \;

run-tests: test-all
test-all:
	./test.sh
test-package:
	go test -race -cover -coverprofile=/tmp/doc-hunt github.com/antham/doc-hunt/$(pkg)
	go tool cover -html=/tmp/doc-hunt
