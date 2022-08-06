GOENV=CGO_ENABLED=0 GOFLAGS="-count=1"
GOCMD=$(GOENV) go
GOTEST=$(GOCMD) test -covermode=atomic -coverprofile=./coverage.out -timeout=20m


clean:
	@rm .file.* || true
	@rm .memtable.* || true
	@rm Test_* || true

dev: clean
	air

test: clean
	$(GOTEST) ./...

see-coverage:
	@go tool cover -html=coverage.out

tidy:
	$(GOCMD) mod tidy