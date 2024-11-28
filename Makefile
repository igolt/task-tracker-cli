tracker-cli: main.go
	go build -o $@ $<

clean:
	rm -f tracker-cli tasks.json

.PHONY: clean
