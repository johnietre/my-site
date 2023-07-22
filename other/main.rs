mod net_http;

fn main() {
    let mut server = net_http::HTTPServer::new("127.0.0.1", 8000);
    server.handle_func("/", handler);
    server.handle_func("/static", dir);
    server.set_default_pattern("/");
    server.start();
}

fn handler(w: net_http::ResponseWriter, r: &mut tnet_http::Request) {
}

fn dir(w: net_http::ResponseWriter, r: &mut tnet_http::Request) {
    w.WriteFile("./static/" + r.getFile());
}
