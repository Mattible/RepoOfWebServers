FROM rust:1.91-alpine AS builder
ARG TARGETOS TARGETARCH

LABEL org.opencontainers.image.authors="mattible@gmail.com" \
    version="0.1.0" \
    repository="https://github.com/Mattible/RepoOfWebServers" \
    description="The Rust WebServer from the RepoOfWebServers project"

WORKDIR /app

COPY src/ src/
COPY Cargo.toml .

COPY . .

RUN apk add --no-cache musl-dev gcc libc-dev build-base

RUN cargo build --release

RUN adduser -D appuser
RUN chown appuser:appuser /app/target/release/main 

FROM scratch AS final
ARG WEBSERVER_PORT=8000
# git sha from env
ARG GITSHA
ENV GITSHA=${GITSHA}

WORKDIR /app
COPY --from=builder /app/target/release/main /app/
COPY --from=builder /etc/passwd /etc/passwd

ENV WEBSERVER_PORT=${WEBSERVER_PORT}
EXPOSE $WEBSERVER_PORT

USER appuser

ENTRYPOINT ["./main"]
