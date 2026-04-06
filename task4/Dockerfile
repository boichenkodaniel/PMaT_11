FROM rust:alpine AS builder

RUN apk add --no-cache musl-dev

WORKDIR /app

COPY src/Cargo.toml ./
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release --target x86_64-unknown-linux-musl
RUN rm -rf src target/x86_64-unknown-linux-musl/release/deps target/x86_64-unknown-linux-musl/release/.fingerprint

COPY src/ ./
RUN cargo build --release --target x86_64-unknown-linux-musl

FROM scratch

COPY --from=builder /app/target/x86_64-unknown-linux-musl/release/server /server

EXPOSE 8080
ENV PORT=8080

ENTRYPOINT ["/server"]
