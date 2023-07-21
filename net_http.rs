use std::io::{self, prelude::*};
use std::net::{SocketAddr, TcpStream};

struct Request {
    type_: ReqType,
    pattern: String,
    full_pattern: String,
    file: String,
    slugs: String,
}

impl Request {
    fn new() -> Self {
    }

    fn with_type_pat() -> Self {
    }

    fn get_type(&self) {
    }

    fn file(&self) -> &str {
        &self.pattern
    }

    fn file(&self) -> &str {
        &self.full_pattern
    }

    fn file(&self) -> &str {
        &self.file
    }
}

struct ResponseWriter {
    sock: TcpStream,
}

impl ResponseWriter {
    fn write_text(&mut self, text: &str) {
        // TODO: Return err?
        let _ = write!(
            self.sock,
            "{}\
            {}\r\n\
            CONTENT-TYPE: text-plain; charset=utf-8\r\n\
            CONTENT-LENGTH: {}\r\n\
            \r\n\
            {}",
            STATUS_OK,
            Self::get_time_string(),
            text.len(),
            text,
        );
    }

    fn write_file(&mut self, filepath: &str) {
        let Ok(mut file) = File::open(filepath) else {
            return;
        };
        let Ok(len) = file.metadata.map(|md| md.len()) else {
            return;
        };
        let ext = filepath.rsplit_once(".").map(|p| p.1).unwrap_or_default();
        let content_type = match ext {
            "html" | "htm" => "CONTENT-TYPE: text/html; charset=utf-8\r\n",
            "css" => "CONTENT-TYPE: text/css; charset=utf-8\r\n",
            "js" => "CONTENT-TYPE: text/javascript; charset=utf-8\r\n",
            //"ico" => "CONTENT-TYPE: image/x-icon\r\n",
            "ico" => "CONTENT-TYPE: image/vnd.microsoft.icon\r\n",
            "png" => "CONTENT-TYPE: image/png\r\n",
            "jpg" | "jpeg" => "CONTENT-TYPE: image/jpeg\r\n",
            "json" => "CONTENT-TYPE: application/json; charset=utf-8\r\n",
            "csv" => "CONTENT-TYPE: text/csv; charset=utf-8\r\n",
            "pdf" => "CONTENT-TYPE: application/pdf\r\n",
            _ => "CONTENT-TYPE: text/plain; charset=utf-8\r\n",
        };
        let _ = write!(
            self.sock,
            "{}\
            {}\r\n\
            {}\r\n\
            CONTENT-LENGTH: {}\r\n\
            \r\n",
            STATUS_OK,
            content_type,
            Self::get_time_string(),
            len,
        );
        let _ = io::copy(&mut file, &mut self.sock);
    }

    fn page_not_found(&mut self, filepath: &str) {
    }

    fn get_time_string() -> String {
    }
}

const BUFFER_LEN: usize = 5_000;
const STATUS_OK: &str = "HTTP/1.1 200 OK\r\n";
const STATUS_NOT_FOUND: &str = "HTTP/1.1 404 Not Found\r\n";
const HTML_404: &str = "\
<!DOCTYPE html>\
<html>\
<head><title>404 Not Found</title></head>\
<body>\
<h1>404 Not Found</h1>\
<p>The requested URL was not found on this server.</p>\
</body></html>";
const DAYS: [&str; 7] = ["Sun", "Mon", "tue", "Wed", "Thu", "Fri", "Sat"];
const MONTHS: [&str; 12] = [
    "Jan", "Feb", "Mar", "Apr",
    "May", "Jun", "Jul", "Aug",
    "Sep", "Oct", "Nov", "Dec",
];

pub struct HTTPServer {
    addr: SocketAddr,
    num_threads: usize,
    running: bool,
    default_pattern: String,
    allow_partial: bool,
    not_found_file_path: String,
    // Mutex
    // CondVar
    threads: Vec<JoinHandle>,
    sock_queue: LinkedList<TcpStream>,
    routes: HashMap<String, RouteHandler>,
}

impl HTTPServer {
    pub fn builder() -> HTTPServerBuilder {
        Default::default()
    }

    pub fn start(&mut self);
    pub fn start_nonblocking(&mut self);
    pub fn stop(&mut self);
    fn run_server(&mut self);
    fn handle_conns(...);
    fn parse_header(header: String) -> Request;
    fn submit_to_pool(&mut self, sock: TcpStream);
}

impl Into<HTTPServerBuilder> for HTTPServer {
    fn into(self) -> HTTPServerBuilder {
        HTTPServerBuilder(self)
    }
}

pub struct HTTPServerBuilder(HTTPServer);

impl HTTPServerBuilder {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn addr(mut self, addr: impl Into<SocketAddr>) -> Self {
        self.0.addr = addr.into();
        self
    }

    pub fn num_threads(mut self, num_threads: usize) -> Self {
        self.0.num_threads = num_threads;
        self
    }

    pub fn default_pattern(mut self, default_pattern: String) -> Self {
        self.0.default_pattern = default_pattern;
        self
    }

    pub fn allow_partial(mut self, allow_partial: bool) -> Self {
        self.0.allow_partial = allow_partial;
        self
    }

    pub fn not_found_file_path(mut self, not_found_file_path: String) -> Self {
        self.0.not_found_file_path = not_found_file_path;
        self
    }

    pub fn build(self) -> HTTPServer {
        self.0
    }
}

impl Default for HTTPServerBuilder {
    fn default() -> Self {
        HTTPServerBuilder(Default::default())
    }
}
