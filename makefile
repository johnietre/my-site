.PHONY: go

go: cmd/my-site/main.go
	CGO_ENABLED=1 go build -o bin/my-site $^
