FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /slurm-mcp ./cmd/slurm-mcp/

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=builder /slurm-mcp /usr/local/bin/slurm-mcp
ENTRYPOINT ["slurm-mcp"]
CMD ["--transport=sse", "--port=8080"]
EXPOSE 8080
