FROM golang:1.21

WORKDIR /work

COPY . /work

RUN make build

ENTRYPOINT ["/work/out"]
