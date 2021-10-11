.PHONY: \
		init \
		test \
		test-coverage

.SILENT:

init:
	rm -f .git/hooks/pre-commit*
	cp ./.githooks/* .git/hooks

test:
	go test -v -coverprofile cover.out -race -count=1 ./...
	go tool cover -func cover.out | grep total | awk '{print "\n>>> Total coverage: " $$3}'
