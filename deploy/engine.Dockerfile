FROM golang:1.25.5 AS build
WORKDIR /build
COPY app ./app
COPY internal ./internal
COPY Makefile .
COPY go* .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN make engine

FROM scratch
COPY --from=build /build/bin/engine /bin/engine
WORKDIR /engine
CMD ["/bin/engine"]

