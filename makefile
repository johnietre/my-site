.PHONY: go cpp

go: server/*.go
	go build -o bin/my-site $^

cpp:
	g++ server.cpp -o server -lpthread -lstdc++fs -std=c++17
