FROM golang:alpine as build

RUN apk --update add ca-certificates git

WORKDIR /build

COPY go.mod go.sum /build/
COPY events/main.go /build/events/
COPY events/cmd /build/events/cmd
COPY events/internal /build/events/internal
COPY events/migrations /build/events/migrations
COPY events/pkg /build/events/pkg

COPY fighters/pkg/model /build/fighters/pkg/model

COPY pkg/discovery /build/pkg/discovery
COPY pkg/logger /build/pkg/logger
COPY pkg/model /build/pkg/model
COPY pkg/sigx /build/pkg/sigx
COPY pkg/pgxs /build/pkg/pgxs
COPY pkg/utils /build/pkg/utils
COPY pkg/utils /build/pkg/utils

COPY gen /build/gen

RUN go mod tidy

RUN GOOS=linux CGO_ENABLED=0 go build -o events-service ./events


FROM alpine:latest as runtime

RUN apk --update add ca-certificates

WORKDIR /app

COPY documentation/openapi/yaml /app/api-docs
COPY events/configs /app/configs
COPY --from=build /build/events-service /app/

ENTRYPOINT ["/app/events-service"]
CMD ["serve", "events", "--env", "prod"]