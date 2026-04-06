FROM golang:1.22-alpine AS builder

ARG TARGETARCH

WORKDIR /app

COPY src/go.mod ./
RUN go mod download

COPY src/ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -ldflags="-s -w" -o /app/server .

FROM scratch

ARG TARGETARCH
COPY --from=builder /app/server /server

EXPOSE 8080
ENV PORT=8080

ENTRYPOINT ["/server"]
