FROM golang:1.24.1@sha256:52ff1b35ff8de185bf9fd26c70077190cd0bed1e9f16a2d498ce907e5c421268 AS cache
WORKDIR /modules
COPY go.mod go.sum ./
RUN go mod download


FROM golang:1.24.1@sha256:52ff1b35ff8de185bf9fd26c70077190cd0bed1e9f16a2d498ce907e5c421268 AS build
COPY --from=cache /go/pkg /go/pkg
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /go/bin/app /app/cmd/app

FROM scratch AS final
LABEL author="a.burashnikov"
COPY --from=build /go/bin/app /app
USER nobody
ENTRYPOINT [ "/app" ]
