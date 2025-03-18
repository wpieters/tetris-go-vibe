FROM debian:stable-slim

ARG TARGETOS
ARG TARGETARCH

# Install stty for terminal control
RUN apt-get update && apt-get install -y coreutils && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY dist/tetris-${TARGETOS}-${TARGETARCH} /app/tetris

# Set TERM for proper terminal handling
ENV TERM=xterm-256color

CMD ["/app/tetris"]
