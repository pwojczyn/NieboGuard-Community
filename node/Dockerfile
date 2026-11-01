FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -ldflags '-extldflags "-static"' -o /bin/nnode ./cmd/nnode

FROM alpine:latest
RUN apk add --no-cache sqlite-libs sqlite
WORKDIR /app
COPY --from=builder /bin/nnode /usr/local/bin/nnode
COPY configs/ /app/configs/
COPY web/ /app/web/
RUN mkdir -p /data /logs
ENV NB_NODE_CONFIG=/app/configs/default.yaml
ENV NB_NODE_WORKDIR=/data
ENV NB_NODE_HOST=0.0.0.0
ENV NB_NODE_PORT=8787
EXPOSE 8787 8787/tcp 8887/tcp
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD /usr/local/bin/nnode --verbose 2>/dev/null || exit 1
ENTRYPOINT ["/usr/local/bin/nnode"]
CMD ["--config=/app/configs/default.yaml"]
