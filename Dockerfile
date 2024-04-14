FROM golang:1.21.1 AS build-stage

WORKDIR /app

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /avito-banner ./cmd/main.go

FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /avito-banner /avito-banner
COPY ./config/config.yaml ./config.yaml

EXPOSE 8080

ENTRYPOINT ["/avito-banner"]
CMD ["--config", "/config.yaml"]
