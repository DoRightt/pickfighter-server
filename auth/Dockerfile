FROM golang:alpine as build

RUN apk --update add ca-certificates git

WORKDIR /build

COPY go.mod go.sum /build/
COPY auth/main.go /build/auth/
COPY auth/cmd /build/auth/cmd
COPY auth/internal /build/auth/internal
COPY auth/migrations /build/auth/migrations
COPY auth/pkg /build/auth/pkg
COPY pkg/discovery /build/pkg/discovery
COPY pkg/logger /build/pkg/logger
COPY pkg/model /build/pkg/model
COPY pkg/sigx /build/pkg/sigx
COPY pkg/pgxs /build/pkg/pgxs
COPY pkg/utils /build/pkg/utils
COPY gen /build/gen

RUN go mod tidy

RUN GOOS=linux CGO_ENABLED=0 go build -o auth-service ./auth


FROM alpine:latest as runtime

RUN apk --update add ca-certificates

WORKDIR /app

COPY documentation/openapi/yaml /app/api-docs
COPY auth/configs /app/configs
COPY auth/hack/dev/certs /app/hack/dev/certs
COPY --from=build /build/auth-service /app/

ENTRYPOINT ["/app/auth-service"]
CMD ["serve", "auth", "--env", "prod"]