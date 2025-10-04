use rust_webserver::WebServer;

fn main() -> std::io::Result<()> {
    let server = WebServer::new();
    server.run_server()
}
