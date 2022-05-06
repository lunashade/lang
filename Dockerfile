FROM alpine:latest

WORKDIR /lang

RUN apk add --update \
    bash \
    vim \
    git \
    make \
    musl-dev \
    go \
    curl \
    llvm \
    clang

ENTRYPOINT ["/bin/bash"]
