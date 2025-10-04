FROM rust:1.88-alpine AS builder
ARG TARGETOS TARGETARCH

LABEL org.opencontainers.image.authors="mattible@gmail.com" \
    version="1.0" \
    repository="https://github.com/Mattible/RepoOfWebServers" \
    description="The Rust WebServer from the RepoOfWebServers project"

WORKDIR /app
COPY . .
RUN apk add --no-cache musl-dev gcc libc-dev build-base
RUN cargo clean
RUN cargo build --release

FROM alpine:latest AS final
ARG WEBSERVER_PORT=8000
WORKDIR /app
COPY --from=builder /app/target/release/main /app/
EXPOSE $WEBSERVER_PORT

ENTRYPOINT ["./main"]