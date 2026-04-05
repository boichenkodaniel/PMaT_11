FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY src/go.mod ./
RUN go mod download

COPY src/ ./

RUN go test -v .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

FROM scratch

WORKDIR /
COPY --from=builder /app/server /server

EXPOSE 8080

ENV PORT=8080

ENTRYPOINT ["/server"]
