lint:
	golint $(go list ./... | grep -v /vendor/)
# Exclude mock files from Go coverage reports

coverage:
	go test -coverprofile=coverage.out ./...
	grep -v "mock_" coverage.out | grep -v "example_main.go" | grep -v "/mock/" > coverage.filtered.out
	go tool cover -func=coverage.filtered.out
