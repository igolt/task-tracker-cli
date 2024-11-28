task-cli: main.go
	go build -o $@ $<

clean:
	rm -f task-cli tasks.json

.PHONY: clean
