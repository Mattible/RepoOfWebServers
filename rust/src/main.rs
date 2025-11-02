use rust_webserver::WebServer;

#[async_std::main]
async fn main() -> std::io::Result<()> {
    let server = WebServer::new();
    server.run_server().await
}
