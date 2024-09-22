FROM golang:1.23-alpine AS builder
ARG TARGETOS TARGETARCH
COPY . .
RUN go mod download
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o main .

FROM scratch AS final
ARG WEBSERVER_PORT=8000
WORKDIR /app
COPY --from=builder /go/main /app/
EXPOSE $WEBSERVER_PORT

CMD ["./main"]