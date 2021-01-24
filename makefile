.PHONY: go cpp

go:
	go build server.go

cpp:
	g++ server.cpp -o server -lpthread -std=c++17
