FROM golang:1.24.1@sha256:52ff1b35ff8de185bf9fd26c70077190cd0bed1e9f16a2d498ce907e5c421268 AS build
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOARCH=amd64
WORKDIR /src
COPY . .
RUN go build -trimpath -ldflags "-s -w" -o /go/bin/tarantool-app


FROM scratch AS binary
LABEL author="a.burashnikov"
COPY --from=build /go/bin/tarantool-app /tarantool-app
ENTRYPOINT [ "/tarantool-app" ]
