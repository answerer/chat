FROM golang:1.7

RUN apt-get update && apt-get install bzr

RUN curl https://glide.sh/get | sh