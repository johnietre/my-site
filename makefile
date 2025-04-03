.PHONY: go

go: cmd/my-site/main.go
	go build -o bin/my-site $^
