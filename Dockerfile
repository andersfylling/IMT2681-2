FROM golang:1.8
MAINTAINER https://github.com/andersfylling

WORKDIR /go/src/github.com/andersfylling/IMT2681-2
COPY . .

# Get Glide for package management
RUN curl https://glide.sh/get | sh
RUN glide install

ENV IMT_DATABASE_USER cloudtech
ENV IMT_DATABASE_PASS cloudtech
