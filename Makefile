GOCMD=go
GOTEST=go test -v

clean:
	rm .file.* || true
	rm .memtable.* || true
	rm Test_* || true

dev: clean
	air

test: clean
	$(GOTEST) ./...

tidy:
	$(GOCMD) mod tidy