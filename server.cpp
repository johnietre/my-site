#include "includes/http_server.hpp"
#include <iostream>
#include <string>
using namespace std;
using namespace Server;

void handler(ResponseWriter w, Request &r) {
  w.WriteFile("index.html");
}

void dir(ResponseWriter w, Request &r) {
  w.WriteFile("./static/" + r.getFile());
}

int main(int argc, char **argv) {
  HTTPServer server("192.168.1.125", 8000);
  server.handleFunc("/", &handler);
  server.handleFunc("/static", &dir);
  server.setDefaultPattern("/");
  server.start();
  return 0;
}
