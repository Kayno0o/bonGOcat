FROM golang:latest

# Install required dependencies
RUN apt-get update && \
    apt-get install -y gcc libc6-dev libx11-dev xorg-dev libxtst-dev libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev build-essential git meson ninja-build bison pkg-config libxml2-dev libwayland-dev wayland-protocols

# Clone libxkbcommon repository and build it
RUN git clone https://github.com/xkbcommon/libxkbcommon.git && \
    cd libxkbcommon && \
    meson setup build --prefix=/usr/local && \
    ninja -C build && \
    ninja -C build install

WORKDIR /app

ENV GOOS=linux
ENV GOARCH=amd64

# Build the Go application with static linking
CMD ["sh", "-c", "go build -o bongo main.go"]
