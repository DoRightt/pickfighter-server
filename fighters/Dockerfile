FROM golang:alpine as build

RUN apk --update add ca-certificates git

WORKDIR /build

COPY go.mod go.sum /build/
COPY fighters/ /build/fighters/

COPY pickfighter/pkg/errors /build/pickfighter/pkg/errors
COPY auth/pkg/version /build/auth/pkg/version

COPY pkg/discovery /build/pkg/discovery
COPY pkg/logger /build/pkg/logger
COPY pkg/model /build/pkg/model
COPY pkg/sigx /build/pkg/sigx
COPY pkg/pgxs /build/pkg/pgxs
COPY pkg/utils /build/pkg/utils
COPY pkg/httplib /build/pkg/httplib

COPY gen /build/gen

RUN go mod tidy

RUN GOOS=linux CGO_ENABLED=0 go build -o fighters-service ./fighters


FROM alpine:latest as runtime

RUN apk --update add ca-certificates

WORKDIR /app

COPY documentation/openapi/yaml /app/api-docs
COPY fighters/configs /app/configs
COPY fighters/data /app/data
COPY --from=build /build/fighters-service /app/

ENTRYPOINT ["/app/fighters-service"]
CMD ["serve", "fighters", "--env", "prod"]