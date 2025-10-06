# Build stage - compile Python to binary
FROM python:3.13.7-alpine AS builder
ARG TARGETOS TARGETARCH

LABEL org.opencontainers.image.authors="mattible@gmail.com" \
    version="1.0" \
    repository="https://github.com/Mattible/RepoOfWebServers" \
    description="The Python WebServer from the RepoOfWebServers project"

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache \
    gcc \
    musl-dev \
    linux-headers \
    binutils \
    upx

# Install PyInstaller and other dependencies
RUN pip install --no-cache-dir pyinstaller

# Copy the server file
COPY server.py .

# Create standalone binary
RUN pyinstaller --onefile \
    --name webserver \
    --distpath /app/dist \
    --workpath /app/build \
    --specpath /app \
    --strip \
    --upx-dir /usr/bin \
    server.py

# Compress the binary further
RUN upx --best --lzma /app/dist/webserver

RUN adduser -D appuser
RUN chown appuser:appuser /app/dist/webserver

# Production stage - minimal container with libc
FROM alpine AS final
ARG WEBSERVER_PORT=8000

# git sha from env
ARG GITSHA
ENV GITSHA=${GITSHA}

WORKDIR /app

# Copy the compiled binary from builder
COPY --from=builder /app/dist/webserver /app/webserver
COPY --from=builder /etc/passwd /etc/passwd

# Expose the port
ENV WEBSERVER_PORT=${WEBSERVER_PORT}
EXPOSE $WEBSERVER_PORT

# Switch to non-root user
USER appuser

# Run the binary
ENTRYPOINT ["/app/webserver"]
