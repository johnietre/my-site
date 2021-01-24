.PHONY: go cpp

go:
	go build server.go

cpp:
	g++ -o server -lpthread server.cpp -std=c++17
