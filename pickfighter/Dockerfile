FROM golang:alpine as build

RUN apk --update add ca-certificates git

WORKDIR /build

COPY go.mod go.sum /build/
COPY pickfighter/main.go /build/pickfighter/
COPY pickfighter/cmd /build/pickfighter/cmd
COPY pickfighter/internal /build/pickfighter/internal
COPY pickfighter/pkg /build/pickfighter/pkg
COPY pickfighter/pkg /build/pickfighter/pkg

COPY pkg/discovery /build/pkg/discovery
COPY pkg/logger /build/pkg/logger
COPY pkg/model /build/pkg/model
COPY pkg/sigx /build/pkg/sigx
COPY pkg/pgxs /build/pkg/pgxs
COPY pkg/utils /build/pkg/utils
COPY pkg/ipaddr /build/pkg/ipaddr
COPY pkg/httplib /build/pkg/httplib

COPY internal/grpcutil /build/internal/grpcutil

COPY auth/pkg/model /build/auth/pkg/model
COPY auth/pkg/version /build/auth/pkg/version
COPY events/pkg/model /build/events/pkg/model
COPY fighters/pkg/model /build/fighters/pkg/model

COPY gen /build/gen

RUN go mod tidy

RUN GOOS=linux CGO_ENABLED=0 go build -o pickfighter-service ./pickfighter


FROM alpine:latest as runtime

RUN apk --update add ca-certificates

WORKDIR /app

COPY documentation/openapi/yaml /app/api-docs
COPY pickfighter/configs /app/configs
COPY pickfighter/hack/dev/certs /app/hack/dev/certs
COPY --from=build /build/pickfighter-service /app/

ENTRYPOINT ["/app/pickfighter-service"]
CMD ["serve", "gateway", "--env", "prod"]