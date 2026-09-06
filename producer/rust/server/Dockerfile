FROM ubuntu:20.04 as builder

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=America

RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    openssl libssl-dev \
    pkg-config \
    python \
    valgrind \
    zlib1g-dev

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

ENV PATH="/root/.cargo/bin:${PATH}"

RUN apt-get update && apt-get install -y cmake

COPY . /app

WORKDIR /app
RUN cargo build --release

FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=America

RUN apt-get update && apt-get install -y \
    openssl libssl-dev \
    zlib1g

COPY --from=builder /app/target/release/server /app/server

CMD ["/app/server"]
