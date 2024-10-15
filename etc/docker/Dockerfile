ARG OS
ARG ARCH

FROM --platform=${OS}/${ARCH} golang:1.21-bullseye AS builder
ARG OS
ARG ARCH
ARG SOURCE
ARG VERSION

# Run makefile to build all the commands
WORKDIR /usr/src/app
COPY . .
RUN apt update && apt install -y libmosquitto-dev
RUN make

# Copy server to /usr/local/bin
FROM --platform=${OS}/${ARCH} debian:bullseye-slim
ARG OS
ARG ARCH
ARG SOURCE
ARG VERSION

RUN apt update && apt install -y libmosquitto-dev
COPY --from=builder /usr/src/app/build/* /usr/local/bin/
COPY etc/docker/server-entrypoint.sh .

# Labels
LABEL org.opencontainers.image.source=${SOURCE}

# Entrypoint when running the server
ENTRYPOINT [ "/server-entrypoint.sh" ]
CMD [ "/usr/local/bin/rotel" ]

