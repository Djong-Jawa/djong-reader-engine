.PHONY: test coverage coverage-report lint

# Run all tests
test:
	go test ./... -v

# Run tests, collect coverage, filter out auto-generated files, enforce 80% minimum
# -coverpkg=./... instruments every package so cross-package lines are not
# falsely reported as uncovered (e.g. resolver files tested by resolver tests).
coverage:
	go test ./... \
		-coverprofile=coverage_raw.out \
		-covermode=set \
		-coverpkg=./... 2>&1
	@grep -v "graph/resolvers/generated.go" coverage_raw.out > coverage.out
	@go tool cover -func=coverage.out | tee /dev/stderr | tail -1 | \
		awk '{gsub(/%/,""); if($$3 < 80.0) { print "FAIL: coverage " $$3 "% is below 80% minimum"; exit 1 } else { print "PASS: coverage " $$3 "%" }}'

# Open HTML coverage report in browser
coverage-report: coverage
	go tool cover -html=coverage.out

# Run vet
lint:
	go vet ./...
