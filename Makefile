clean:
	rm .file.* || true
	rm .memtable.* || true

dev: clean
	air

tidy:
	go mod tidy