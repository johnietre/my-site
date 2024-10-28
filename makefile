.PHONY: go cpp

go: cmd/my-site/main.go
	go build -o bin/my-site $^
