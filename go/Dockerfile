FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder
ARG TARGETOS TARGETARCH BUILDPLATFORM

LABEL org.opencontainers.image.authors="mattible@gmail.com" \
    version="0.1.0" \
    repository="https://github.com/Mattible/RepoOfWebServers" \
    description="The Golang WebServer from the RepoOfWebServers project"

COPY . .
RUN go mod download
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o main .

RUN adduser -D appuser
RUN chown appuser:appuser /go/main

FROM scratch AS final
ARG WEBSERVER_PORT=8000
# git sha from env
ARG GITSHA
ENV GITSHA=${GITSHA}

WORKDIR /app
COPY --from=builder /go/main /app/
COPY --from=builder /etc/passwd /etc/passwd

ENV WEBSERVER_PORT=${WEBSERVER_PORT}
EXPOSE $WEBSERVER_PORT

USER appuser

ENTRYPOINT ["./main"]