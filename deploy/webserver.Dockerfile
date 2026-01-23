FROM golang:1.25.5 AS build
WORKDIR /build
COPY app ./app
COPY internal ./internal
COPY Makefile .
COPY go* .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN make webserver

FROM scratch
COPY --from=build /build/bin/webserver /bin/webserver
WORKDIR /webserver
CMD ["/bin/webserver"]

